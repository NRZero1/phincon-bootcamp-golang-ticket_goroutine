package handler

import (
	"github.com/gin-gonic/gin"
)

type TicketHandlerInterface interface {
	TicketSave
	TicketFindById
	TicketGetAll
}

type TicketSave interface {
	Save(context *gin.Context)
}

type TicketFindById interface {
	FindById(context *gin.Context)
}

type TicketGetAll interface {
	GetAll(context *gin.Context)
}