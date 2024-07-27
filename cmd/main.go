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
var userHandler handler.UserHandlerInterface

func initEventHandler() handler.EventHandlerInterface {
	repo := repoImplement.NewEventRepository()
	usecase := useCaseImplement.NewEventUseCase(repo)
	handler := handlerImplement.NewEventHandler(usecase)
	return handler
}

func initUserHandler() handler.UserHandlerInterface {
	repo := repoImplement.NewUserRepository()
	usecase := useCaseImplement.NewUserUseCase(repo)
	handler := handlerImplement.NewUserHandler(usecase)
	return handler
}

func init() {
	eventHandler = initEventHandler()
	userHandler = initUserHandler()
}

func main() {
	router := middleware.NewRouter()

	router.AddRoute("GET", "/event/", eventHandler.GetAll)
	router.AddRoute("GET", "/event/{id}", eventHandler.FindById)
	router.AddRoute("POST", "/event/", eventHandler.Save)

	router.AddRoute("GET", "/user/", userHandler.GetAll)
	router.AddRoute("GET", "/user/{id}", userHandler.FindById)
	router.AddRoute("POST", "/user/", userHandler.Save)

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