package api

import (
	"time"

	"github.com/mhghw/fara-message/db"
)

type HTTPError struct {
	Message string `json:"message"`
}

type ContactResponse struct {
	Contacts []Contact `json:"contact,omitempty"`
}
type Contact struct {
	ID          string `json:"id"`
	UserName    string `json:"user_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

func convertContactTableToContact(contactTable db.ContactTable) Contact {
	dateOfBirth := contactTable.Contact.DateOfBirth.Format("2006-01-02")
	gender := ConvertGenderToString(contactTable.Contact.Gender)
	contact := Contact{
		ID:          contactTable.ContactID,
		UserName:    contactTable.Contact.Username,
		FirstName:   contactTable.Contact.FirstName,
		LastName:    contactTable.Contact.LastName,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
	}
	return contact
}
func ConvertGenderToString(genderNumber int8) string {
	var gender string
	switch genderNumber {
	case 0:
		gender = "Male"
	case 1:
		gender = "Female"
	case 2:
		gender = "NonBinary"
	}
	return gender
}
func convertUserTableToRegisterForm(user db.UserTable) RegisterForm {
	gender := ConvertGenderToString(user.Gender)
	formattedTime := user.DateOfBirth.Format(time.DateOnly)
	result := RegisterForm{
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Password:        user.Password,
		ConfirmPassword: user.Password,
		Gender:          gender,
		Email:           user.Email,
		DateOfBirth:     formattedTime,
	}
	return result
}
