package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type TicketOrderRepositoryInterface interface {
	TicketOrderSave
	TicketOrderFindById
	TicketOrderGetAll
}

type TicketOrderSave interface {
	Save(context context.Context, ticket *domain.TicketOrder) (error)
}

type TicketOrderFindById interface {
	FindByID(context context.Context, id int) (domain.TicketOrder, error)
}

type TicketOrderGetAll interface {
	GetAll(context context.Context) ([]domain.TicketOrder, error)
}