//
// Date: 3/3/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
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
	Session      Session   `json:"-"`
	LastActivity time.Time `json:"last_activity"`
}

/* End File */
