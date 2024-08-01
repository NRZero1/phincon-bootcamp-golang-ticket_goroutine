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
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return err
		}

		query := `INSERT INTO events (event_name) VALUES ($1) RETURNING event_id`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return err
		}

		defer stmt.Close()

		errScan := stmt.QueryRowContext(ctx, event.EventName).Scan(&event.EventID)
		log.Trace().Msg("Query ran")

		if errScan != nil {
			trx.Rollback()
			return errScan
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
	log.Debug().Msg(fmt.Sprintf("Event repo find by id received id with value %d", id))
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch event because of timeout with message: %s", ctx.Err()))
		return domain.Event{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch event")
	
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin Transaction")

		if err != nil {
			return domain.Event{}, err
		}

		query := "SELECT event_id, event_name FROM events WHERE event_id=$1"
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return domain.Event{}, err
		}

		defer stmt.Close()

		var event domain.Event
		errScan := stmt.QueryRowContext(ctx, id).Scan(
			&event.EventID,
			&event.EventName,
		)
		log.Trace().Msg("Query ran")

		if errScan != nil {
			if errScan == sql.ErrNoRows {
				log.Error().Msg(fmt.Sprintf("Event with ID %d not found", id))
				return domain.Event{}, fmt.Errorf("no Event found with ID %d", id)
			}
			return domain.Event{}, errScan
		}

		if err = trx.Commit(); err != nil {
			return domain.Event{}, err
		}

		log.Info().Msg("Event repo find by id completed")
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
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return []domain.Event{}, err
		}

		query := "SELECT * FROM events"
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return []domain.Event{}, err
		}

		defer stmt.Close()

		res, err := stmt.QueryContext(ctx)
		log.Trace().Msg("Query ran")

		if err != nil {
			return []domain.Event{}, err
		}

		defer res.Close()

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