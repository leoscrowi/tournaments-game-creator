package usecase

import (
	"context"
	"github.com/google/uuid"
	"time"
	"tournaments-core/internal/domain/models"
	"tournaments-core/internal/domain/ports/repository"
	"tournaments-core/internal/domain/ports/usecase"
)

type resultsUseCase struct {
	resultRepository repository.ResultsRepository
	contextTimeout   time.Duration
}

func NewResultsUseCase(r repository.ResultsRepository, timeout time.Duration) usecase.ResultsUseCase {
	return &resultsUseCase{
		resultRepository: r,
		contextTimeout:   timeout,
	}
}

func (ru *resultsUseCase) FetchById(ctx context.Context, id uuid.UUID) (models.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()
	return ru.resultRepository.FetchById(ctx, id)
}

func (ru *resultsUseCase) DeleteById(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()
	return ru.resultRepository.DeleteById(ctx, id)
}

func (ru *resultsUseCase) Create(ctx context.Context, r *models.Result) error {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()
	return ru.resultRepository.Create(ctx, r)
}

func (ru *resultsUseCase) Update(ctx context.Context, updated *models.Result) error {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()
	return ru.resultRepository.Update(ctx, updated)
}
