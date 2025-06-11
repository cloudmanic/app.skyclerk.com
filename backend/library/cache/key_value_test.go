//
// Date: 2/27/2017
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cache

import (
	"testing"

	"app.skyclerk.com/backend/models"
	"github.com/nbio/st"
)


//
// Test - Set 01
//
func TestSet01(t *testing.T) {
	// Start the db connection.
	testDB, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(testDB, dbName)

	// Set the cache database
	db = testDB
	
	// Store something in cache
	Set("sc-testing-1", "Skyclerk is DaBomb.com")

	// Get stored value.
	result := ""
	found, err := Get("sc-testing-1", &result)

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, found, true)
	st.Expect(t, result, "Skyclerk is DaBomb.com")
}

//
// Test - Set 02
//
func TestSet02(t *testing.T) {
	// Start the db connection.
	testDB, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(testDB, dbName)

	// Set the cache database
	db = testDB
	
	// Get a value we know we do not have
	result := ""
	found, _ := Get("sc-testing-not-found", &result)

	// Verify the data was return as expected
	st.Expect(t, found, false)
	st.Expect(t, result, "")
}

//
// Test - Set 03
//
func TestSet03(t *testing.T) {
	// Start the db connection.
	testDB, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(testDB, dbName)

	// Set the cache database
	db = testDB

	// Create an Billing model.
	b := models.Billing{
		Status: "Active",
	}

	// Store the struct in the cache
	Set("sc-testing-2", b)

	// Get a value we know we do not have
	result := models.Billing{}
	found, _ := Get("sc-testing-2", &result)

	// Verify the data was return as expected
	st.Expect(t, found, true)
	st.Expect(t, result.Status, "Active")
}

/* End File */
