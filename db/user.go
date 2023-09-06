package db

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
)

// User contains data of each user account in the application
type User struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique"`
	Email       string `gorm:"varchar(50),unique"`
	Password    string `gorm:"varchar(64)"`
	FirstName   string `gorm:"varchar(50)"`
	LastName    string `gorm:"varchar(50)"`
	PhoneNumber string `gorm:"varchar(30),unique"`
	Gender      string `gorm:"varchar(10)"`
}

// CreateNewUser creates a new user in the database using GormDB
func (gdb *GormDB) CreateNewUser(u *User) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err != nil {
		return err
	} else {
		u.Password = string(pw)
	}

	// Check if there is a duplicate username
	var count int64
	gdb.db.Model(&User{}).Where(&User{Username: u.Username}).Count(&count)
	if count > 0 {
		return errors.New("this username is already taken")
	}

	// Check if there is a duplicate email
	gdb.db.Model(&User{}).Where(&User{Email: u.Email}).Count(&count)
	if count > 0 {
		return errors.New("this email already exists")
	}

	// Check if there is a duplicate phone number
	gdb.db.Model(&User{}).Where(&User{PhoneNumber: u.PhoneNumber}).Count(&count)
	if count > 0 {
		return errors.New("this phone number already exists")
	}

	// Check if the email format is correct
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	// Check if the phone number format is correct (this format only accepts phone numbers
	// that start with 0 and have 11 digits)
	if len(u.PhoneNumber) != 11 || u.PhoneNumber[0] != '0' {
		return errors.New("invalid phone number format")
	}

	return gdb.db.Create(u).Error
}

// GetUserByUsername returns the User that has the exact same username
func (gdb *GormDB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := gdb.db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
