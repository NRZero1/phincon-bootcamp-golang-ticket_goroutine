package impl

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type TicketUseCase struct {
	repoTicket repository.TicketRepositoryInterface
	repoEvent repository.EventRepositoryInterface
}

func NewTicketUseCase(repoTicket repository.TicketRepositoryInterface, repoEvent repository.EventRepositoryInterface) (usecase.TicketUseCaseInterface) {
	return TicketUseCase{
		repoTicket: repoTicket,
		repoEvent: repoEvent,
	}
}

func (uc TicketUseCase) Save(ctx context.Context, ticket domain.Ticket) (dto.TicketResponse, error) {
	log.Trace().Msg("Entering ticket usecase save")
	log.Info().Msg("Attempting to check event if exist")
	foundEvent, errEvent := uc.repoEvent.FindByID(ctx, ticket.EventID)

	if errEvent != nil {
		return dto.TicketResponse{}, errEvent
	}

	log.Info().Msg("Attempting to call ticket repo save")
	savedTicket, errSaved := uc.repoTicket.Save(ctx, &ticket)

	if errSaved != nil {
		return dto.TicketResponse{}, errSaved
	}

	ticketResponse := dto.TicketResponse {
		TicketID: savedTicket.TicketID,
		EventDetails: dto.EventResponse{
			EventID: foundEvent.EventID,
			EventName: foundEvent.EventName,
		},
		Name: savedTicket.Name,
		Price: savedTicket.Price,
		Stock: savedTicket.Stock,
		Type: savedTicket.Type,
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
	foundEvent, errEvent := uc.repoEvent.FindByID(ctx, foundTicket.EventID)

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
		foundEvent, err := uc.repoEvent.FindByID(ctx, v.EventID)

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