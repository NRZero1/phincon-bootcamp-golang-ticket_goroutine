package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type TicketRepositoryInterface interface {
	TicketSave
	TicketFindById
	TicketGetAll
}

type TicketSave interface {
	Save(context context.Context, ticket *domain.Ticket) (domain.Ticket, error)
}

type TicketFindById interface {
	FindByID(context context.Context, id int) (domain.Ticket, error)
}

type TicketGetAll interface {
	GetAll(context context.Context) ([]domain.Ticket, error)
}