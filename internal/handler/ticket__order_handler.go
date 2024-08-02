package handler

import (
	"github.com/gin-gonic/gin"
)

type TicketOrderHandlerInterface interface {
	TicketOrderSave
	TicketOrderFindById
	TicketOrderGetAll
}

type TicketOrderSave interface {
	Save(context *gin.Context)
}

type TicketOrderFindById interface {
	FindById(context *gin.Context)
}

type TicketOrderGetAll interface {
	GetAll(context *gin.Context)
}