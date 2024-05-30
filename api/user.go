package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhghw/fara-message/db"
)

type UsernameType struct {
	Username string `json:"username"`
}

func ReadUserHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("failed to validate token: %v", err)
		return
	}
	log.Println(userID)
	userTable, err := db.Mysql.ReadUser(userID)
	if err != nil {
		log.Printf("failed to read user: %v", err)
		return
	}
	RegisterForm := convertUserTableToRegisterForm(userTable)
	c.JSON(http.StatusOK, RegisterForm)

}

func UpdateUserHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	userID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("error validating user: %v", err)
		c.JSON(400, "error validating user")

		return
	}
	oldUser, err := db.Mysql.ReadUser(userID)
	if err != nil {
		log.Printf("error reading user:%v", err)
		c.JSON(400, "error reading user")
		return
	}
	oldUserRegisterForm := convertUserTableToRegisterForm(oldUser)
	var newInfoRequest RegisterForm
	err = c.BindJSON(&newInfoRequest)
	if err != nil {
		log.Printf("error binding JSON:%v", err)
		c.Status(400)
		return
	}
	newInfoRequest.DateOfBirth = oldUserRegisterForm.DateOfBirth

	user, err := convertRegisterFormToUser(newInfoRequest)
	if err != nil {
		log.Printf("error converting registerForm to user:%v", err)
		c.Status(400)
		return
	}
	newUserTable := db.ConvertUserToUserTable(user)
	fmt.Println(newUserTable)
	err = db.Mysql.UpdateUser(userID, newUserTable)
	if err != nil {
		log.Printf("error updating user:%v", err)
		c.Status(400)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})

}

func DeleteUserHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.Status(400)
		return
	}

	err = db.Mysql.DeleteUser(userID)
	if err != nil {
		log.Printf("failed to delete user:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("failed to delete user:%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})

}

func addContactHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.Status(400)
		return
	}
	contactID := c.Param("id")
	if err := db.Mysql.AddContact(userID, contactID); err != nil {
		log.Printf("failed to add contact:%v", err)
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(200, "contact added successfully")
}
func GetUserContactsHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.Status(400)
		return
	}
	contactsDB, err := db.Mysql.GetUserContacts(userID)
	if err != nil {
		log.Printf("failed to get contacts:%v", err)
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	var contacts []Contact
	for _, v := range contactsDB {
		contacts = append(contacts, convertContactTableToContact(v))
	}

	contactResponse := ContactResponse{
		Contacts: contacts,
	}

	c.JSON(200, contactResponse)
}

func DeleteContactHandler(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	userID, err := ValidateToken(authorizationHeader)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.Status(400)
		return
	}
	contactID := c.Param("id")
	if err := db.Mysql.DeleteContact(userID, contactID); err != nil {
		log.Printf("failed to delete contact:%v", err)
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(200, "contact deleted successfully")
}
