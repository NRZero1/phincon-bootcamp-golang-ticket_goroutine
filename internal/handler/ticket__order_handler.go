package handler

import (
	"net/http"
)

type TicketOrderHandlerInterface interface {
	TicketOrderSave
	TicketOrderFindById
	TicketOrderGetAll
}

type TicketOrderSave interface {
	Save(responseWritter http.ResponseWriter, request *http.Request)
}

type TicketOrderFindById interface {
	FindById(responseWritter http.ResponseWriter, request *http.Request)
}

type TicketOrderGetAll interface {
	GetAll(responseWritter http.ResponseWriter, request *http.Request)
}