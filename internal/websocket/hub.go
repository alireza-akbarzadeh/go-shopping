package websocket

import (
	"sync"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
	"github.com/gorilla/websocket"
)

// WebSocket event types
const (
	EventOrderUpdate     = "order_update"
	EventPaymentSuccess  = "payment_success"
	EventInventoryUpdate = "inventory_update"
	EventChatMessage     = "chat_message"
	EventNotification    = "notification"
	EventUserOnline      = "user_online"
	EventUserOffline     = "user_offline"
)

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	UserID    uint        `json:"user_id,omitempty"`
	RoomID    string      `json:"room_id,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID     uint
	UserID uint
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte
	Rooms  map[string]bool // rooms this client is subscribed to
}

// Hub manages WebSocket connections and broadcasting
type Hub struct {
	clients     map[*Client]bool
	broadcast   chan []byte
	register    chan *Client
	unregister  chan *Client
	rooms       map[string]map[*Client]bool // room -> clients
	userClients map[uint]map[*Client]bool   // userID -> clients
	mu          sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		rooms:       make(map[string]map[*Client]bool),
		userClients: make(map[uint]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// addClient adds a client to the hub
func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// Add to user clients map
	if h.userClients[client.UserID] == nil {
		h.userClients[client.UserID] = make(map[*Client]bool)
	}
	h.userClients[client.UserID][client] = true

	utils.Log.WithField("user_id", client.UserID).Info("WebSocket client connected")
}

// removeClient removes a client from the hub
func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove from user clients map
		if userClients, exists := h.userClients[client.UserID]; exists {
			delete(userClients, client)
			if len(userClients) == 0 {
				delete(h.userClients, client.UserID)
			}
		}

		// Remove from rooms
		for roomID := range client.Rooms {
			if roomClients, exists := h.rooms[roomID]; exists {
				delete(roomClients, client)
				if len(roomClients) == 0 {
					delete(h.rooms, roomID)
				}
			}
		}

		utils.Log.WithField("user_id", client.UserID).Info("WebSocket client disconnected")
	}
}

// broadcastMessage broadcasts a message to all clients
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// SendToUser sends a message to all clients of a specific user
func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if userClients, exists := h.userClients[userID]; exists {
		for client := range userClients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// SendToRoom sends a message to all clients in a specific room
func (h *Hub) SendToRoom(roomID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if roomClients, exists := h.rooms[roomID]; exists {
		for client := range roomClients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	client.Rooms[roomID] = true
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if roomClients, exists := h.rooms[roomID]; exists {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			delete(h.rooms, roomID)
		}
	}
	delete(client.Rooms, roomID)
}

// GetUserClientCount returns the number of active connections for a user
func (h *Hub) GetUserClientCount(userID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if userClients, exists := h.userClients[userID]; exists {
		return len(userClients)
	}
	return 0
}
