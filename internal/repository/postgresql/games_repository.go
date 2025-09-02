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

type gamesRepository struct {
	db *sql.DB
}

func NewGamesRepository(connect string) (repository.GamesRepository, error) {
	db, err := sql.Open("postgres", connect)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &gamesRepository{db}, nil
}

func (r *gamesRepository) Create(ctx context.Context, g *models.Game) error {
	const op = "postgresql.GamesRepository.Create"

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: Failed to begin transaction: %w", op, err)
	}

	query := `
	INSERT INTO game_creator.games (game_id, game_start, game_type_id)
	VALUES ($1, $2, $3)
	`

	_, err = tx.ExecContext(ctx, query, g.GameID.String(), g.GameStart, g.GameTypeID.String())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: Failed to insert into games: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: Failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *gamesRepository) FetchById(ctx context.Context, id uuid.UUID) (models.Game, error) {
	const op = "postgresql.GamesRepository.FetchById"

	query := `
	SELECT game_id, game_start, game_type_id
	FROM game_creator.games WHERE game_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var game models.Game
	err := row.Scan(&game.GameID, &game.GameStart, &game.GameTypeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Game{}, fmt.Errorf("%s: Game not found", op)
		}
		return models.Game{}, fmt.Errorf("%s: Failed to get game from db: %w", op, err)
	}

	return game, nil
}

func (r *gamesRepository) Update(ctx context.Context, updated *models.Game) error {
	const op = "postgresql.GamesRepository.Update"

	query := `
	UPDATE game_creator.games
	SET 
	    game_start=COALESCE($1, game_start),
	    game_type_id=COALESCE($2, game_type_id)
	WHERE game_id=$3
	`

	var nullTime sql.NullTime
	if !updated.GameStart.IsZero() {
		nullTime = sql.NullTime{Time: updated.GameStart, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query,
		nullTime,
		updated.GameTypeID,
		updated.GameID,
	)

	if err != nil {
		return fmt.Errorf("%s: failed to update game: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: game with id %s not found", op, updated.GameID)
	}

	return nil
}

func (r *gamesRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	const op = "postgresql.GamesRepository.DeleteById"

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: Failed to begin transaction: %w", op, err)
	}

	query := `
	DELETE FROM game_creator.games WHERE game_id = $1
	`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: Failed to delete from games: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: Failed to commit transaction: %w", op, err)
	}

	return nil
}
