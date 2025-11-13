package repository

import (
	"context"
	"database/sql"

	"pr-manager-service/internal/domain"
	uc "pr-manager-service/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

var _ uc.UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetUser(ctx context.Context, userId string) (*domain.User, error) {
	getUserSQL := `
		SELECT user_id, username, is_active
		FROM users
		WHERE user_id = $1
	`
	var u domain.User
	err := r.pool.QueryRow(ctx, getUserSQL, userId).Scan(&u.UserId, &u.UserName, &u.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, string, error) {
	updateSQL := `
		UPDATE users
		SET is_active = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
		RETURNING user_id, username, is_active
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, updateSQL, userId, isActive).
		Scan(&u.UserId, &u.UserName, &u.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", sql.ErrNoRows
		}
		return nil, "", err
	}

	getTeamSQL := `
		SELECT team_id
		FROM memberships
		WHERE user_id = $1
		ORDER BY team_id
		LIMIT 1
	`
	var teamName string
	err = r.pool.QueryRow(ctx, getTeamSQL, userId).Scan(&teamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &u, "", nil
		}
		return nil, "", err
	}

	return &u, teamName, nil
}

func (r *UserRepository) GetTeamName(ctx context.Context, userId string) (string, error) {
	getTeamSQL := `
		SELECT team_id
		FROM memberships
		WHERE user_id = $1
		ORDER BY team_id
		LIMIT 1
	`

	var teamName string
	err := r.pool.QueryRow(ctx, getTeamSQL, userId).Scan(&teamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	return teamName, nil
}

