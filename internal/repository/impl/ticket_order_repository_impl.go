package impl

import (
	"context"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type TicketOrderRepository struct {
	mtx sync.Mutex
	ticket_orders map[int]domain.TicketOrder
	nextId int
}

func NewTicketOrderRepository() repository.TicketOrderRepositoryInterface {
	return &TicketOrderRepository {
		ticket_orders: map[int]domain.TicketOrder{},
		nextId: 1,
	}
}

func (repo *TicketOrderRepository) Save(ctx context.Context, ticketOrder *domain.TicketOrder) (domain.TicketOrder, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket order repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save ticket order because of timeout with message: %s", ctx.Err()))
		return domain.TicketOrder{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new ticket order")
		if _, exists := repo.ticket_orders[ticketOrder.OrderID]; exists {
			log.Error().Msg(fmt.Sprintf("ticket order with ID %d already exist", ticketOrder.OrderID))
			return domain.TicketOrder{}, fmt.Errorf("ticket order with ID %d already exist", ticketOrder.OrderID)
		}

		ticketOrder.OrderID = repo.nextId
		repo.ticket_orders[ticketOrder.OrderID] = *ticketOrder
		repo.nextId++
		log.Info().Msg("New ticket order saved")
		return repo.ticket_orders[ticketOrder.OrderID], nil
	}
}

func (repo *TicketOrderRepository) FindByID(ctx context.Context, id int) (domain.TicketOrder, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket order repository find by id")
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket order because of timeout with message: %s", ctx.Err()))
		return domain.TicketOrder{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket order")
		if foundTicketOrder, exists := repo.ticket_orders[id]; exists {
			log.Info().Msg("Fetching completed")
			return foundTicketOrder, nil
		}
		log.Error().Msg(fmt.Sprintf("Ticket order with ID %d not found", id))
		return domain.TicketOrder{}, fmt.Errorf("no Ticket order found with ID %d", id)
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
		listOfTicketOrders := make([]domain.TicketOrder, 0, len(repo.ticket_orders))

		for _, v := range repo.ticket_orders {
			listOfTicketOrders = append(listOfTicketOrders, v)
		}

		log.Info().Msg("Fetching completed")
		return listOfTicketOrders, nil
	}
}