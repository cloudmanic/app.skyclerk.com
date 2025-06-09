//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/services"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dockerMysqlContainerName = "mariadb-10.20-testing.app.skyclerk.com"
)

//
// NewTestDB - Start the Test DB connection.
//
func NewTestDB(dbName string) (*DB, string, error) {
	var db *gorm.DB
	var err error

	// Check which database driver to use for testing
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite3" // Default to SQLite for testing
	}

	switch driver {
	case "sqlite3":
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

	case "mysql":
		// Make sure our docker mysql container for testing is running.
		if !isDockerMysqlRunning() {
			log.Fatal(errors.New("Docker testing Mysql container is not running. Please run scripts/start_testing_db.sh or run the docker-compose.yml."))
		}

		// If dbName is empty we create our own.
		if len(dbName) == 0 {
			dbName = "sk_" + helpers.RandStr(10)
		}

		// Create database
		createTestDatabase(dbName)

		// Connect to Mysql but do not connect to a database
		db, err = gorm.Open("mysql", "root:foobar@tcp(127.0.0.1:9907)/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

		if err != nil {
			services.Error(err)
			log.Fatal(err)
		}

	default:
		log.Fatal(errors.New("Unsupported database driver for testing: " + driver))
	}

	// Run doMigrations
	doMigrations(db)

	// Clear all tables.
	truncateAllTables(db)

	// Return db connection.
	return &DB{db}, dbName, nil
}

//
// createTestDatabase - Create test database.
//
func createTestDatabase(name string) {
	// Connect to DB
	db, err := sql.Open("mysql", "root:foobar@tcp(127.0.0.1:9907)/?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create DB
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)

	if err != nil {
		log.Fatal(err)
	}
}

//
// TestingTearDown - used to delete database we created
//
func TestingTearDown(db *DB, dbName string) {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite3"
	}

	switch driver {
	case "sqlite3":
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
		
	case "mysql":
		// special case for testing_db
		if dbName == "testing_db" {
			db.Close()
			return
		}

		// Drop table
		db.Exec("DROP DATABASE IF EXISTS " + dbName + ";")
		db.Close()
	}
}

//
// TruncateAllTables - Clear all tables. DONT FORGET TO UPDATE ClearAccount() in accounts.go
//
func truncateAllTables(db *gorm.DB) {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite3"
	}

	// SQLite doesn't support TRUNCATE, so we use DELETE
	if driver == "sqlite3" {
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
	} else {
		// MySQL TRUNCATE statements
		db.Exec("TRUNCATE TABLE accounts;")
		db.Exec("TRUNCATE TABLE users;")
		db.Exec("TRUNCATE TABLE activities;")
		db.Exec("TRUNCATE TABLE applications;")
		db.Exec("TRUNCATE TABLE sessions;")
		db.Exec("TRUNCATE TABLE acct_to_users;")
		db.Exec("TRUNCATE TABLE invites;")
		db.Exec("TRUNCATE TABLE Labels;")
		db.Exec("TRUNCATE TABLE Ledger;")
		db.Exec("TRUNCATE TABLE Files;")
		db.Exec("TRUNCATE TABLE Contacts;")
		db.Exec("TRUNCATE TABLE Categories;")
		db.Exec("TRUNCATE TABLE SnapClerk;")
		db.Exec("TRUNCATE TABLE LabelsToLedger;")
		db.Exec("TRUNCATE TABLE FilesToLedger;")
		db.Exec("TRUNCATE TABLE billings;")
		db.Exec("TRUNCATE TABLE forgot_passwords;")
		db.Exec("TRUNCATE TABLE connected_accounts;")
	}
}

//
// isDockerMysqlRunning - verify our testing mysql instance is running in docker.
//
func isDockerMysqlRunning() bool {
	// Command to get the status of our mysql docker container.
	command := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Status}}|{{.Ports}}|{{.Names}}")
	output, err := command.CombinedOutput()
	if err != nil {
		return false
	}

	// Parse the output of the command
	outputString := string(output)
	outputString = strings.TrimSpace(outputString)
	dockerPsResponse := strings.Split(outputString, "\n")

	// Loop through the response to find the container we care about.
	for _, response := range dockerPsResponse {
		containerStatusData := strings.Split(response, "|")
		containerStatus := containerStatusData[1]
		containerName := containerStatusData[3]

		// dis we find the cotainer we wanted?
		if containerName == dockerMysqlContainerName {
			if strings.HasPrefix(containerStatus, "Up ") {
				return true
			}
		}
	}

	return false
}

/* End File */
