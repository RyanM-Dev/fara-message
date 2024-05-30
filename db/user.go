package db

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func (d *Database) CreateUser(user User) error {
	userTable := ConvertUserToUserTable(user)
	result := d.db.Create(&userTable)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (d *Database) ReadUserByUsername(username string) (UserTable, error) {
	var user UserTable
	result := d.db.Where("username	 = ?", username).First(&user)
	if result.Error != nil {
		return user, fmt.Errorf("failed to read user: %w", result.Error)
	}
	return user, nil
}

func (d *Database) ReadUser(ID string) (UserTable, error) {
	var user UserTable
	result := d.db.Select("id", "username", "first_name", "last_name", "gender", "email", "date_of_birth", "created_time").
		Where("id = ?", ID).
		First(&user)
	if result.Error != nil {
		return user, fmt.Errorf("failed to read user: %w", result.Error)
	}
	return user, nil
}

func (d *Database) UpdateUser(ID string, newInfo UserTable) error {

	result := d.db.Model(&UserTable{}).Where("ID=?", ID).Updates(UserTable{Username: newInfo.Username, FirstName: newInfo.FirstName, LastName: newInfo.LastName, Password: newInfo.Password, Gender: newInfo.Gender, DateOfBirth: newInfo.DateOfBirth})
	if result.Error != nil {
		return fmt.Errorf("failed to Update user: %w", result.Error)
	}
	return nil
}

func (d *Database) DeleteUser(ID string) error {
	var user UserTable
	result := d.db.First(&user, "ID=?", ID)
	if result.Error != nil {
		return fmt.Errorf("failed to find user to delete: %w", result.Error)
	}
	result = d.db.Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	return nil
}

func (d *Database) AddContact(userID, contactID string) error {
	contact := ContactTable{
		UserTableID: userID,
		ContactID:   contactID,
	}
	repeatedContact, err := d.isContactExist(userID, contactID)
	if err != nil {
		return fmt.Errorf("failed to check if contact exists: %w", err)
	}
	if repeatedContact {
		return errors.New("contact already exists")
	}
	if err := d.db.Create(&contact).Error; err != nil {
		log.Println("create contact")
		return fmt.Errorf("failed to add contact:%v", err)
	}
	return nil
}

func (d *Database) isContactExist(userID, contactID string) (bool, error) {
	if userID == contactID {
		return false, errors.New("user id  and contact id are the same")
	}
	contact := ContactTable{
		UserTableID: userID,
		ContactID:   contactID,
	}
	result := d.db.First(&contact, "user_table_id = ? AND contact_id = ?", userID, contactID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("record not found")
			return false, nil
		} else {
			return false, fmt.Errorf("failed to find contact:%v", result.Error)

		}

	}
	return true, nil
}

func (d *Database) GetContact(userID, contactID string) (ContactTable, error) {
	var contact ContactTable
	contactExistence, err := d.isContactExist(userID, contactID)
	if err != nil {
		return ContactTable{}, fmt.Errorf("failed to check if contact exists: %w", err)
	}
	if contactExistence {
		if err := d.db.Model(&ContactTable{}).Where("user_table_id = ? AND contact_id = ?", userID, contactID).First(&contact).Error; err != nil {
			return ContactTable{}, fmt.Errorf("failed to get contact: %w", err)
		}
	}
	return ContactTable{}, errors.New("contact not found")
}

func (d *Database) GetUserContacts(userID string) ([]ContactTable, error) {
	var contacts []ContactTable
	if err := d.db.Model(&ContactTable{}).Preload("Contact").Where("user_table_id = ?", userID).Find(&contacts).Error; err != nil {
		return contacts, fmt.Errorf("failed to get contacts:%w", err)
	}
	for i := range contacts {
		user, err := d.ReadUser(contacts[i].ContactID)
		if err != nil {
			return []ContactTable{}, fmt.Errorf("failed to fill contact:%w", err)
		}
		contacts[i].Contact = user
	}
	return contacts, nil
}

func (d *Database) DeleteContact(userID, contactID string) error {
	var contact ContactTable
	contactExistence, err := d.isContactExist(userID, contactID)
	if err != nil {
		return fmt.Errorf("failed to check if contact exists: %w", err)
	}
	if contactExistence {
		if err := d.db.Model(&ContactTable{}).Where("user_table_id = ? AND contact_id = ?", userID, contactID).Delete(&contact).Error; err != nil {
			return fmt.Errorf("failed to delete contact: %w", err)
		}
		return nil
	}
	return errors.New("contact for delete not found")
}
