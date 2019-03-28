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
type AcctUsersLu struct {
	AcctUsersLuId        uint      `gorm:"primary_key;column:AcctUsersLuId" json:"_"`
	AcctUsersLuAcctId    uint      `gorm:"column:AcctUsersLuAcctId" json:"_"`
	AcctUsersLuUserId    uint      `gorm:"column:AcctUsersLuUserId" json:"_"`
	AcctUsersLuUpdatedAt time.Time `gorm:"column:AcctUsersLuUpdatedAt" sql:"not null" json:"_"`
	AcctUsersLuCreatedAt time.Time `gorm:"column:AcctUsersLuCreatedAt" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (AcctUsersLu) TableName() string {
	return "AcctUsersLu"
}

/* End File */
