package handler

import (
	"ticket_goroutine/internal/handler"
	handlerImpl "ticket_goroutine/internal/handler/impl_gin"
	providerUseCase "ticket_goroutine/internal/provider/usecase"
)

var (
	EventHandler handler.EventHandlerInterface
	UserHandler handler.UserHandlerInterface
	TicketHandler handler.TicketHandlerInterface
	TicketOrderHandler handler.TicketOrderHandlerInterface
)

func InitHandler() {
	providerUseCase.InitUseCase()
	EventHandler = handlerImpl.NewEventHandler(providerUseCase.EventUseCase)
	UserHandler = handlerImpl.NewUserHandler(providerUseCase.UserUseCase)
	TicketHandler = handlerImpl.NewTicketHandler(providerUseCase.TicketUseCase)
	TicketOrderHandler = handlerImpl.NewTicketOrderHandler(providerUseCase.TicketOrderUseCase)
}