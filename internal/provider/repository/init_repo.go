package repository

import (
	"database/sql"
	"ticket_goroutine/internal/repository"
	repoImplement "ticket_goroutine/internal/repository/impl_db"
)

var (
	EventRepository repository.EventRepositoryInterface
	UserRepository repository.UserRepositoryInterface
	TicketRepository repository.TicketRepositoryInterface
	TicketOrderRepository repository.TicketOrderRepositoryInterface
)

func InitRepository(database *sql.DB) {
	EventRepository = repoImplement.NewEventRepository(database)
	UserRepository = repoImplement.NewUserRepository(database)
	TicketRepository = repoImplement.NewTicketRepository(database)
	TicketOrderRepository = repoImplement.NewTicketOrderRepository(database)
}