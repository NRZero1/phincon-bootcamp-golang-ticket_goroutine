package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type TicketRepositoryInterface interface {
	TicketSave
	TicketFindById
	TicketGetAll
	TicketDeduct
	TicketRestore
}

type TicketSave interface {
	Save(context context.Context, ticket *domain.Ticket) (error)
}

type TicketFindById interface {
	FindByID(context context.Context, id int) (domain.Ticket, error)
}

type TicketGetAll interface {
	GetAll(context context.Context) ([]domain.Ticket, error)
}

type TicketDeduct interface {
	Deduct(context context.Context, id int, amount int) (domain.Ticket, error)
}

type TicketRestore interface {
	Restore(context context.Context, id int, amount int)
}