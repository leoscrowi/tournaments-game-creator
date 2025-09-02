package repository

import (
	"context"
	"github.com/google/uuid"
	"tournaments-core/internal/domain/models"
)

type GamesRepository interface {
	FetchById(ctx context.Context, id uuid.UUID) (models.Game, error)
	Update(ctx context.Context, updated *models.Game) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, g *models.Game) error
}
