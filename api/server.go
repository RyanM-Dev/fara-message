package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RunWebServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	router := gin.New()
	router.POST("/user/register", Register)
	router.POST("/user/change_password", changePassword)
	router.POST("/login", authenticateUser)
	router.Use(AuthMiddleware)
	// router.POST("/user/create",CreateUserHandler)
	router.POST("/user/read", ReadUserHandler)
	router.POST("/user/update", UpdateUserHandler)
	router.POST("/user/delete", DeleteUserHandler)
	router.POST("/user/edit", editUser)
	router.POST("/send/message", SendMessageHandler)
	router.DELETE("/delete/message", DeleteMessageHandler)
	router.POST("/new_direct_chat", NewDirectChat)
	router.POST("/new_group_chat", NewGroupChat)
	router.GET("/chat/:id/messages", GetChatMessages)
	err := router.Run(addr)
	return err
}