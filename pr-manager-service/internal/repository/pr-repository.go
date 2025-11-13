package repository

import (
	"context"
	"database/sql"

	"pr-manager-service/internal/domain"
	uc "pr-manager-service/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

var _ uc.PullRequestRepositoryInterface = (*PullRequestRepository)(nil)

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error {
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

	insertPrSQL := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertPrSQL, pr.PullRequestId, pr.PullRequestName, pr.AuthorId)
	if err != nil {
		return err
	}

	insertReviewerSQL := `
		INSERT INTO reviewer_assignments (user_id, pull_request_id, slot)
		VALUES ($1, $2, $3)
	`

	for i, reviewerId := range pr.AssignedReviewers {
		slot := i + 1
		_, err = tx.Exec(ctx, insertReviewerSQL, reviewerId, pr.PullRequestId, slot)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PullRequestRepository) GetPullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	getPrSQL := `
		SELECT pull_request_id, pull_request_name, author_id, status_id
		FROM pull_requests
		WHERE pull_request_id = $1
	`
	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, getPrSQL, prId).
		Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.StatusId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	getReviewersSQL := `
		SELECT user_id
		FROM reviewer_assignments
		WHERE pull_request_id = $1
		ORDER BY slot
	`
	rows, err := r.pool.Query(ctx, getReviewersSQL, prId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PullRequestRepository) MergePullRequest(ctx context.Context, prId string) (*domain.PullRequest, error) {
	updateSQL := `
		UPDATE pull_requests
		SET status_id = 2,
		    mergedAt  = COALESCE(mergedAt, CURRENT_TIMESTAMP),
		    updated_at = CURRENT_TIMESTAMP
		WHERE pull_request_id = $1
		RETURNING pull_request_id, pull_request_name, author_id, status_id
	`

	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, updateSQL, prId).
		Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.StatusId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Get reviewers
	getReviewersSQL := `
		SELECT user_id
		FROM reviewer_assignments
		WHERE pull_request_id = $1
		ORDER BY slot
	`
	rows, err := r.pool.Query(ctx, getReviewersSQL, prId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PullRequestRepository) GetAllPrByUserId(ctx context.Context, userId string) ([]domain.PullRequest, error) {
	querySQL := `
		SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status_id
		FROM pull_requests p
		JOIN reviewer_assignments r ON p.pull_request_id = r.pull_request_id
		WHERE r.user_id = $1
		ORDER BY p.created_at
	`
	rows, err := r.pool.Query(ctx, querySQL, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		err = rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.StatusId)
		if err != nil {
			return nil, err
		}
		result = append(result, pr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PullRequestRepository) ReplaceReviewer(ctx context.Context, prId string, oldUserId string, newUserId string) error {
	updateSQL := `
		UPDATE reviewer_assignments
		SET user_id = $3
		WHERE pull_request_id = $1 AND user_id = $2
	`
	ct, err := r.pool.Exec(ctx, updateSQL, prId, oldUserId, newUserId)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PullRequestRepository) GetActiveTeamMembers(ctx context.Context, teamName string) ([]domain.User, error) {
	querySQL := `
		SELECT u.user_id, u.username, u.is_active
		FROM memberships m
		JOIN users u ON u.user_id = m.user_id
		WHERE m.team_id = $1
		  AND u.is_active = true
	`
	rows, err := r.pool.Query(ctx, querySQL, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		err = rows.Scan(&u.UserId, &u.UserName, &u.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
