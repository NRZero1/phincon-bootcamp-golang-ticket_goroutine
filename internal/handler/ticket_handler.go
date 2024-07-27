package handler

import (
	"net/http"
)

type TicketHandlerInterface interface {
	TicketSave
	TicketFindById
	TicketGetAll
}

type TicketSave interface {
	Save(responseWritter http.ResponseWriter, request *http.Request)
}

type TicketFindById interface {
	FindById(responseWritter http.ResponseWriter, request *http.Request)
}

type TicketGetAll interface {
	GetAll(responseWritter http.ResponseWriter, request *http.Request)
}