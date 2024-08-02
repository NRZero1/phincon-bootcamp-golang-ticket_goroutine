package handler

import (
	"github.com/gin-gonic/gin"
)

type EventHandlerInterface interface {
	EventSave
	EventFindById
	EventGetAll
}

type EventSave interface {
	Save(context *gin.Context)
}

type EventFindById interface {
	FindById(context *gin.Context)
}

type EventGetAll interface {
	GetAll(context *gin.Context)
}