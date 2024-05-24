package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type Message struct {
	ID      string `json:"id"`
	ChatID  string `json:"chatID"`
	Content string `json:"content"`
}

func SendMessageHandler(c *gin.Context) {
	var message Message
	if err := c.BindJSON(&message); err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}

	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error get ID:%v", err)
		c.Status(400)
		return
	}

	err = db.Mysql.SendMessage(userID, message.ChatID, message.Content)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "message sent successfully",
	})
}

func DeleteMessageHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error get ID:%v", err)
		c.Status(400)
		return
	}

	var message Message
	if err := c.BindJSON(&message); err != nil {
		log.Printf("error binding JSON:%v", err)
	}
	messageID, err := strconv.Atoi(message.ID)
	if err != nil {
		log.Printf("error converting string ID to int: %v", err)
		c.JSON(400, "error converting string ID to int")
		return
	}
	dbMessage, err := db.Mysql.GetUserMessage(messageID, userID)
	if err != nil {
		log.Printf("error finding message with provided message and user ID: %v", err)
		c.JSON(400, "error finding message with provided message and user ID")
		return

	}
	err = db.Mysql.DeleteMessage(dbMessage)
	if err != nil {
		log.Printf("error:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "message deleted successfully",
	})
}
