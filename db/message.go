package db

import "fmt"

func (d *Database) SendMessage(senderID string, chatID string, content string) (Message, error) {
	var message Message
	message.UserTableID = senderID
	message.ChatTableID = chatID
	message.Content = content
	result := d.db.Create(&message)
	if result.Error != nil {
		return Message{}, fmt.Errorf("error sending message: %w", result.Error)
	}
	err := d.db.Preload("UserTable").Preload("ChatTable").First(&message, message.ID).Error
	if err != nil {
		return Message{}, fmt.Errorf("error retrieving complete message: %w", err)
	}
	return message, nil
}

func (d *Database) DeleteMessage(message Message) error {
	result := d.db.Where("ID=?", message.ID).Delete(&message)
	if result.Error != nil {
		return fmt.Errorf("error deleting message: %w", result.Error)
	}
	return nil
}

func (d *Database) GetUserMessage(messageID int, userID string) (Message, error) {
	var message Message
	if err := d.db.Preload("UserTable").Preload("chat_tables").Where("ID=?", messageID).Where("user_table_id = ?", userID).First(&message).Error; err != nil {
		return message, err
	}
	return message, nil
}
