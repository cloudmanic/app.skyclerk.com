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
	"os/exec"
	"strings"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/services"
	"github.com/jinzhu/gorm"
)

const (
	dockerMysqlContainerName = "skyclerk_com_testing"
)

//
// NewTestDB - Start the Test DB connection.
//
func NewTestDB(dbName string) (*DB, string, error) {
	var err error

	// Make sure our docker mysql container for testing is running.
	if !isDockerMysqlRunning() {
		log.Fatal(errors.New("Docker testing Mysql container is not running. Please run scripts/start_testing_db.sh."))
	}

	// If dbName is empty we create our own.
	if len(dbName) == 0 {
		dbName = "sk_" + helpers.RandStr(10)
	}

	// Create database
	createTestDatabase(dbName)

	// Connect to Mysql but do not connect to a database
	db, err := gorm.Open("mysql", "root:foobar@tcp(127.0.0.1:9907)/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
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
	// special case for testing_db
	if dbName == "testing_db" {
		db.Close()
		return
	}

	// Drop table
	db.Exec("DROP DATABASE IF EXISTS " + dbName + ";")
	db.Close()
}

//
// TruncateAllTables - Clear all tables. DONT FORGET TO UPDATE ClearAccount() in accounts.go
//
func truncateAllTables(db *gorm.DB) {
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
	db.Exec("TRUNCATE TABLE acct_to_billings;")
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
