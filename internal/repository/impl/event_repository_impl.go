package impl

import (
	"fmt"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"
)

type EventRepository struct {
	events map[int]domain.Event
}

func NewEventRepository() repository.EventRepositoryInterface {
	return EventRepository {
		events: map[int]domain.Event{},
	}
}

func (repo EventRepository) Save(event *domain.Event) (domain.Event, error) {
	if _, exists := repo.events[event.EventID]; exists {
		return domain.Event{}, fmt.Errorf("event with ID %d already exist", event.EventID)
	}

	repo.events[event.EventID] = *event
	return repo.events[event.EventID], nil
}

func (repo EventRepository) FindByID(id int) (domain.Event, error) {
	if foundEvent, exists := repo.events[id]; exists {
		return foundEvent, nil
	}

	return domain.Event{}, fmt.Errorf("No Event found with ID %d", id)
}

func (repo EventRepository) GetAll() ([]domain.Event) {
	listOfEvents := make([]domain.Event, 0, len(repo.events))

	for _, v := range repo.events {
		listOfEvents = append(listOfEvents, v)
	}

	return listOfEvents
}