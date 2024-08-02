package impl

import (
	"context"
	"fmt"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type TicketOrderUseCase struct {
	repoTicketOrder repository.TicketOrderRepositoryInterface
	useCaseTicket usecase.TicketUseCaseInterface
	useCaseEvent usecase.EventUseCaseInterface
	useCaseUser usecase.UserUseCaseInterface
}

func NewTicketOrderUseCase(repoTicketOrder repository.TicketOrderRepositoryInterface,
	useCaseTicket usecase.TicketUseCaseInterface,
	useCaseEvent usecase.EventUseCaseInterface,
	useCaseUser usecase.UserUseCaseInterface) (usecase.TicketOrderUseCaseInterface) {
	return TicketOrderUseCase{
		repoTicketOrder: repoTicketOrder,
		useCaseTicket: useCaseTicket,
		useCaseEvent: useCaseEvent,
		useCaseUser: useCaseUser,
	}
}

func (uc TicketOrderUseCase) Save(ctx context.Context, ticketOrder domain.TicketOrder) (dto.TicketOrderResponseSave, error) {
	log.Trace().Msg("Entering ticket order usecase save")

	log.Info().Msg("Attempting to check ticket if exist")
	foundTicket, errTicket := uc.useCaseTicket.FindById(ctx, ticketOrder.TicketID)

	if errTicket != nil {
		return dto.TicketOrderResponseSave{}, errTicket
	}

	log.Info().Msg("Attempting to check event if exist")
	foundEvent, errEvent := uc.useCaseEvent.FindById(ctx, foundTicket.EventDetails.EventID)

	if errEvent != nil {
		return dto.TicketOrderResponseSave{}, errEvent
	}

	log.Info().Msg("Attempting to check user if exist")
	foundUser, errUser := uc.useCaseUser.FindById(ctx, ticketOrder.UserID)

	if errUser != nil {
		return dto.TicketOrderResponseSave{}, errUser
	}

	if foundTicket.Stock < ticketOrder.Amount {
		return dto.TicketOrderResponseSave{}, fmt.Errorf("not enough stock")
	}

	var totalPrice float64 = float64(ticketOrder.Amount) * foundTicket.Price

	if foundUser.Balance < totalPrice {
		return dto.TicketOrderResponseSave{}, fmt.Errorf("not enough balance")
	}

	deductedTicket, errDeduct := uc.useCaseTicket.Deduct(ctx, foundTicket.TicketID, ticketOrder.Amount)

	if errDeduct != nil {
		return dto.TicketOrderResponseSave{}, errDeduct
	}

	reducedUserBalance, errBalance := uc.useCaseUser.ReduceBalance(ctx, foundUser.UserID, totalPrice)
	
	if errBalance != nil {
		uc.useCaseTicket.Restore(ctx, deductedTicket.TicketID, ticketOrder.Amount)
	}

	log.Info().Msg("Attempting to call ticket order repo save")
	ticketSave := domain.TicketOrder {
		OrderID: ticketOrder.OrderID,
		TicketID: ticketOrder.TicketID,
		UserID: ticketOrder.UserID,
		Amount: ticketOrder.Amount,
		TotalPrice: totalPrice,
	}
	errSaved := uc.repoTicketOrder.Save(ctx, &ticketSave)

	if errSaved != nil {
		return dto.TicketOrderResponseSave{}, errSaved
	}

	ticketOrderResponse := dto.TicketOrderResponseSave {
		OrderID: ticketSave.OrderID,
		TicketDetails: dto.TicketResponse{
			TicketID: deductedTicket.TicketID,
			EventDetails: dto.EventResponse{
				EventID: foundEvent.EventID,
				EventName: foundEvent.EventName,
			},
			Name: deductedTicket.Name,
			Price: deductedTicket.Price,
			Stock: deductedTicket.Stock,
			Type: deductedTicket.Type,
		},
		UserDetails: dto.UserResponse{
			UserID: reducedUserBalance.UserID,
			Email: reducedUserBalance.Email,
			Name: reducedUserBalance.Name,
			PhoneNumber: reducedUserBalance.PhoneNumber,
			RemainingBalance: reducedUserBalance.Balance,
		},
		Amount: ticketSave.Amount,
		TotalPrice: totalPrice,
	}
	return ticketOrderResponse, nil
}

func (uc TicketOrderUseCase) FindById(ctx context.Context, id int) (dto.TicketOrderResponse, error) {
	log.Trace().Msg("Entering ticket order usecase find by id")
	log.Info().Msg("Attempting to call ticket order repo to check if ticket order exist")
	foundTicketOrder, errTicketOrder := uc.repoTicketOrder.FindByID(ctx, id)

	if errTicketOrder != nil {
		return dto.TicketOrderResponse{}, errTicketOrder
	}
	
	log.Info().Msg("Attempting to call ticket use case to check if ticket exist")
	foundTicket, errTicket := uc.useCaseTicket.FindById(ctx, foundTicketOrder.TicketID)

	if errTicket != nil {
		return dto.TicketOrderResponse{}, errTicket
	}

	log.Info().Msg("Attempting to call event use case to check if event exist")
	foundEvent, errEvent := uc.useCaseEvent.FindById(ctx, foundTicket.EventDetails.EventID)

	if errEvent != nil {
		return dto.TicketOrderResponse{}, errEvent
	}

	log.Info().Msg("Attempting to call user use case to check if user exist")
	foundUser, errUser := uc.useCaseUser.FindById(ctx, foundTicketOrder.UserID)

	if errUser != nil {
		return dto.TicketOrderResponse{}, errUser
	}

	ticketOrderResponse := dto.TicketOrderResponse {
		OrderID: foundTicketOrder.OrderID,
		TicketDetails: dto.TicketResponseTicketOrder{
			TicketID: foundTicket.TicketID,
			EventDetails: dto.EventResponse{
				EventID: foundEvent.EventID,
				EventName: foundEvent.EventName,
			},
			Name: foundTicket.Name,
			Price: foundTicket.Price,
			Type: foundTicket.Type,
		},
		UserDetails: dto.UserResponseTicketOrder{
			UserID: foundUser.UserID,
			Email: foundUser.Email,
			Name: foundUser.Name,
			PhoneNumber: foundUser.PhoneNumber,
		},
		Amount: foundTicketOrder.Amount,
		TotalPrice: foundTicketOrder.TotalPrice,
	}
	return ticketOrderResponse, nil
}

func (uc TicketOrderUseCase) GetAll(ctx context.Context) ([]dto.TicketOrderResponse, error) {
	log.Trace().Msg("Entering ticket order usecase get all")
	log.Info().Msg("Attempting to call fetch all ticket")
	listOfTicketOrder, err := uc.repoTicketOrder.GetAll(ctx)

	if err != nil {
		return []dto.TicketOrderResponse{}, err
	}

	var listOfTicketOrderResponse []dto.TicketOrderResponse

	for _, v := range listOfTicketOrder {
		log.Info().Msg("Attempting to call ticket use case to find ticket by id")
		foundTicket, errTicket := uc.useCaseTicket.FindById(ctx, v.TicketID)

		if errTicket != nil {
			return []dto.TicketOrderResponse{}, errTicket
		}

		log.Info().Msg("Attempting to call event use case to find event by id")
		foundEvent, errEvent := uc.useCaseEvent.FindById(ctx, foundTicket.EventDetails.EventID)

		if errEvent != nil {
			return []dto.TicketOrderResponse{}, errEvent
		}

		log.Info().Msg("Attempting to call user use case to find user by id")
		foundUser, errUser := uc.useCaseUser.FindById(ctx, v.UserID)

		if errUser != nil {
			return []dto.TicketOrderResponse{}, errUser
		}

		ticketOrderResponse := dto.TicketOrderResponse {
			OrderID: v.OrderID,
			TicketDetails: dto.TicketResponseTicketOrder{
				TicketID: foundTicket.TicketID,
				EventDetails: dto.EventResponse{
					EventID: foundEvent.EventID,
					EventName: foundEvent.EventName,
				},
				Name: foundTicket.Name,
				Price: foundTicket.Price,
				Type: foundTicket.Type,
			},
			UserDetails: dto.UserResponseTicketOrder{
				UserID: foundUser.UserID,
				Email: foundUser.Email,
				Name: foundUser.Name,
				PhoneNumber: foundUser.PhoneNumber,
			},
			Amount: v.Amount,
			TotalPrice: v.TotalPrice,
		}

		listOfTicketOrderResponse = append(listOfTicketOrderResponse, ticketOrderResponse)
	}
	return listOfTicketOrderResponse, nil
}