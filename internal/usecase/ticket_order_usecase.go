package usecase

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
)

type TicketOrderUseCaseInterface interface {
	TicketOrderSave
	TicketOrderFindById
	TicketOrderGetAll
}

type TicketOrderSave interface {
	Save(context context.Context, ticket domain.TicketOrder) (dto.TicketOrderResponseSave, error)
}

type TicketOrderFindById interface {
	FindById(context context.Context, id int) (dto.TicketOrderResponse, error)
}

type TicketOrderGetAll interface {
	GetAll(context context.Context) ([]dto.TicketOrderResponse, error)
}