package usecase

import "ticket_goroutine/internal/domain"

type EventUseCaseInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(event domain.Event) (domain.Event, error)
}

type EventFindById interface {
	FindById(id int) (domain.Event, error)
}

type EventGetAll interface {
	GetAll() ([]domain.Event)
}