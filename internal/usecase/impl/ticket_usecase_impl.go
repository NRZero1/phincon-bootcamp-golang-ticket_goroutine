package impl

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type TicketUseCase struct {
	repo repository.TicketRepositoryInterface
}

func NewTicketUseCase(repo repository.TicketRepositoryInterface) (usecase.TicketUseCaseInterface) {
	return TicketUseCase{
		repo: repo,
	}
}

func (uc TicketUseCase) Save(ctx context.Context, ticket domain.Ticket) (domain.Ticket, error) {
	log.Trace().Msg("Entering ticket usecase save")
	return uc.repo.Save(ctx, &ticket)
}

func (uc TicketUseCase) FindById(ctx context.Context, id int) (domain.Ticket, error) {
	log.Trace().Msg("Entering ticket usecase find by id")
	return uc.repo.FindByID(ctx, id)
}

func (uc TicketUseCase) GetAll(ctx context.Context) ([]domain.Ticket, error) {
	log.Trace().Msg("Entering ticket usecase get all")
	return uc.repo.GetAll(ctx)
}