package implmap

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"ticket_goroutine/internal/domain"
// 	"ticket_goroutine/internal/repository"

// 	"github.com/rs/zerolog/log"
// )

// type EventRepository struct {
// 	mtx sync.Mutex
// 	events map[int]domain.Event
// }

// func NewEventRepository() repository.EventRepositoryInterface {
// 	repo := &EventRepository {
// 		events: map[int]domain.Event{},
// 	}

// 	repo.initData()
// 	return repo
// }

// func (repo *EventRepository) initData() {
// 	repo.mtx.Lock()
//     defer repo.mtx.Unlock()

// 	repo.events[1] = domain.Event{
// 		EventID: 1,
// 		EventName: "Test Event 1",
// 	}

// 	repo.events[2] = domain.Event{
// 		EventID: 2,
// 		EventName: "Test Event 2",
// 	}
// }

// func (repo *EventRepository) Save(ctx context.Context, event *domain.Event) (domain.Event, error) {
// 	repo.mtx.Lock()
// 	defer repo.mtx.Unlock()

// 	log.Trace().Msg("Inside event repository save")

// 	select {
// 	case <- ctx.Done():
// 		log.Error().Msg(fmt.Sprintf("Error when trying to save event because of timeout with message: %s", ctx.Err()))
// 		return domain.Event{}, ctx.Err()
// 	default:
// 		log.Trace().Msg("Attempting to save new event")
// 		if _, exists := repo.events[event.EventID]; exists {
// 			log.Error().Msg(fmt.Sprintf("Event with ID %d already exist", event.EventID))
// 			return domain.Event{}, fmt.Errorf("event with ID %d already exist", event.EventID)
// 		}

// 		repo.events[event.EventID] = *event
// 		log.Info().Msg("New event saved")
// 		return repo.events[event.EventID], nil
// 	}
// }

// func (repo *EventRepository) FindByID(ctx context.Context, id int) (domain.Event, error) {
// 	repo.mtx.Lock()
// 	defer repo.mtx.Unlock()

// 	log.Trace().Msg("Inside event repository find by id")

// 	select {
// 	case <- ctx.Done():
// 		log.Error().Msg(fmt.Sprintf("Error when trying to fetch event because of timeout with message: %s", ctx.Err()))
// 		return domain.Event{}, ctx.Err()
// 	default:
// 		log.Trace().Msg("Attempting to fetch event")
// 		if foundEvent, exists := repo.events[id]; exists {
// 			log.Info().Msg("Fetching completed")
// 			return foundEvent, nil
// 		}
// 		log.Error().Msg(fmt.Sprintf("Event with ID %d not found", id))
// 		return domain.Event{}, fmt.Errorf("no Event found with ID %d", id)
// 	}
// }

// func (repo *EventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
// 	repo.mtx.Lock()
// 	defer repo.mtx.Unlock()

// 	log.Trace().Msg("Inside event repository get all")

// 	select {
// 	case <- ctx.Done():
// 		log.Error().Msg(fmt.Sprintf("Error when trying to fetch event because of timeout with message: %s", ctx.Err()))
// 		return []domain.Event{}, ctx.Err()
// 	default:
// 		log.Trace().Msg("Attempting to fetch event")
// 		listOfEvents := make([]domain.Event, 0, len(repo.events))

// 		for _, v := range repo.events {
// 			listOfEvents = append(listOfEvents, v)
// 		}

// 		log.Info().Msg("Fetching completed")
// 		return listOfEvents, nil
// 	}
// }