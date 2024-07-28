package main

import (
	"fmt"
	"net/http"
	"ticket_goroutine/internal/handler"
	handlerImplement "ticket_goroutine/internal/handler/impl"
	"ticket_goroutine/internal/middleware"
	"ticket_goroutine/internal/repository"
	repoImplement "ticket_goroutine/internal/repository/impl"
	"ticket_goroutine/internal/usecase"
	useCaseImplement "ticket_goroutine/internal/usecase/impl"

	"github.com/rs/zerolog/log"
)

var eventRepo repository.EventRepositoryInterface
var userRepo repository.UserRepositoryInterface
var ticketRepo repository.TicketRepositoryInterface
var ticketOrderRepo repository.TicketOrderRepositoryInterface

var eventUseCase usecase.EventUseCaseInterface
var userUseCase usecase.UserUseCaseInterface
var ticketUseCase usecase.TicketUseCaseInterface
var ticketOrderUseCase usecase.TicketOrderUseCaseInterface

var eventHandler handler.EventHandlerInterface
var userHandler handler.UserHandlerInterface
var ticketHandler handler.TicketHandlerInterface
var ticketOrderHandler handler.TicketOrderHandlerInterface

func initEventHandler() handler.EventHandlerInterface {
	eventRepo = repoImplement.NewEventRepository()
	eventUseCase = useCaseImplement.NewEventUseCase(eventRepo)
	handler := handlerImplement.NewEventHandler(eventUseCase)
	return handler
}

func initUserHandler() handler.UserHandlerInterface {
	userRepo = repoImplement.NewUserRepository()
	userUseCase = useCaseImplement.NewUserUseCase(userRepo)
	handler := handlerImplement.NewUserHandler(userUseCase)
	return handler
}

func initTicketHandler() handler.TicketHandlerInterface {
	ticketRepo = repoImplement.NewTicketRepository()

	if eventUseCase == nil {
        log.Fatal().Msg("Event use case is not initialized")
    }

	ticketUseCase = useCaseImplement.NewTicketUseCase(ticketRepo, eventUseCase)

	handler := handlerImplement.NewTicketHandler(ticketUseCase)
	return handler
}

func initTicketOrderHandler() handler.TicketOrderHandlerInterface {
	ticketOrderRepo = repoImplement.NewTicketOrderRepository()
	ticketOrderUseCase = useCaseImplement.NewTicketOrderUseCase(ticketOrderRepo, ticketUseCase, eventUseCase, userUseCase)
	handler := handlerImplement.NewTicketOrderHandler(ticketOrderUseCase)
	
	return handler
}

func init() {
	eventHandler = initEventHandler()
	userHandler = initUserHandler()
	ticketHandler = initTicketHandler()
	ticketOrderHandler = initTicketOrderHandler()
}

func main() {
	router := middleware.NewRouter()

	router.AddRoute("GET", "/event/", eventHandler.GetAll)
	router.AddRoute("GET", "/event/{id}", eventHandler.FindById)
	router.AddRoute("POST", "/event/", eventHandler.Save)

	router.AddRoute("GET", "/user/", userHandler.GetAll)
	router.AddRoute("GET", "/user/{id}", userHandler.FindById)
	router.AddRoute("POST", "/user/", userHandler.Save)

	router.AddRoute("GET", "/ticket/", ticketHandler.GetAll)
	router.AddRoute("GET", "/ticket/{id}", ticketHandler.FindById)
	router.AddRoute("POST", "/ticket/", ticketHandler.Save)

	router.AddRoute("GET", "/ticket-order/", ticketOrderHandler.GetAll)
	router.AddRoute("GET", "/ticket-order/{id}", ticketOrderHandler.FindById)
	router.AddRoute("POST", "/ticket-order/", ticketOrderHandler.Save)

	middlewares := middleware.CreateStack(
		middleware.Logging,
	)
	
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middlewares(router),
	}

	fmt.Println("Server is running in port ", server.Addr)
	err := server.ListenAndServe()

	if err != nil {
		panic(err.Error())
	}
}