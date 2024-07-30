package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticket_goroutine/internal/handler"
	handlerImplement "ticket_goroutine/internal/handler/impl"
	"ticket_goroutine/internal/middleware"
	"ticket_goroutine/internal/provider/db"
	"ticket_goroutine/internal/repository"
	repoImplement "ticket_goroutine/internal/repository/impl_db"
	"ticket_goroutine/internal/usecase"
	useCaseImplement "ticket_goroutine/internal/usecase/impl"
	"time"

	"github.com/rs/zerolog/log"
)

var eventRepo repository.EventRepositoryInterface
// var userRepo repository.UserRepositoryInterface
// var ticketRepo repository.TicketRepositoryInterface
// var ticketOrderRepo repository.TicketOrderRepositoryInterface

var eventUseCase usecase.EventUseCaseInterface
var userUseCase usecase.UserUseCaseInterface
// var ticketUseCase usecase.TicketUseCaseInterface
// var ticketOrderUseCase usecase.TicketOrderUseCaseInterface

var eventHandler handler.EventHandlerInterface
// var userHandler handler.UserHandlerInterface
// var ticketHandler handler.TicketHandlerInterface
// var ticketOrderHandler handler.TicketOrderHandlerInterface

var database *sql.DB

func initDB() *sql.DB {
	db, err := db.NewConnection("postgresql").GetConnection("postgres", "root", "localhost", "5432", "phincon_golang_ticket")

	if err != nil {
		return nil
	}

	return db
}

func initEventHandler() handler.EventHandlerInterface {
	eventRepo = repoImplement.NewEventRepository(database)

	if eventRepo == nil {
		log.Fatal().Msg("Event repo didn't get initialized")
	}
	eventUseCase = useCaseImplement.NewEventUseCase(eventRepo)
	handler := handlerImplement.NewEventHandler(eventUseCase)
	return handler
}

// func initUserHandler() handler.UserHandlerInterface {
// 	userRepo = repoImplement.NewUserRepository()
// 	userUseCase = useCaseImplement.NewUserUseCase(userRepo)
// 	handler := handlerImplement.NewUserHandler(userUseCase)
// 	return handler
// }

// func initTicketHandler() handler.TicketHandlerInterface {
// 	ticketRepo = repoImplement.NewTicketRepository()

// 	if eventUseCase == nil {
//         log.Fatal().Msg("Event use case is not initialized")
//     }

// 	ticketUseCase = useCaseImplement.NewTicketUseCase(ticketRepo, eventUseCase)

// 	handler := handlerImplement.NewTicketHandler(ticketUseCase)
// 	return handler
// }

// func initTicketOrderHandler() handler.TicketOrderHandlerInterface {
// 	ticketOrderRepo = repoImplement.NewTicketOrderRepository()
// 	ticketOrderUseCase = useCaseImplement.NewTicketOrderUseCase(ticketOrderRepo, ticketUseCase, eventUseCase, userUseCase)
// 	handler := handlerImplement.NewTicketOrderHandler(ticketOrderUseCase)
	
// 	return handler
// }

func init() {
	database = initDB()
	eventHandler = initEventHandler()
	// userHandler = initUserHandler()
	// ticketHandler = initTicketHandler()
	// ticketOrderHandler = initTicketOrderHandler()
}

func main() {
	router := middleware.NewRouter()

	router.AddRoute("GET", "/event/", eventHandler.GetAll)
	router.AddRoute("GET", "/event/{id}", eventHandler.FindById)
	router.AddRoute("POST", "/event/", eventHandler.Save)

	// router.AddRoute("GET", "/user/", userHandler.GetAll)
	// router.AddRoute("GET", "/user/{id}", userHandler.FindById)
	// router.AddRoute("POST", "/user/", userHandler.Save)

	// router.AddRoute("GET", "/ticket/", ticketHandler.GetAll)
	// router.AddRoute("GET", "/ticket/{id}", ticketHandler.FindById)
	// router.AddRoute("POST", "/ticket/", ticketHandler.Save)

	// router.AddRoute("GET", "/ticket-order/", ticketOrderHandler.GetAll)
	// router.AddRoute("GET", "/ticket-order/{id}", ticketOrderHandler.FindById)
	// router.AddRoute("POST", "/ticket-order/", ticketOrderHandler.Save)

	middlewares := middleware.CreateStack(
		middleware.Logging,
	)
	
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middlewares(router),
	}

	go func() {
		fmt.Println("Server is running in port ", server.Addr)
		err := server.ListenAndServe()

		if err != nil {
			panic(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("Shutting down the server...")
	errDbClose := database.Close()

	if errDbClose != nil {
		log.Fatal().Msg(fmt.Sprintf("Database shutdown error: %s", errDbClose.Error()))
	}

    // Set a timeout for shutdown (for example, 5 seconds).
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal().Msg(fmt.Sprintf("Server shutdown error: %v", err))
    }
    log.Info().Msg("Server gracefully stopped")
}