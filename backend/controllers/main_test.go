//
// Date: 2025-01-06
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2025 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"os"
	"testing"

	"app.skyclerk.com/backend/library/cache"
	"app.skyclerk.com/backend/models"
)

// TestMain sets up the test environment for all controller tests
func TestMain(m *testing.M) {
	// Create a test database for cache initialization
	db, dbName, err := models.NewTestDB("")
	if err != nil {
		panic(err)
	}
	
	// Initialize cache with test database
	cache.SetDB(db)
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	models.TestingTearDown(db, dbName)
	
	// Exit
	os.Exit(code)
}