package routes

import (
	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/controllers"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(router *gin.Engine, wsController *controllers.WebSocketController, cfg *config.Config) {
	ws := router.Group("/ws")
	ws.Use(middleware.AuthMiddleware(cfg))
	{
		// WebSocket connection endpoint
		ws.GET("/connect", wsController.Connect)

		// Notification endpoints
		ws.GET("/notifications", wsController.GetNotifications)
		ws.PUT("/notifications/:id/read", wsController.MarkNotificationAsRead)
		ws.PUT("/notifications/read-all", wsController.MarkAllNotificationsAsRead)

		// Chat endpoints
		ws.POST("/chat/rooms", wsController.CreateChatRoom)
		ws.POST("/chat/rooms/:room_id/messages", wsController.SendChatMessage)
		ws.GET("/chat/rooms/:room_id/messages", wsController.GetChatMessages)
	}
}
