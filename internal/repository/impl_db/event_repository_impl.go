package impldb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type EventRepository struct {
	mtx sync.Mutex
	db *sql.DB
}

func NewEventRepository(database *sql.DB) repository.EventRepositoryInterface {
	return &EventRepository{
		db: database,
	}
}

func (repo *EventRepository) Save(ctx context.Context, event *domain.Event) (error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside event repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save event because of timeout with message: %s", ctx.Err()))
		return ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new event")

		trx, err := repo.db.BeginTx(ctx, nil)

		if err != nil {
			return err
		}

		query := `INSERT INTO events (event_name) VALUES ($1) RETURNING event_id`

		stmt, err := trx.PrepareContext(ctx, query)

		if err != nil {
			return err
		}

		errScan := stmt.QueryRowContext(ctx, event.EventName).Scan(&event.EventID)

		if errScan != nil {
			trx.Rollback()
			return err
		}

		if err = trx.Commit(); err != nil {
			return err
		}
		
		log.Info().Msg("New event saved")
		return nil
	}
}

func (repo *EventRepository) FindByID(ctx context.Context, id int) (domain.Event, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside event repository find by id")
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch event because of timeout with message: %s", ctx.Err()))
		return domain.Event{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch event")
	
		trx, err := repo.db.BeginTx(ctx, nil)

		if err != nil {
			return domain.Event{}, err
		}

		query := "SELECT event_id, event_name FROM events WHERE event_id=$1"

		stmt, err := trx.PrepareContext(ctx, query)

		if err != nil {
			return domain.Event{}, err
		}

		var event domain.Event
		errScan := stmt.QueryRowContext(ctx, id).Scan(
			&event.EventID,
			&event.EventName,
		)

		if errScan != nil {
			if errScan == sql.ErrNoRows {
				log.Error().Msg(fmt.Sprintf("Event with ID %d not found", id))
				return domain.Event{}, fmt.Errorf("no Event found with ID %d", id)
			}
			return event, nil
		}

		if err = trx.Commit(); err != nil {
			return domain.Event{}, err
		}

		return event, nil
	}
}

func (repo *EventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside event repository get all")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch event because of timeout with message: %s", ctx.Err()))
		return []domain.Event{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch event")
		
		trx, err := repo.db.BeginTx(ctx, nil)

		if err != nil {
			return []domain.Event{}, err
		}

		query := "SELECT * FROM events"

		stmt, err := trx.PrepareContext(ctx, query)

		if err != nil {
			return []domain.Event{}, err
		}

		res, err := stmt.QueryContext(ctx)

		if err != nil {
			return []domain.Event{}, err
		}

		var listOfEvents []domain.Event

		for res.Next() {
			var event domain.Event
			res.Scan(&event.EventID, &event.EventName)

			listOfEvents = append(listOfEvents, event)
		}

		log.Info().Msg("Fetching completed")
		return listOfEvents, nil
	}
}