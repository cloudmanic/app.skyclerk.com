//
// Date: 3/3/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

type Application struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `sql:"not null" json:"name"`
	ClientId  string    `sql:"not null" json:"client_id"`
	Secret    string    `sql:"not null" json:"secret"`
	GrantType string    `sql:"not null;type:ENUM('password', 'authorization_code');default:'password'" json:"grant_type"`
}

//
// Validate a client id and grant type.
//
func (db *DB) ValidateClientIdGrantType(clientId string, grantType string) (Application, error) {

	var u Application

	if db.Where("client_id = ? AND grant_type= ?", clientId, grantType).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil

}

/* End File */
