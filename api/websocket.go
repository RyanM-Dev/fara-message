package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mhghw/fara-message/db"
)

type Hub struct {
	clients     map[string]*Client
	chatClients map[string][]*Client
	broadcast   chan Message
	register    chan *Client
	unregister  chan *Client
	mu          sync.Mutex
}

type Client struct {
	user    db.UserTable
	hub     *Hub
	conn    *websocket.Conn
	send    chan db.Message
	receive chan db.Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var users = make(map[string]*websocket.Conn)

func newHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan Message),
		chatClients: make(map[string][]*Client),
	}
}

func (h *Hub) run() {

	for {
		select {
		case client := <-h.register:
			h.clients[client.user.ID] = client
			h.addClientToChatClients(client)
			log.Println("Client registered")
		case client := <-h.unregister:
			if _, ok := h.clients[client.user.ID]; ok {
				delete(h.clients, client.user.ID)
				close(client.send)
				close(client.receive)
				log.Println("Client unregistered")

			}
		case message := <-h.broadcast:
			log.Println("Message received from broadcast")
			h.sendMessage(message)
		}
	}
}

func serveWs(hub *Hub, c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error get ID:%v", err)
		c.Status(400)
		return
	}
	user, err := db.Mysql.ReadUser(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("failed to upgrade connection to websocket:%v", err)
	}
	client := &Client{hub: hub,
		conn:    conn,
		user:    user,
		send:    make(chan db.Message),
		receive: make(chan db.Message),
	}
	users[client.user.ID] = conn
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}

func (h *Hub) sendMessage(msg Message) error {
	exist, err := db.Mysql.CheckChatMemberExists(msg.SenderID, msg.ChatID)
	if err != nil {
		return fmt.Errorf("failed to check chat member existence:%v", err)
	}
	if !exist {
		return fmt.Errorf("chat member does not exist")
	}
	dbMessage, err := db.Mysql.SendMessage(msg.SenderID, msg.ChatID, msg.Content)
	if err != nil {
		log.Printf("failed to write message in database: %v", err)
		return err
	}
	clients, ok := h.chatClients[msg.ChatID]

	if !ok {
		log.Println("no clients exist in the chat map")
	}
	log.Println("client size:", len(clients))
	for _, client := range clients {
		log.Println("starting sending message to clients")
		client.receive <- dbMessage
		log.Println("Message received")
	}

	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading json: %v", err)
			break
		}
		msg.SenderID = c.user.ID
		c.hub.broadcast <- msg
		log.Println("Message sent to broadcast")

	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case msg := <-c.receive:
			log.Println("Message received from client")
			if msg.UserTableID != c.user.ID {
				err := c.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("error writing json: %v", err)
					return
				}

			}

		}
	}
}

func (h *Hub) addClientToChatClients(client *Client) error {
	user := client.user
	userChatMembers, err := db.Mysql.GetUsersChatMembers(user.ID)
	if err != nil {
		log.Println("failed to get users chat members")
		return err
	}
	for _, v := range userChatMembers {
		h.chatClients[v.ChatTableID] = append(h.chatClients[v.ChatTableID], client)
		log.Println("client activated in chat map")
	}
	log.Println("client added to chat map")
	return nil
}
