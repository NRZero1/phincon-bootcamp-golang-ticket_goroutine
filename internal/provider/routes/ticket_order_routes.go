package routes

import (
	"ticket_goroutine/internal/handler"

	"github.com/gin-gonic/gin"
)

func TicketOrderRoutes(routerGroup *gin.RouterGroup, ticketOrderHandler handler.TicketOrderHandlerInterface) {
	routerGroup.GET("/", ticketOrderHandler.GetAll)
	routerGroup.GET("/:id", ticketOrderHandler.FindById)
	routerGroup.POST("/", ticketOrderHandler.Save)
}