//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/services"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)


//
// NewTestDB - Start the Test DB connection.
//
func NewTestDB(dbName string) (*DB, string, error) {
	var db *gorm.DB
	var err error

	// If dbName is empty we create our own.
	if len(dbName) == 0 {
		dbName = "sk_test_" + helpers.RandStr(10) + ".db"
	} else if !strings.HasSuffix(dbName, ".db") {
		dbName = dbName + ".db"
	}

	// Create test database path in cache/sqlite directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	cacheDir := filepath.Join(basepath, "..", "cache", "sqlite")
	
	// Create directory if it doesn't exist
	os.MkdirAll(cacheDir, 0755)
	
	dbPath := filepath.Join(cacheDir, dbName)

	// Connect to SQLite
	db, err = gorm.Open("sqlite3", dbPath)

	if err != nil {
		services.Error(err)
		log.Fatal(err)
	}

	// Enable WAL mode for better performance and concurrency
	if err := db.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
		services.Error(err)
		log.Fatal(err)
	}

	// Set synchronous to NORMAL for better performance while maintaining durability
	if err := db.Exec("PRAGMA synchronous=NORMAL").Error; err != nil {
		services.Error(err)
		log.Fatal(err)
	}

	// Set busy timeout to 5 seconds to handle concurrent access better
	if err := db.Exec("PRAGMA busy_timeout=5000").Error; err != nil {
		services.Error(err)
		log.Fatal(err)
	}

	// Run doMigrations
	doMigrations(db)

	// Clear all tables.
	truncateAllTables(db)

	// Return db connection.
	return &DB{db}, dbName, nil
}


//
// TestingTearDown - used to delete database we created
//
func TestingTearDown(db *DB, dbName string) {
	// Close the database connection first
	db.Close()
	
	// For SQLite, remove the database file
	if strings.HasSuffix(dbName, ".db") {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		cacheDir := filepath.Join(basepath, "..", "cache", "sqlite")
		dbPath := filepath.Join(cacheDir, dbName)
		
		// Remove the database file
		os.Remove(dbPath)
	}
}

//
// TruncateAllTables - Clear all tables. DONT FORGET TO UPDATE ClearAccount() in accounts.go
//
func truncateAllTables(db *gorm.DB) {
	// SQLite doesn't support TRUNCATE, so we use DELETE
	// Delete in order to respect foreign key constraints
	db.Exec("DELETE FROM LabelsToLedger;")
	db.Exec("DELETE FROM FilesToLedger;")
	db.Exec("DELETE FROM SnapClerk;")
	db.Exec("DELETE FROM Ledger;")
	db.Exec("DELETE FROM Labels;")
	db.Exec("DELETE FROM Files;")
	db.Exec("DELETE FROM Contacts;")
	db.Exec("DELETE FROM Categories;")
	db.Exec("DELETE FROM activities;")
	db.Exec("DELETE FROM acct_to_users;")
	db.Exec("DELETE FROM invites;")
	db.Exec("DELETE FROM sessions;")
	db.Exec("DELETE FROM applications;")
	db.Exec("DELETE FROM billings;")
	db.Exec("DELETE FROM forgot_passwords;")
	db.Exec("DELETE FROM connected_accounts;")
	db.Exec("DELETE FROM users;")
	db.Exec("DELETE FROM accounts;")
	
	// Reset autoincrement sequences for SQLite
	db.Exec("DELETE FROM sqlite_sequence;")
}


/* End File */
