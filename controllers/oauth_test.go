//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"testing"

	"github.com/cloudmanic/skyclerk.com/library/test"
	"github.com/cloudmanic/skyclerk.com/models"
)

//
// Create a test AcctUser
//
func createTestUserInDB(db *models.DB) {
	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account := test.GetRandomAccount(33)
	account.OwnerId = user.Id
	db.Save(&account)

	app := test.GetRandomApplication()
	db.Save(&app)

	// Set new user account to users
	lu := models.AcctToUsers{AcctId: account.Id, UserId: user.Id}
	db.Save(&lu)
}

//
// TestDoOauthToken01 test to make sure auth works.
//
func TestDoOauthToken01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	createTestUserInDB(db)
}

/* End File */
