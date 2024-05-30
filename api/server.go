package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RunWebServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	router := gin.New()
	router.POST("/register", RegisterHandler)
	router.POST("/login", loginHandler)
	router.Use(AuthMiddlewareHandler)
	router.POST("/user/info", ReadUserHandler)
	router.POST("user/change_password", changePassword)
	router.POST("/user/update", UpdateUserHandler)
	router.DELETE("/user/delete", DeleteUserHandler)
	router.POST("/user/edit", editUser)
	router.POST("/user/contact/:id", addContactHandler)
	router.DELETE("/user/contact/:id", DeleteContactHandler)
	router.GET("/user/contact", GetUserContactsHandler)

	router.POST("/send/message", SendMessageHandler)
	router.DELETE("/delete/message", DeleteMessageHandler)
	router.POST("/chat/direct", NewDirectChatHandler)
	router.POST("/chat/group", NewGroupChatHandler)
	router.GET("/chat/:id", GetChatMessagesHandler)
	router.GET("/user/chat/list", GetUsersChatsHandler)
	err := router.Run(addr)
	return err
}
