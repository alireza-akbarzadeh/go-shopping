package controllers

import (
	"net/http"

	"github.com/alireza-akbarzadeh/shopping-platform/config"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSController struct {
	hub *services.Hub
	cfg *config.Config
}

func NewWSController(hub *services.Hub, cfg *config.Config) *WSController {
	return &WSController{hub: hub, cfg: cfg}
}

func (wc *WSController) HandleWebSocket(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		utils.UnauthorizedResponse(c, "missing token")
		return
	}

	claims, err := utils.ValidateToken(tokenString, wc.cfg.JWT.Secret)
	if err != nil {
		utils.UnauthorizedResponse(c, "invalid token")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &services.Client{
		Conn:   conn,
		UserID: claims.UserID,
		Send:   make(chan []byte, 256),
		Hub:    wc.hub,
	}

	wc.hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
