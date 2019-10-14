//
// Date: 3/3/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	"html/template"

	"golang.org/x/crypto/bcrypt"

	"app.skyclerk.com/backend/library/checkmail"
	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/sendy"
	"app.skyclerk.com/backend/library/slack"
	"app.skyclerk.com/backend/services"
)

// User struct
type User struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `sql:"not null" json:"-"`
	UpdatedAt    time.Time `sql:"not null" json:"-"`
	FirstName    string    `sql:"not null" json:"first_name"`
	LastName     string    `sql:"not null" json:"last_name"`
	Email        string    `sql:"not null" json:"email"`
	Password     string    `sql:"not null" json:"-"`
	Md5Password  string    `sql:"not null" json:"-"`
	Md5Salt      string    `sql:"not null" json:"-"`
	Status       string    `sql:"not null;type:ENUM('Active', 'Disable');default:'Active'" json:"-"`
	LastActivity time.Time `sql:"not null" json:"last_activity"`
	SignupIp     string    `sql:"not null" json:"-"`
	Accounts     []Account `json:"accounts"`
	Session      Session   `json:"-"`
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
	const errMsg = "Sorry, we were unable to find our account."

	var user User
	var session Session

	// See if we already have this user.
	user, err := t.GetUserByEmail(email)

	if err != nil {
		return user, Session{}, errors.New(errMsg)
	}

	// Validate password here by comparing hashes nil means success
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// Check MD5 login
		if (len(user.Md5Salt) > 0) && (len(user.Md5Password) > 0) {
			passMd5 := helpers.GetMd5(password + user.Md5Salt)
			if passMd5 != user.Md5Password {
				return user, Session{}, errors.New(errMsg)
			} else {
				// Create new hash
				hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

				if err != nil {
					return user, Session{}, errors.New(err.Error() + "LoginUserByEmailPass - Unable to create password hash (password hash)")
				}

				// Remove MD5 password and add new hash
				user.Password = string(hash)
				user.Md5Salt = ""
				user.Md5Password = ""
				t.New().Save(&user)
			}
		} else {
			return user, Session{}, errors.New(errMsg)
		}
	}

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
// Create a new user.
//
func (t *DB) CreateUser(first string, last string, email string, password string, appId uint, userAgent string, ipAddress string) (User, error) {
	// Lets do some validation
	if err := t.ValidateCreateUser(first, last, email, false); err != nil {
		return User{}, err
	}

	// Make sure the password is at least 6 chars long
	if err := t.ValidatePassword(password); err != nil {
		return User{}, err
	}

	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		services.InfoMsg(err.Error() + "CreateUser - Unable to create password hash (password hash)")
		return User{}, err
	}

	// Install user into the database
	var _first = template.HTMLEscapeString(first)
	var _last = template.HTMLEscapeString(last)

	user := User{FirstName: _first, LastName: _last, Email: email, Password: string(hash), Status: "Active", LastActivity: time.Now(), SignupIp: ipAddress}
	t.Create(&user)

	// Log user creation.
	services.InfoMsg("CreateUser - Created a new user account - " + first + " " + last + " " + email)

	// Create a session so we get an access_token
	session, err := t.CreateSession(user.Id, appId, userAgent, ipAddress)

	if err != nil {
		services.InfoMsg(err.Error() + "CreateUser - Unable to create session in CreateSession()")
		return User{}, err
	}

	// Add the session to the user object.
	user.Session = session

	// Do post register stuff
	t.doPostUserRegisterStuff(user, ipAddress)

	// Return the user.
	return user, nil
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

//
// ValidateCreateUser a create user action. We do not always get a first name and last name from google.
// so we make the validation optional with them.
//
func (t *DB) ValidateCreateUser(first string, last string, email string, googleAuth bool) error {
	// Are first and last name fields empty
	if (!googleAuth) && (len(first) == 0) && (len(last) == 0) {
		return errors.New("First name and last name fields are required.")
	}

	// Are first name empty
	if (!googleAuth) && len(first) == 0 {
		return errors.New("First name field is required.")
	}

	// Are last name empty
	if (!googleAuth) && len(last) == 0 {
		return errors.New("Last name field is required.")
	}

	// Lets validate the email address
	if err := t.ValidateEmailAddress(email); err != nil {
		return err
	}

	// See if we already have this user.
	_, err := t.GetUserByEmail(email)

	if err == nil {
		return errors.New("Looks like you already have an account.")
	}

	// Return happy.
	return nil
}

// ------------------ Helper Functions --------------------- //

//
// Do post user register stuff.
//
func (t *DB) doPostUserRegisterStuff(user User, ipAddress string) {
	// Subscribe new user to mailing lists.
	go sendy.Subscribe("trial", user.Email, user.FirstName, user.LastName, "", "", ipAddress, "No")
	go sendy.Subscribe("subscribers", user.Email, user.FirstName, user.LastName, "No", "", ipAddress, "No")

	// Tell slack about this.
	go slack.Notify("#events", "New Skyclerk User Account : "+user.Email)
}

/* End File */
