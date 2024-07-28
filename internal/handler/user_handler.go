package handler

import (
	"net/http"
)

type UserHandlerInterface interface {
	UserSave
	UserFindById
	UserGetAll
}

type UserSave interface {
	Save(responseWritter http.ResponseWriter, request *http.Request)
}

type UserFindById interface {
	FindById(responseWritter http.ResponseWriter, request *http.Request)
}

type UserGetAll interface {
	GetAll(responseWritter http.ResponseWriter, request *http.Request)
}