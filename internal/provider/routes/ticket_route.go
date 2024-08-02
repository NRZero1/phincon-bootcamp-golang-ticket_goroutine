package routes

import (
	"ticket_goroutine/internal/handler"

	"github.com/gin-gonic/gin"
)

func TicketRoutes(routerGroup *gin.RouterGroup, ticketHandler handler.TicketHandlerInterface) {
	routerGroup.GET("/", ticketHandler.GetAll)
	routerGroup.GET("/:id", ticketHandler.FindById)
	routerGroup.POST("/", ticketHandler.Save)
}