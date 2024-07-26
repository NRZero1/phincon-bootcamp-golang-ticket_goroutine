package repository

import "ticket_goroutine/internal/domain"

type EventRepositoryInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(event *domain.Event) (domain.Event, error)
}

type EventFindById interface {
	FindByID(id int) (domain.Event, error)
}

type EventGetAll interface {
	GetAll() ([]domain.Event)
}