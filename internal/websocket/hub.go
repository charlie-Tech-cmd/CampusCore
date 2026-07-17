package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Client represents a single active, open WebSocket client channel connection
type Client struct {
	UserID string
	Send   chan []byte // Buffered channel for outbound broadcast payloads
}

// EventPayload defines the structured data layout transmitted to the frontend interface
type EventPayload struct {
	Type    string      `json:"type"` // e.g., "result_approved", "payment_cleared", "ticket_update"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Hub orchestrates the stateful registration, removal, and routing of live WebSocket frames
type Hub struct {
	// Map of active connections tracking clients by their authenticated institutional UserID
	clients   map[string][]*Client
	clientsMu sync.RWMutex

	// Inbound message channels to handle asynchronous event registrations safely
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan EventPayload
}

// NewHub instantiates a thread-safe state management engine for real-time streams
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan EventPayload, 250), // Generous queue buffer space to capture bursts
	}
}

// Run executes our channel selection engine in a blocking background thread loop
func (h *Hub) Run() {
	log.Println("⚡ Real-time WebSocket Hub engine safely spinning up...")
	for {
		select {
		case client := <-h.Register:
			h.clientsMu.Lock()
			h.clients[client.UserID] = append(h.clients[client.UserID], client)
			h.clientsMu.Unlock()
			log.Printf("🔌 Real-time client connected: User %s (Active profiles tracked: %d)",
				client.UserID, len(h.clients[client.UserID]),
			)

		case client := <-h.Unregister:
			h.clientsMu.Lock()
			if connections, exists := h.clients[client.UserID]; exists {
				// Filter out the disconnecting client pointer cleanly from the array stack
				var updatedConnections []*Client
				for _, conn := range connections {
					if conn != client {
						updatedConnections = append(updatedConnections, conn)
					}
				}

				if len(updatedConnections) == 0 {
					delete(h.clients, client.UserID)
				} else {
					h.clients[client.UserID] = updatedConnections
				}
				close(client.Send) // Prevent memory or channel leaks
			}
			h.clientsMu.Unlock()
			log.Printf("🔌 Socket dropped or closed: User %s disconnected safely.", client.UserID)

		case event := <-h.Broadcast:
			// Serialize payload block defensively before transport transmission loops
			payloadBytes, err := json.Marshal(event)
			if err != nil {
				log.Printf("❌ Real-time serialization error: %v", err)
				continue
			}

			h.clientsMu.RLock()
			// Fan-out the real-time update dynamically to all listening browser tabs globally
			for _, connections := range h.clients {
				for _, client := range connections {
					select {
					case client.Send <- payloadBytes:
					default:
						// Defensive step: Drop unresponsive client buffers to prevent blocking the Hub loop
						log.Printf("⚠️ Client buffer saturated for user %s, discarding frame.", client.UserID)
					}
				}
			}
			h.clientsMu.RUnlock()
		}
	}
}

// SendToUser dispatches a targeted notification event packet exclusively to a single designated recipient account
func (h *Hub) SendToUser(userID string, event EventPayload) {
	payloadBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Real-time direct transmission marshalling error: %v", err)
		return
	}

	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()

	connections, exists := h.clients[userID]
	if !exists || len(connections) == 0 {
		return // User is currently offline; gracefully pass through
	}

	for _, client := range connections {
		select {
		case client.Send <- payloadBytes:
		default:
			// Non-blocking fallback if an individual buffer gets full
		}
	}
}
