//
// Date: 3/3/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	"github.com/cloudmanic/skyclerk.com/library/checkmail"
	"github.com/cloudmanic/skyclerk.com/library/helpers"
	"github.com/cloudmanic/skyclerk.com/services"
)

// User struct
type User struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	FirstName    string    `sql:"not null" json:"first_name"`
	LastName     string    `sql:"not null" json:"last_name"`
	Email        string    `sql:"not null" json:"email"`
	Md5Password  string    `sql:"not null" json:"-"`
	Md5Salt      string    `sql:"not null" json:"-"`
	Status       string    `sql:"not null;type:ENUM('Active', 'Disable');default:'Active'" json:"-"`
	LastActivity time.Time `json:"last_activity"`
	Accounts     []Account `json:"accounts"`
}

//
// GetUserById - Get a user by Id.
//
func (t *DB) GetUserById(id uint) (User, error) {
	var u User

	if t.Where("id = ?", id).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Add in accounts TOOD(user): Clean this up to be more GORMY
	aTU := []AcctToUsers{}
	t.Where("user_id = ?", id).Order("acct_id DESC").Find(&aTU)

	for _, row := range aTU {
		a := Account{}
		t.Find(&a, row.AcctId)
		u.Accounts = append(u.Accounts, a)
	}

	// Return the user.
	return u, nil
}

//
// GetUserByEmail - Get a user by email.
//
func (t *DB) GetUserByEmail(email string) (User, error) {
	var u User

	if t.Where("email = ?", email).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil
}

//
// LoginUserByEmailPass - Login a user in by email and password. The userAgent is a way to marking what device this
// login request came from. Same with ipAddress.
//
func (t *DB) LoginUserByEmailPass(email string, password string, appId uint, userAgent string, ipAddress string) (User, Session, error) {
	var user User
	var session Session

	// See if we already have this user.
	user, err := t.GetUserByEmail(email)

	if err != nil {
		return user, Session{}, errors.New("Sorry, we were unable to find our account.")
	}

	// Do MD5 login
	passMd5 := helpers.GetMd5(password + user.Md5Salt)
	if passMd5 != user.Md5Password {
		return user, Session{}, errors.New("Sorry, we were unable to find our account.")
	}

	// TODO(spicer): Support non-md5 passwords
	// // Validate password here by comparing hashes nil means success
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	//
	// if err != nil {
	// 	return user, err
	// }

	// Create a session so we get an access_token (if we passed in an appId)
	if appId > 0 {
		s, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

		if err != nil {
			services.Error(err)
			return User{}, Session{}, err
		}

		session = s
	}

	return user, session, nil
}

//
// ValidateUserLogin - Validate a login user action.
//
func (t *DB) ValidateUserLogin(email string, password string) error {

	// Make sure the password is at least 6 chars long
	if len(password) < 6 {
		return errors.New("The password filed must be at least 6 characters long.")
	}

	// Lets validate the email address
	if err := t.ValidateEmailAddress(email); err != nil {
		return err
	}

	// See if we already have this user.
	_, err := t.GetUserByEmail(email)

	if err != nil {
		return errors.New("Sorry, we were unable to find our account.")
	}

	// Return happy.
	return nil
}

//
// ValidatePassword - Validate password.
//
func (t *DB) ValidatePassword(password string) error {

	// Make sure the password is at least 6 chars long
	if len(password) < 6 {
		return errors.New("The password filed must be at least 6 characters long.")
	}

	// Return happy.
	return nil

}

//
// ValidateEmailAddress - Validate an email address
//
func (t *DB) ValidateEmailAddress(email string) error {

	// Check length
	if len(email) == 0 {
		return errors.New("Email address field is required.")
	}

	// Check format
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("Email address is not a valid format.")
	}

	// Return happy.
	return nil

}

/* End File */
