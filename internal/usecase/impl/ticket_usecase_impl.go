package impl

import (
	"context"
	"errors"
	"fmt"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type TicketUseCase struct {
	repoTicket repository.TicketRepositoryInterface
	useCaseEvent usecase.EventUseCaseInterface
}

func NewTicketUseCase(repoTicket repository.TicketRepositoryInterface, ucEvent usecase.EventUseCaseInterface) (usecase.TicketUseCaseInterface) {
	return TicketUseCase{
		repoTicket: repoTicket,
		useCaseEvent: ucEvent,
	}
}

func (uc TicketUseCase) Save(ctx context.Context, ticket domain.Ticket) (dto.TicketResponse, error) {
	log.Trace().Msg("Entering ticket usecase save")
	log.Info().Msg("Attempting to check event if exist")

	if uc.useCaseEvent == nil {
        log.Error().Msg("useCaseEvent is nil")
        return dto.TicketResponse{}, errors.New("useCaseEvent is not initialized")
    }

	foundEvent, errEvent := uc.useCaseEvent.FindById(ctx, ticket.EventID)

	if errEvent != nil {
		return dto.TicketResponse{}, errEvent
	}

	log.Info().Msg("Attempting to call ticket repo save")
	errSaved := uc.repoTicket.Save(ctx, &ticket)

	if errSaved != nil {
		return dto.TicketResponse{}, errSaved
	}

	ticketResponse := dto.TicketResponse {
		TicketID: ticket.TicketID,
		EventDetails: dto.EventResponse{
			EventID: foundEvent.EventID,
			EventName: foundEvent.EventName,
		},
		Name: ticket.Name,
		Price: ticket.Price,
		Stock: ticket.Stock,
		Type: ticket.Type,
	}
	return ticketResponse, nil
}

func (uc TicketUseCase) FindById(ctx context.Context, id int) (dto.TicketResponse, error) {
	log.Trace().Msg("Entering ticket usecase find by id")
	log.Info().Msg("Attempting to call ticket repo to check if ticket exist")
	foundTicket, errTicket := uc.repoTicket.FindByID(ctx, id)

	if errTicket != nil {
		return dto.TicketResponse{}, errTicket
	}
	
	log.Info().Msg("Attempting to call event repo to check if event exist")
	foundEvent, errEvent := uc.useCaseEvent.FindById(ctx, foundTicket.EventID)

	if errEvent != nil {
		return dto.TicketResponse{}, errEvent
	}

	ticketResponse := dto.TicketResponse {
		TicketID: foundTicket.TicketID,
		EventDetails: dto.EventResponse{
			EventID: foundEvent.EventID,
			EventName: foundEvent.EventName,
		},
		Name: foundTicket.Name,
		Price: foundTicket.Price,
		Stock: foundTicket.Stock,
		Type: foundTicket.Type,
	}
	return ticketResponse, nil
}

func (uc TicketUseCase) GetAll(ctx context.Context) ([]dto.TicketResponse, error) {
	log.Trace().Msg("Entering ticket usecase get all")
	log.Info().Msg("Attempting to call fetch all ticket")
	listOfTicket, err := uc.repoTicket.GetAll(ctx)

	if err != nil {
		return []dto.TicketResponse{}, err
	}

	var listOfTicketResponse []dto.TicketResponse

	for _, v := range listOfTicket {
		foundEvent, err := uc.useCaseEvent.FindById(ctx, v.EventID)

		if err != nil {
			return []dto.TicketResponse{}, err
		}

		ticketResponse := dto.TicketResponse {
			TicketID: v.TicketID,
			EventDetails: dto.EventResponse{
				EventID: foundEvent.EventID,
				EventName: foundEvent.EventName,
			},
			Name: v.Name,
			Price: v.Price,
			Stock: v.Stock,
			Type: v.Type,
		}

		listOfTicketResponse = append(listOfTicketResponse, ticketResponse)
	}
	return listOfTicketResponse, nil
}

func (uc TicketUseCase) Deduct(ctx context.Context, id int, amount int) (domain.Ticket, error) {
	log.Trace().Msg("Entering ticket usecase deduct")
	log.Info().Msg("Attempting to call ticket repo to deduct stock")
	ticket, errTicket := uc.repoTicket.Deduct(ctx, id, amount)

	if errTicket != nil {
		return domain.Ticket{}, fmt.Errorf("")
	}

	return ticket, nil
}

func (uc TicketUseCase) Restore(ctx context.Context, id int, amount int) {
	log.Trace().Msg("Entering ticket usecase restore")
	log.Info().Msg("Attempting to call ticket repo to restore stock")
	uc.repoTicket.Restore(ctx, id, amount)
}
