package main

import (
	"fmt"
	"net/http"
	"ticket_goroutine/internal/middleware"
)

func main() {
	router := middleware.NewRouter()
	// r := repository.NewBookRepository()
	// uc := usecase.NewBookUseCase(r)
	// h := handler.NewBookHandler(uc)

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