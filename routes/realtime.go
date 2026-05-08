package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRealtimeRoutes(router *gin.RouterGroup, wsCtrl *controllers.WSController) {
	router.GET("/ws", wsCtrl.HandleWebSocket)
}
