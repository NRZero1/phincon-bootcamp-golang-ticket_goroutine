package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type EventRepositoryInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(context context.Context, event *domain.Event) (error)
}

type EventFindById interface {
	FindByID(context context.Context, id int) (domain.Event, error)
}

type EventGetAll interface {
	GetAll(context context.Context) ([]domain.Event, error)
}