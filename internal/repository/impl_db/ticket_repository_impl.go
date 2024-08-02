package impldb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/domain/dto"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type TicketRepository struct {
	mtx sync.Mutex
	db *sql.DB
}

func NewTicketRepository(database *sql.DB) repository.TicketRepositoryInterface {
	return &TicketRepository{
		db: database,
	}
}

func (repo *TicketRepository) Save(ctx context.Context, ticket *domain.Ticket) (error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save ticket because of timeout with message: %s", ctx.Err()))
		return ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new ticket")
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return err
		}

		log.Trace().Msg("Checking if ticket with the same Event ID and Type already exist")

		query := `
			SELECT
				e.event_id,
				e.event_name,
				t.ticket_id,
				t.name,
				t.price,
				t.type
			FROM tickets t
			JOIN events e
			ON t.event_id=e.event_id
			WHERE t.event_id=$1`

		log.Trace().Msg("Query is set")
		
		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return err
		}

		defer stmt.Close()

		res, err := stmt.QueryContext(ctx, ticket.EventID)
		log.Trace().Msg("Query ran")

		if err != nil {
			return err
		}

		defer res.Close()

		var foundTicket dto.TicketResponseTicketOrder
		for res.Next() {
			log.Trace().Msg("Inside checking event and type")
			log.Debug().Msg(fmt.Sprintf("%+v", res.Next()))

			res.Scan(&foundTicket.EventDetails.EventID, &foundTicket.EventDetails.EventName, &foundTicket.TicketID, &foundTicket.Name, &foundTicket.Price, &foundTicket.Type)
			log.Debug().
				Int("Ticket ID: ", foundTicket.TicketID).
				Int("Event ID: ", foundTicket.EventDetails.EventID).
				Str("Event Name: ", foundTicket.EventDetails.EventName).
				Str("Name: ", foundTicket.Name).
				Float64("Price: ", foundTicket.Price).
				Str("Type: ", foundTicket.Type).
				Msg("")

			if foundTicket.EventDetails.EventID == ticket.EventID && foundTicket.Type == ticket.Type {
				log.Error().Msg(fmt.Sprintf("ticket for Event ID %d with Type %s already exist", ticket.EventID, ticket.Type))
				return fmt.Errorf("ticket with Event ID %d with Type %s already exist", ticket.EventID, ticket.Type)
			}
		}

		query2 := `INSERT INTO tickets (event_id, name, price, stock, type) VALUES ($1, $2, $3, $4, $5) RETURNING ticket_id`
		log.Trace().Msg("Query 2 is set")

		stmt2, err := trx.PrepareContext(ctx, query2)
		log.Trace().Msg("Prepared statement 2 created with context")

		if err != nil {
			return err
		}

		defer stmt2.Close()

		errScan := stmt2.QueryRowContext(ctx, ticket.EventID, ticket.Name, ticket.Price, ticket.Stock, ticket.Type).Scan(&ticket.TicketID)
		log.Trace().Msg("Query 2 ran")

		if errScan != nil {
			trx.Rollback()
			return errScan
		}

		if err = trx.Commit(); err != nil {
			return err
		}

		log.Info().Msg("New ticket saved")
		return nil
	}
}

func (repo *TicketRepository) FindByID(ctx context.Context, id int) (domain.Ticket, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository find by id")
	log.Debug().Msg(fmt.Sprintf("Ticket repo find by id receive id with value %d", id))

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket because of timeout with message: %s", ctx.Err()))
		return domain.Ticket{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket")
		
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return domain.Ticket{}, err
		}

		query := `
			SELECT
				t.ticket_id,
				t.event_id,
				t.name,
				t.price,
				t.stock,
				t.type
			FROM tickets t
			WHERE
				t.ticket_id=$1
		`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared Statement created with context")

		if err != nil {
			return domain.Ticket{}, err
		}

		defer stmt.Close()

		var ticket domain.Ticket

		errScan := stmt.QueryRowContext(ctx, id).Scan(
			&ticket.TicketID,
			&ticket.EventID,
			&ticket.Name,
			&ticket.Price,
			&ticket.Stock,
			&ticket.Type,
		)

		if errScan != nil {
			if errScan == sql.ErrNoRows {
				log.Error().Msg(fmt.Sprintf("Ticket with ID %d not found", id))
				return domain.Ticket{}, fmt.Errorf("no Ticket found with ID %d", id)
			}
			return domain.Ticket{}, err
		}

		if err = trx.Commit(); err != nil {
			return domain.Ticket{}, err
		}

		log.Info().Msg("Ticket repo find by id completed")

		return ticket, nil
	}
}

func (repo *TicketRepository) GetAll(ctx context.Context) ([]domain.Ticket, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository get all")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket because of timeout with message: %s", ctx.Err()))
		return []domain.Ticket{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket")

		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return []domain.Ticket{}, err
		}

		query := `SELECT ticket_id, event_id, name, price, type FROM tickets`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return []domain.Ticket{}, err
		}

		defer stmt.Close()

		res, err := stmt.QueryContext(ctx)
		log.Trace().Msg("Query ran")

		if err != nil {
			return []domain.Ticket{}, err
		}

		defer res.Close()

		var listOfTickets []domain.Ticket

		for res.Next() {
			var ticket domain.Ticket

			res.Scan(&ticket.TicketID, &ticket.EventID, &ticket.Name, &ticket.Price, &ticket.Type)

			listOfTickets = append(listOfTickets, ticket)
		}
		
		log.Info().Msg("Fetching completed")
		return listOfTickets, nil
	}
}

func (repo *TicketRepository) Deduct(ctx context.Context, id int, amount int) (domain.Ticket, error) {
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository deduct")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to deduct ticket because of timeout with message: %s", ctx.Err()))
		return domain.Ticket{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to deduct ticket")

		foundTicket, err := repo.FindByID(ctx, id)
		log.Debug().Msg(fmt.Sprintf("Ticket ID found with ID %d", foundTicket.TicketID))

		if err != nil {
			return domain.Ticket{}, err
		}

		repo.mtx.Lock()
		log.Debug().Msg(fmt.Sprintf("Stock before deduct: %d", foundTicket.Stock))
		foundTicket.Stock = foundTicket.Stock - amount

		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return domain.Ticket{}, err
		}

		query := `
			UPDATE
				tickets
			SET
				stock=$1
			WHERE
				ticket_id=$2
			RETURNING
				*
		`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return domain.Ticket{}, err
		}

		var ticket domain.Ticket

		errScan := stmt.QueryRowContext(ctx, foundTicket.Stock, id).Scan(
			&ticket.TicketID,
			&ticket.EventID,
			&ticket.Name,
			&ticket.Price,
			&ticket.Stock,
			&ticket.Type,
		)

		if errScan != nil {
			return domain.Ticket{}, errScan
		}

		log.Debug().Msg(fmt.Sprintf("Stock after deduct: %d", ticket.Stock))
		log.Info().Msg("Successfully deduct stock")

		return ticket, nil
	}
}

func (repo *TicketRepository) Restore(ctx context.Context, id int, amount int) {
	log.Trace().Msg("Inside ticket repository restore")
	defer repo.mtx.Unlock()

	repo.mtx.Lock()
	log.Info().Msg("Trying to restore stock")
	
	trx, err := repo.db.Begin()
	log.Trace().Msg("Begin transaction")

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Error trying to built transaction in restore stock with message: %s", err.Error()))
	}

	query := `
		SELECT
			t.ticket_id,
			t.event_id,
			t.name,
			t.price,
			t.stock,
			t.type
		FROM tickets t
		WHERE
			t.ticket_id=$1
	`
	log.Trace().Msg("Query is set")

	stmt, err := trx.Prepare(query)
	log.Trace().Msg("Prepared statement created")

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Error trying to create prepared statement in restore stock with message: %s", err.Error()))
	}

	defer stmt.Close()

	var foundTicket domain.Ticket

	errScan := stmt.QueryRow(id).Scan(
		&foundTicket.TicketID,
		&foundTicket.EventID,
		&foundTicket.Name,
		&foundTicket.Price,
		&foundTicket.Stock,
		&foundTicket.Type,
	)

	if errScan != nil {
		if errScan == sql.ErrNoRows {
			log.Error().Msg(fmt.Sprintf("Ticket with ID %d not found", id))
			log.Fatal().Msg(fmt.Sprintf("no Ticket found with ID %d", id))
		}
		log.Fatal().Msg(fmt.Sprintf("Error when trying to scan find by id ticket with message: %s", errScan.Error()))
	}

	log.Debug().Msg(fmt.Sprintf("Stock before restore: %d", foundTicket.Stock))
	foundTicket.Stock = foundTicket.Stock + amount

	query2 := `
		UPDATE
			tickets
		SET
			stock=$1
		WHERE
			ticket_id=$2
		RETURNING
			ticket_id,
			stock
	`
	log.Trace().Msg("Query update is set")

	stmt2, err := trx.Prepare(query2)

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Error trying to create prepared statement in restore stock with message: %s", err.Error()))
	}

	errScan2 := stmt2.QueryRow(foundTicket.Stock, id).Scan(&foundTicket.TicketID, &foundTicket.Stock)

	if errScan2 != nil {
		log.Fatal().Msg(fmt.Sprintf("Error when trying to scan update ticket with message: %s", errScan2.Error()))
	}

	log.Debug().Msg(fmt.Sprintf("Stock after restore: %d", foundTicket.Stock))
	log.Info().Msg("Stock restored successfully")
}