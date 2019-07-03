//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// AcctUsersLu struct
type AcctToUsers struct {
	Id        uint      `gorm:"primary_key" json:"_"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AcctId    uint      `sql:"not null"  json:"_"`
	UserId    uint      `sql:"not null"  json:"_"`
}

//
// GetUsersByAccount returns a list of users within your account.
//
func (t *DB) GetUsersByAccount(accountId uint) []User {
	// SQL String
	sql := "SELECT users.* FROM acct_to_users "
	sql = sql + "JOIN users ON users.id = acct_to_users.user_id "
	sql = sql + "WHERE acct_to_users.acct_id = ? ORDER BY users.id"

	// Struct we return
	rt := []User{}

	// Run query
	t.Raw(sql, accountId).Scan(&rt)

	// Return happy.
	return rt
}

/* End File */
