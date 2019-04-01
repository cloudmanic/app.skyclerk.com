//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

// AcctUsersLu struct
type AcctToUsers struct {
	Id        uint      `gorm:"primary_key" json:"_"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AcctId    uint      `sql:"not null"  json:"_"`
	UserId    uint      `sql:"not null"  json:"_"`
}

/* End File */
