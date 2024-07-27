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
	return &TicketRepository {
		tickets: map[int]domain.Ticket{},
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