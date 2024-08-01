package usecase

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
)

type TicketUseCaseInterface interface {
	TicketSave
	TicketFindById
	TicketGetAll
	TicketDeduct
	// TicketRestore
}

type TicketSave interface {
	Save(context context.Context, ticket domain.Ticket) (dto.TicketResponse, error)
}

type TicketFindById interface {
	FindById(context context.Context, id int) (dto.TicketResponse, error)
}

type TicketGetAll interface {
	GetAll(context context.Context) ([]dto.TicketResponse, error)
}

type TicketDeduct interface {
	Deduct(context context.Context, id int, amount int) (domain.Ticket, error)
}

type TicketRestore interface {
	Restore(context context.Context, id int, amount int)
}