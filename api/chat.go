package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
	"github.com/rs/xid"
)

type GroupChatRequest struct {
	ChatName string         `json:"chatName"`
	Users    []db.UserTable `json:"users"`
}
type DirectChatRequest struct {
	UserID string `json:"username"`
}

func NewDirectChatHandler(c *gin.Context) {
	var requestBody DirectChatRequest

	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Print("failed to bind json, ", err)
		return
	}
	tokenString := c.GetHeader("Authorization")

	hostUserID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("failed to find user by token: %v", err)
	}
	destinationUserTable, err := db.Mysql.ReadUserByUsername(requestBody.UserID)
	if err != nil {
		log.Printf("failed to read user: %v", err)
	}

	hostUserTable, err := db.Mysql.ReadUser(hostUserID)

	if err != nil {
		log.Printf("failed to read user: %v", err)
	}
	var userTables []db.UserTable
	userTables = append(userTables, hostUserTable, destinationUserTable)
	chatID, err := db.Mysql.NewChat("", db.Direct, userTables)
	if err != nil {
		log.Print("failed to create chat, ", err)
		return
	}

	UpdateChatClientList(userTables, chatID)

	log.Print("direct chat created")
	c.JSON(200, chatID)
}

func NewGroupChatHandler(c *gin.Context) {
	var requestBody GroupChatRequest
	err := c.BindJSON(&requestBody)
	if err != nil {
		log.Print("failed to bind json, ", err)
		return
	}
	tokenString := c.GetHeader("Authorization")

	userID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("failed to find user by token: %v", err)
	}
	userTables := []db.UserTable{}
	for _, v := range requestBody.Users {
		user, err := db.Mysql.ReadUserByUsername(v.Username)
		if err != nil {
			log.Printf("failed to read user: %v", err)
			return
		}
		userTables = append(userTables, user)
	}

	validUser := false
	for _, user := range userTables {
		if user.ID == userID {
			validUser = true
		}
	}
	if !validUser {
		log.Printf("you're not allowed")
		c.JSON(400, "Invalid token")
		return
	}
	if len(userTables) == 0 {
		log.Print("failed to create chat: no users provided")
		return
	}
	log.Println(requestBody.ChatName)
	chatID, err := db.Mysql.NewChat(requestBody.ChatName, db.Group, userTables)
	if err != nil {
		log.Printf("failed to create chat: %v", err)
		return
	}
	UpdateChatClientList(userTables, chatID)

	c.JSON(200, chatID)
}

func AddToGroupChatHandler(c *gin.Context) {
	type requestBody struct {
		ChatID string `json:"chat_id"`
		UserID string `json:"user_id"`
	}
	var reqBody requestBody
	tokenString := c.GetHeader("Authorization")

	userID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("failed to find user by token: %v", err)
	}
	err = c.BindJSON(&reqBody)
	if err != nil {
		log.Print("failed to bind json, ", err)
		return
	}
	err = db.Mysql.AddToGroupChat(reqBody.ChatID, reqBody.UserID, userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	userTable, err := db.Mysql.ReadUser(reqBody.UserID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("no user found: %v", err),
		})
	}
	var userTables = []db.UserTable{userTable}
	UpdateChatClientList(userTables, reqBody.ChatID)

	c.JSON(200, gin.H{
		"status": "user added successfully",
	})
}

func GetChatMessagesHandler(c *gin.Context) {
	chatID := c.Param("id")
	messages, err := db.Mysql.GetChatMessages(chatID)
	if err != nil {
		log.Printf("failed to get messages: %v", err)
		c.JSON(400, "failed to get chat messages")
		return
	}
	c.JSON(200, messages)

}

func GetUsersChatsHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	userID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("failed to find user by token: %v", err)
	}
	chatMembers, err := db.Mysql.GetUsersChatMembers(userID)
	if err != nil {
		log.Printf("failed to get users chats: %v", err)
		return
	}

	result, err := db.Mysql.GetUsersChatIDAndChatName(chatMembers)
	if err != nil {
		log.Printf("failed to get users chats: %v", err)
		return
	}
	c.JSON(200, result)
}

func DirectChatIDGenerator(users []db.User) (string, error) {
	var userIDs []xid.ID
	err := errors.New("too many users provided")
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	if len(userIDs) > 2 {
		log.Println("too many users for direct chat")
		return "", err
	}
	var concatenatedID string
	if userIDs[0].Compare(userIDs[1]) < 0 {
		concatenatedID = userIDs[0].String() + userIDs[1].String()
	} else {
		concatenatedID = userIDs[1].String() + userIDs[0].String()
	}
	hashID := hash(concatenatedID)
	return hashID, nil

}

func UpdateChatClientList(users []db.UserTable, chatID string) {
	hub := GetHub()
	hub.mu.Lock()
	defer hub.mu.Unlock()
	var clients []*Client
	for _, user := range users {
		client, exists := hub.clients[user.ID]
		if exists {
			clients = append(clients, client)
		} else {
			log.Printf("User %s is not connected", user.ID)
		}
	}
	if _, ok := hub.chatClients[chatID]; !ok {
		hub.chatClients[chatID] = []*Client{}
	}
	for _, client := range clients {
		hub.chatClients[chatID] = append(hub.chatClients[chatID], client)
	}
}
