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

type TicketOrderRepository struct {
	mtx sync.Mutex
	db *sql.DB
}

func NewTicketOrderRepository(database *sql.DB) repository.TicketOrderRepositoryInterface {
	return &TicketOrderRepository {
		db: database,
	}
}

func (repo *TicketOrderRepository) Save(ctx context.Context, ticketOrder *domain.TicketOrder) (error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket order repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save ticket order because of timeout with message: %s", ctx.Err()))
		return ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new ticket order")
		
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return err
		}

		query := `
			INSERT INTO
				ticket_order
				(ticket_id, user_id, amount, total_price)
			VALUES
			($1, $2, $3, $4)
			RETURNING order_id
		`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return err
		}

		defer stmt.Close()

		errScan := stmt.QueryRowContext(ctx, ticketOrder.TicketID, ticketOrder.UserID, ticketOrder.Amount, ticketOrder.TotalPrice).Scan(&ticketOrder.OrderID)
		log.Trace().Msg("Query ran")

		if errScan != nil {
			trx.Rollback()
			return errScan
		}

		if err = trx.Commit(); err != nil {
			return err
		}

		log.Info().Msg("New ticket order saved")
		return nil
	}
}

func (repo *TicketOrderRepository) FindByID(ctx context.Context, id int) (domain.TicketOrder, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket order repository find by id")
	log.Debug().Msg(fmt.Sprintf("Ticket order repo find by id received id value is %d", id))

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket order because of timeout with message: %s", ctx.Err()))
		return domain.TicketOrder{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket order")
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return domain.TicketOrder{}, err
		}

		query := `
			SELECT
				t.order_id,
				t.ticket_id,
				t.user_id,
				t.amount,
				t.total_price
			FROM ticket_order t
			WHERE
				t.order_id=$1
		`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return domain.TicketOrder{}, err
		}

		defer stmt.Close()

		var ticketOrder domain.TicketOrder

		errScan := stmt.QueryRowContext(ctx, id).Scan(
			&ticketOrder.OrderID,
			&ticketOrder.TicketID,
			&ticketOrder.UserID,
			&ticketOrder.Amount,
			&ticketOrder.TotalPrice,
		)

		if errScan != nil {
			if errScan == sql.ErrNoRows {
				log.Error().Msg(fmt.Sprintf("Ticket order with ID %d not found", id))
				return domain.TicketOrder{}, fmt.Errorf("no Ticket order found with ID %d", id)
			}
			return domain.TicketOrder{}, err
		}

		if err = trx.Commit(); err != nil {
			return domain.TicketOrder{}, err
		}

		log.Info().Msg("Ticket order repo find by id completed")
		return ticketOrder, nil
	}
}

func (repo *TicketOrderRepository) GetAll(ctx context.Context) ([]domain.TicketOrder, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket order repository get all")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket order because of timeout with message: %s", ctx.Err()))
		return []domain.TicketOrder{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket order")
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return []domain.TicketOrder{}, err
		}

		query := `SELECT * FROM ticket_order`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return []domain.TicketOrder{}, err
		}

		defer stmt.Close()

		res, err := stmt.QueryContext(ctx)
		log.Trace().Msg("Query ran")

		if err != nil {
			return []domain.TicketOrder{}, err
		}

		defer res.Close()

		var listOfTicketOrders []domain.TicketOrder

		for res.Next() {
			var ticketOrder domain.TicketOrder

			res.Scan(&ticketOrder.TicketID, &ticketOrder.OrderID, &ticketOrder.UserID, &ticketOrder.Amount, &ticketOrder.TotalPrice)

			listOfTicketOrders = append(listOfTicketOrders, ticketOrder)
		}
		
		log.Info().Msg("Fetching completed")
		return listOfTicketOrders, nil
	}
}