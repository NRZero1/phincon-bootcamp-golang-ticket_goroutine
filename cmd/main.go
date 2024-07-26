package main

import (
	"fmt"
	"net/http"
	"ticket_goroutine/internal/handler"
	handlerImplement "ticket_goroutine/internal/handler/impl"
	"ticket_goroutine/internal/middleware"
	repoImplement "ticket_goroutine/internal/repository/impl"
	useCaseImplement "ticket_goroutine/internal/usecase/impl"
)

var eventHandler handler.EventHandlerInterface

func initEventHandler() handler.EventHandlerInterface {
	repo := repoImplement.NewEventRepository()
	usecase := useCaseImplement.NewEventUseCase(repo)
	handler := handlerImplement.NewEventHandler(usecase)
	return handler
}

func init() {
	eventHandler = initEventHandler()
}

func main() {
	router := middleware.NewRouter()

	router.AddRoute("GET", "/event/", eventHandler.GetAll)
	router.AddRoute("GET", "/event/{id}", eventHandler.FindById)
	router.AddRoute("POST", "/event/", eventHandler.Save)

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