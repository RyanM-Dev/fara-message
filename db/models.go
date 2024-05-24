package db

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/rs/xid"
)

type Message struct {
	ID          int
	UserTableID string `gorm:"type:varchar(255)"`
	UserTable   UserTable
	ChatTableID string `gorm:"type:varchar(255)"`
	ChatTable   ChatTable
	Content     string
}

type Gender struct {
	gender int
}

var (
	Male      = Gender{gender: 0}
	Female    = Gender{gender: 1}
	NonBinary = Gender{gender: 2}
)

type User struct {
	ID          xid.ID
	Username    string
	FirstName   string
	LastName    string
	Password    string
	Gender      Gender
	DateOfBirth time.Time
	CreatedTime time.Time
	DeletedTime time.Time
}
type UserTable struct {
	ID          string `gorm:"type:varchar(255)"`
	Username    string
	FirstName   string
	LastName    string
	Password    string
	Gender      int8
	DateOfBirth time.Time
	CreatedTime time.Time
	DeletedTime sql.NullTime
}

type Chat struct {
	ID          string
	Name        string
	CreatedTime time.Time
	DeletedTime time.Time
	Type        ChatType
}
type ChatTable struct {
	ID          string `gorm:"type:varchar(255)"`
	Name        string
	CreatedTime time.Time
	DeletedTime sql.NullTime
	Type        int8
}
type ChatMember struct {
	UserTableID string `gorm:"type:varchar(255)"`
	UserTable   UserTable
	ChatTableID string `gorm:"type:varchar(255)"`
	ChatTable   ChatTable
	JoinedTime  time.Time
	LeftTime    sql.NullTime
}

type ChatType struct {
	chatType int
}
type ChatIDAndChatName struct {
	ChatID   string
	ChatName string
}

func (c *ChatType) Int() int {
	return c.chatType
}

var (
	Direct  = ChatType{chatType: 0}
	Group   = ChatType{chatType: 1}
	Unknown = ChatType{chatType: -1}
)

type ContactTable struct {
	UserTableID string `gorm:"type:varchar(255)"`
	UserTable   UserTable
	ContactID   string `gorm:"type:varchar(255)"`
	Contact     UserTable
}

func ConvertUserToUserTable(user User) UserTable {
	var gender int8
	switch user.Gender.gender {
	case 0:
		gender = 0
	case 1:
		gender = 1
	case 2:
		gender = 2
	}
	userTable := UserTable{
		ID:          user.ID.String(),
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    user.Password,
		Gender:      gender,
		DateOfBirth: user.DateOfBirth,
		CreatedTime: user.CreatedTime,
	}
	return userTable
}

func ConvertChatToChatTable(chat Chat) ChatTable {
	var chatType int8
	switch chat.Type {
	case Direct:
		chatType = 0
	case Group:
		chatType = 1
	default:
		chatType = -1
	}

	result := ChatTable{
		ID:          chat.ID,
		Name:        chat.Name,
		CreatedTime: chat.CreatedTime,
		Type:        chatType,
	}
	return result
}

func hashDB(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}
