package models

import (
	"errors"
	"time"

	"github.com/cloudmanic/skyclerk.com/backend/library/helpers"
	"github.com/cloudmanic/skyclerk.com/backend/services"
)

// Session struct
type Session struct {
	Id            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserId        uint   `sql:"not null;index:UserId"`
	ApplicationId uint   `sql:"not null;index:ApplicationId"  json:"-"`
	UserAgent     string `sql:"not null"`
	AccessToken   string `sql:"not null"`
	LastIpAddress string `sql:"not null"`
	LastActivity  time.Time
}

//
// Get by Access token.
//
func (t *DB) GetByAccessToken(accessToken string) (Session, error) {
	// Session
	var sess Session

	if t.First(&sess, "access_token = ?", accessToken).RecordNotFound() {
		return Session{}, errors.New("Access Token Not Found - Unable to Authenticate (#001)")
	}

	// TODO(spicer): figure out why this does now work.
	// // Double check because of case sensitivity
	// if sess.AccessToken == accessToken {
	// 	return Session{}, errors.New("Access Token Not Found - Unable to Authenticate (#002)")
	// }

	// Return happy
	return sess, nil
}

//
// Create new session. A user can have more than one session. Typically it is one session per browser or device.
// We return the session object. The big thing here is we create the access token for this session.
//
func (db *DB) CreateSession(UserId uint, appId uint, UserAgent string, LastIpAddress string) (Session, error) {
	// Create an access token.
	access_token, err := helpers.GenerateRandomString(50)

	if err != nil {
		services.Error(err)
		return Session{}, err
	}

	// Save the session into the database.
	session := Session{UserId: UserId, ApplicationId: appId, UserAgent: UserAgent, AccessToken: access_token, LastIpAddress: LastIpAddress}
	db.Create(&session)

	// Return the session.
	return session, nil
}

/* End File */
