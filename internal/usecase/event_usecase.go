package usecase

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type EventUseCaseInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(context context.Context, event domain.Event) (domain.Event, error)
}

type EventFindById interface {
	FindById(context context.Context, id int) (domain.Event, error)
}

type EventGetAll interface {
	GetAll(context context.Context) ([]domain.Event, error)
}