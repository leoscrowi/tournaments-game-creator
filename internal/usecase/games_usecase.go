package usecase

import (
	"context"
	"github.com/google/uuid"
	"time"
	"tournaments-core/internal/domain/models"
	"tournaments-core/internal/domain/ports/repository"
	"tournaments-core/internal/domain/ports/usecase"
)

type gamesUseCase struct {
	gamesRepository repository.GamesRepository
	contextTimeout  time.Duration
}

func NewGamesUseCase(gamesRepository repository.GamesRepository, timeout time.Duration) usecase.GamesUseCase {
	return &gamesUseCase{
		gamesRepository: gamesRepository,
		contextTimeout:  timeout,
	}
}

func (gu *gamesUseCase) FetchById(ctx context.Context, id uuid.UUID) (models.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()
	return gu.gamesRepository.FetchById(ctx, id)
}

func (gu *gamesUseCase) Update(ctx context.Context, updated *models.Game) error {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()
	return gu.gamesRepository.Update(ctx, updated)
}

func (gu *gamesUseCase) DeleteById(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()
	return gu.gamesRepository.DeleteById(ctx, id)
}

func (gu *gamesUseCase) Create(ctx context.Context, g *models.Game) error {
	ctx, cancel := context.WithTimeout(ctx, gu.contextTimeout)
	defer cancel()
	return gu.gamesRepository.Create(ctx, g)

}
