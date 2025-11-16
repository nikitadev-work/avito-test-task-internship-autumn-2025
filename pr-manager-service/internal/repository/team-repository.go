package repository

import (
	"context"
	"database/sql"

	"pr-manager-service/internal/domain"
	uc "pr-manager-service/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

var _ uc.TeamRepositoryInterface = (*TeamRepository)(nil)

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, teamName string, members []domain.User) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// Create team
	createTeamSQL := `
		INSERT INTO teams (team_name)
		VALUES ($1)
	`
	_, err = tx.Exec(ctx, createTeamSQL, teamName)
	if err != nil {
		return err
	}

	upsertUserSQL := `
		INSERT INTO users (user_id, username, is_active)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			username  = EXCLUDED.username,
			is_active = EXCLUDED.is_active
	`
	insertMembershipSQL := `
		INSERT INTO memberships (user_id, team_name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, team_name) DO NOTHING
	`

	// Create users and memberships for them
	for _, u := range members {
		_, err = tx.Exec(ctx, upsertUserSQL, u.UserId, u.UserName, u.IsActive)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, insertMembershipSQL, u.UserId, teamName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, []domain.User, error) {
	getTeamSQL := `
		SELECT team_name
		FROM teams
		WHERE team_name = $1
	`
	// Check if the team exists
	var t domain.Team
	err := r.pool.QueryRow(ctx, getTeamSQL, teamName).Scan(&t.TeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, sql.ErrNoRows
		}
		return nil, nil, err
	}

	getMembersSQL := `
		SELECT u.user_id, u.username, u.is_active
		FROM memberships m
		JOIN users u ON u.user_id = m.user_id
		WHERE m.team_name = $1
		ORDER BY u.username
	`
	// Get all members
	rows, err := r.pool.Query(ctx, getMembersSQL, teamName)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var members []domain.User
	for rows.Next() {
		var u domain.User
		err = rows.Scan(&u.UserId, &u.UserName, &u.IsActive)
		if err != nil {
			return nil, nil, err
		}
		members = append(members, u)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return &t, members, nil
}
