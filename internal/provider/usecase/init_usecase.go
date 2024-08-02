package usecase

import (
	providerRepo "ticket_goroutine/internal/provider/repository"
	"ticket_goroutine/internal/usecase"
	useCaseImpl "ticket_goroutine/internal/usecase/impl"
)

var (
	EventUseCase usecase.EventUseCaseInterface
	UserUseCase usecase.UserUseCaseInterface
	TicketUseCase usecase.TicketUseCaseInterface
	TicketOrderUseCase usecase.TicketOrderUseCaseInterface
)

func InitUseCase() {
	EventUseCase = useCaseImpl.NewEventUseCase(providerRepo.EventRepository)
	UserUseCase = useCaseImpl.NewUserUseCase(providerRepo.UserRepository)
	TicketUseCase = useCaseImpl.NewTicketUseCase(providerRepo.TicketRepository, EventUseCase)
	TicketOrderUseCase = useCaseImpl.NewTicketOrderUseCase(providerRepo.TicketOrderRepository, TicketUseCase, EventUseCase, UserUseCase)
}