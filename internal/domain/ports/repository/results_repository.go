package repository

import (
	"context"
	"github.com/google/uuid"
	"tournaments-core/internal/domain/models"
)

type ResultsRepository interface {
	FetchById(ctx context.Context, id uuid.UUID) (models.Result, error)
	Update(ctx context.Context, updated *models.Result) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, r *models.Result) error
}
