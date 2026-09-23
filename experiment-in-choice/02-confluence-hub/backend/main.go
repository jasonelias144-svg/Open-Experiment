package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"://github.com"
)
// Message holds structural network communications datatype Message struct {
	RoomID     string `json:"roomId"`
	SenderType string `json:"senderType"` // "human" | "ai" | "system"
	SenderID   string `json:"senderId"`
	ContextType string `json:"contextType"` // "scientific" | "artistic" | "tech" | "phi" | "social"
	Content    string `json:"content"`
	BranchID   string `json:"branchId,omitempty"`
}
// Client definition representing a single pipeline nodetype Client struct {
	ID         string
	SenderType string
	Conn       *websocket.Conn
	Send       chan []byte
}
// Room holds client mappings for active roomstype Room struct {
	ID      string
	Clients map[*Client]bool
	Mu      sync.RWMutex
}
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *ClientRegistration
	Unregister chan *ClientRegistration
	Broadcast  chan Message
	Mu         sync.RWMutex
}
type ClientRegistration struct {
	Client *Client
	RoomID string
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *ClientRegistration),
		Unregister: make(chan *ClientRegistration),
		Broadcast:  make(chan Message),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case reg := <-h.Register:
			h.Mu.Lock()
			room, exists := h.Rooms[reg.RoomID]
			if !exists {
				room = &Room{ID: reg.RoomID, Clients: make(map[*Client]bool)}
				h.Rooms[reg.RoomID] = room
			}
			h.Mu.Unlock()

			room.Mu.Lock()
			room.Clients[reg.Client] = true
			room.Mu.Unlock()
			fmt.Printf("[HUB] Node [%s:%s] mounted to Room [%s]\n", reg.Client.SenderType, reg.Client.ID, reg.RoomID)

		case unreg := <-h.Unregister:
			h.Mu.RLock()
			room, exists := h.Rooms[unreg.RoomID]
			h.Mu.RUnlock()

			if exists {
				room.Mu.Lock()
				if _, ok := room.Clients[unreg.Client]; ok {
					delete(room.Clients, unreg.Client)
					close(unreg.Client.Send)
					fmt.Printf("[HUB] Node [%s] demounted from Room [%s]\n", unreg.Client.ID, unreg.RoomID)
				}
				room.Mu.Unlock()
			}

		case msg := <-h.Broadcast:
			h.Mu.RLock()
			room, exists := h.Rooms[msg.RoomID]
			h.Mu.RUnlock()

			if exists {
				payload, _ := json.Marshal(msg)
				room.Mu.RLock()
				for client := range room.Clients {
					select {
					case client.Send <- payload:
					default:
						// Handle channel blockages gracefully
						go func(c *Client) { h.Unregister <- &ClientRegistration{Client: c, RoomID: msg.RoomID} }(client)
					}
				}
				room.Mu.RUnlock()
			}
		}
	}
}
func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		
		query := r.URL.Query()
		roomID := query.Get("room")
		senderType := query.Get("type") // "human", "ai"
		clientID := query.Get("id")

		if roomID == "" || senderType == "" || clientID == "" {
			conn.Close()
			return
		}

		client := &Client{ID: clientID, SenderType: senderType, Conn: conn, Send: make(chan []byte, 256)}
		hub.Register <- &ClientRegistration{Client: client, RoomID: roomID}

		// Handle client socket listening layers
		go func() {
			defer func() {
				hub.Unregister <- &ClientRegistration{Client: client, RoomID: roomID}
				conn.Close()
			}()
			for {
				_, msgBytes, err := conn.ReadMessage()
				if err != nil {
					break
				}
				var baseMsg Message
				if err := json.Unmarshal(msgBytes, &baseMsg); err == nil {
					baseMsg.RoomID = roomID
					baseMsg.SenderID = clientID
					baseMsg.SenderType = senderType
					hub.Broadcast <- baseMsg
				}
			}
		}()

		// Handle client write buffers
		go func() {
			for msg := range client.Send {
				_ = conn.WriteMessage(websocket.TextMessage, msg)
			}
		}()
	})

	fmt.Println("[INIT] Confluence Hub Hub active on :8080/ws")
	_ = http.ListenAndServe(":8080", nil)
}
