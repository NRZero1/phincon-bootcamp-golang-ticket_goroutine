package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticket_goroutine/internal/middleware"
	"ticket_goroutine/internal/provider/db"
	"ticket_goroutine/internal/provider/handler"
	"ticket_goroutine/internal/provider/repository"
	"ticket_goroutine/internal/provider/routes"
	"ticket_goroutine/internal/provider/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initDB() *sql.DB {
	db, err := db.NewConnection("postgresql").GetConnection("postgres", "root", "localhost", "5432", "phincon_golang_ticket")

	if err != nil {
		return nil
	}

	return db
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	router := gin.New()

	// middlewares := middleware.CreateStack(
	// 	middleware.Logging,
	// )

	db := initDB()

	repository.InitRepository(db)
	usecase.InitUseCase()
	handler.InitHandler()

	router.Use(gin.Recovery(), middleware.Logging())
	
	globalGroup := router.Group("")
	{
		routes.EventRoutes(globalGroup.Group("/event"), handler.EventHandler)
		routes.UserRoutes(globalGroup.Group("/user"), handler.UserHandler)
		routes.TicketRoutes(globalGroup.Group("/ticket"), handler.TicketHandler)
		routes.TicketOrderRoutes(globalGroup.Group("/ticket-order"), handler.TicketOrderHandler)
	}

	db.SetMaxOpenConns(25)
	
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Server is running in port ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msg(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("Shutting down the server...")
	errDbClose := db.Close()

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