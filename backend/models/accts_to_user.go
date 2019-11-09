//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

// AcctToUsers struct NOTE: We do not run this in our migrations. It is auto creeated.
type AcctToUsers struct {
	AccountId uint `sql:"not null"  json:"_"`
	UserId    uint `sql:"not null"  json:"_"`
}

//
// GetUsersByAccount returns a list of users within your account.
//
func (t *DB) GetUsersByAccount(accountId uint) []User {
	// SQL String
	sql := "SELECT users.* FROM acct_to_users "
	sql = sql + "JOIN users ON users.id = acct_to_users.user_id "
	sql = sql + "WHERE acct_to_users.account_id = ? ORDER BY users.id"

	// Struct we return
	rt := []User{}

	// Run query
	t.Raw(sql, accountId).Scan(&rt)

	// Return happy.
	return rt
}

/* End File */
