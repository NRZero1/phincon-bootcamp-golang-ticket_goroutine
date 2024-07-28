package impl

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type EventUseCase struct {
	repo repository.EventRepositoryInterface
}

func NewEventUseCase(repo repository.EventRepositoryInterface) (usecase.EventUseCaseInterface) {
	return EventUseCase{
		repo: repo,
	}
}

func (uc EventUseCase) Save(ctx context.Context, event domain.Event) (domain.Event, error) {
	log.Trace().Msg("Entering event usecase save")
	return uc.repo.Save(ctx, &event)
}

func (uc EventUseCase) FindById(ctx context.Context, id int) (domain.Event, error) {
	log.Trace().Msg("Entering event usecase find by id")
	return uc.repo.FindByID(ctx, id)
}

func (uc EventUseCase) GetAll(ctx context.Context) ([]domain.Event, error) {
	log.Trace().Msg("Entering event usecase get all")
	return uc.repo.GetAll(ctx)
}