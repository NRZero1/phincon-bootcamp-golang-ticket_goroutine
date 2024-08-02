package routes

import (
	"ticket_goroutine/internal/handler"

	"github.com/gin-gonic/gin"
)

func EventRoutes(routerGroup *gin.RouterGroup, eventHandler handler.EventHandlerInterface) {
	routerGroup.GET("/", eventHandler.GetAll)
	routerGroup.GET("/:id", eventHandler.FindById)
	routerGroup.POST("/", eventHandler.Save)
}