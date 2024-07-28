package handler

import (
	"net/http"
)

type EventHandlerInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(responseWritter http.ResponseWriter, request *http.Request)
}

type EventFindById interface {
	FindById(responseWritter http.ResponseWriter, request *http.Request)
}

type EventGetAll interface {
	GetAll(responseWritter http.ResponseWriter, request *http.Request)
}