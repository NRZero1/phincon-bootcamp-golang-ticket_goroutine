package impl

import (
	"context"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type TicketRepository struct {
	mtx sync.Mutex
	tickets map[int]domain.Ticket
}

func NewTicketRepository() repository.TicketRepositoryInterface {
	repo := &TicketRepository {
		tickets: map[int]domain.Ticket{},
	}

	repo.initData()
	return repo
}

func (repo *TicketRepository) initData() {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.tickets[1] = domain.Ticket{
		TicketID: 1,
		EventID: 1,
		Name: "Test Ticket 1",
		Price: 5000,
		Stock: 10,
		Type: "VIP",
	}

	repo.tickets[2] = domain.Ticket{
		TicketID: 2,
		EventID: 1,
		Name: "Test Ticket 2",
		Price: 250,
		Stock: 100,
		Type: "CAT 1",
	}
}

func (repo *TicketRepository) Save(ctx context.Context, ticket *domain.Ticket) (domain.Ticket, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save ticket because of timeout with message: %s", ctx.Err()))
		return domain.Ticket{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new ticket")
		if _, exists := repo.tickets[ticket.TicketID]; exists {
			log.Error().Msg(fmt.Sprintf("ticket with ID %d already exist", ticket.TicketID))
			return domain.Ticket{}, fmt.Errorf("ticket with ID %d already exist", ticket.TicketID)
		}

		log.Trace().Msg("Checking if ticket with the same Event ID and Type already exist")
		for _, v := range repo.tickets {
			if foundTicket, exists := repo.tickets[v.TicketID]; exists {
				if foundTicket.EventID == ticket.EventID && foundTicket.Type == ticket.Type {
					log.Error().Msg(fmt.Sprintf("ticket for Event ID %d with Type %s already exist", ticket.TicketID, ticket.Type))
					return domain.Ticket{}, fmt.Errorf("ticket with ID %d with Type %s already exist", ticket.TicketID, ticket.Type)
				}
			}
		}

		repo.tickets[ticket.TicketID] = *ticket
		log.Info().Msg("New ticket saved")
		return repo.tickets[ticket.TicketID], nil
	}
}

func (repo *TicketRepository) FindByID(ctx context.Context, id int) (domain.Ticket, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside ticket repository find by id")
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch ticket because of timeout with message: %s", ctx.Err()))
		return domain.Ticket{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch ticket")
		if foundTicket, exists := repo.tickets[id]; exists {
			log.Info().Msg("Fetching completed")
			return foundTicket, nil
		}
		log.Error().Msg(fmt.Sprintf("Ticket with ID %d not found", id))
		return domain.Ticket{}, fmt.Errorf("no Ticket found with ID %d", id)
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
		listOfTickets := make([]domain.Ticket, 0, len(repo.tickets))

		for _, v := range repo.tickets {
			listOfTickets = append(listOfTickets, v)
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

		repo.tickets[foundTicket.TicketID] = foundTicket

		log.Debug().Msg(fmt.Sprintf("Stock after deduct: %d", foundTicket.Stock))
		log.Info().Msg("Successfully deduct stock")

		return repo.tickets[foundTicket.TicketID], nil
	}
}

func (repo *TicketRepository) Restore(ctx context.Context, id int, amount int) {
	log.Trace().Msg("Inside ticket repository restore")
	defer repo.mtx.Unlock()

	repo.mtx.Lock()
	log.Info().Msg("Trying to restore stock")
	foundTicket, exist := repo.tickets[id]

	if exist {
		log.Debug().Msg(fmt.Sprintf("Stock before restore: %d", foundTicket.Stock))
		foundTicket.Stock = foundTicket.Stock + amount
		repo.tickets[id] = foundTicket
		log.Debug().Msg(fmt.Sprintf("Stock after restore: %d", foundTicket.Stock))
		log.Info().Msg("Stock restored successfully")
	}
}