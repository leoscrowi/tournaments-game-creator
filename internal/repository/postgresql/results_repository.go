package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"tournaments-core/internal/domain/models"
	"tournaments-core/internal/domain/ports/repository"
)

type resultsRepository struct {
	db *sql.DB
}

func NewResultsRepository(connect string) (repository.ResultsRepository, error) {
	db, err := sql.Open("postgres", connect)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &resultsRepository{db}, nil
}

func (r *resultsRepository) Create(ctx context.Context, res *models.Result) error {
	const op = "postgresql.ResultsRepository.Create"

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: Failed to begin transaction: %w", op, err)
	}

	query := `
	INSERT INTO game_creator.results (result_id, game_id, winner_id, comment)
	VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, query, res.ResultID, res.GameID, res.WinnerID, res.Comment)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: Failed to insert into results: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: Failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *resultsRepository) FetchById(ctx context.Context, id uuid.UUID) (models.Result, error) {
	const op = "postgresql.ResultsRepository.FetchById"

	query := `
	SELECT result_id, game_id, winner_id, comment
	FROM game_creator.results WHERE result_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var result models.Result
	err := row.Scan(&result.ResultID, &result.GameID, &result.WinnerID, &result.Comment)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Result{}, fmt.Errorf("%s: Result not found", op)
		}
		return models.Result{}, fmt.Errorf("%s: Failed to get result from db: %w", op, err)
	}

	return result, nil
}

func (r *resultsRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	const op = "postgresql.ResultsRepository.DeleteById"

	query := `
	DELETE FROM game_creator.results WHERE result_id = $1
	`

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: Failed to begin transaction: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: Failed to delete from results: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: Failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *resultsRepository) Update(ctx context.Context, updated *models.Result) error {
	const op = "postgresql.ResultsRepository.Update"

	query := `
        UPDATE game_creator.results 
        SET game_id = $1, 
            winner_id = $2, 
            comment = $3
        WHERE result_id = $4
    `

	result, err := r.db.ExecContext(ctx, query,
		updated.GameID,
		updated.WinnerID,
		updated.Comment,
		updated.ResultID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: result with id %s not found", op, updated.ResultID)
	}

	return nil
}
