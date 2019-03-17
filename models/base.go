//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"flag"
	"go/build"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"
)

//
// Start up the controller.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/skyclerk.com/.env")
}

//
// NewDB Setup the db connection.
//
func NewDB() (*DB, error) {

	dbName := os.Getenv("DB_DATABASE")

	// Is this a testing run?
	if flag.Lookup("test.v") != nil {
		dbName = os.Getenv("DB_DATABASE") + "_testing"
	}

	// Connect to Mysql
	db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@"+os.Getenv("DB_HOST")+"/"+dbName+"?charset=utf8&parseTime=true&loc=UTC")

	if err != nil {
		return nil, err
	}

	// Ping make sure the server is up.
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	// Run migrations
	db.AutoMigrate(&LabelsToLedger{}) // Must be first.
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Application{})
	db.AutoMigrate(&Label{})
	db.AutoMigrate(&Ledger{})
	db.AutoMigrate(&Contact{})
	db.AutoMigrate(&Category{})

	// Is this a testing run?
	if flag.Lookup("test.v") != nil {
		ClearTestingData(db)
	}

	// Return db connection.
	return &DB{db}, nil
}

//
// Clear testing data.
//
func ClearTestingData(db *gorm.DB) {
	// Clear tables
	db.Exec("TRUNCATE TABLE Users;")
	db.Exec("TRUNCATE TABLE Applications;")
	db.Exec("TRUNCATE TABLE GoSessions;")
	db.Exec("TRUNCATE TABLE Labels;")
	db.Exec("TRUNCATE TABLE Ledger;")
	db.Exec("TRUNCATE TABLE Contacts;")
	db.Exec("TRUNCATE TABLE Categories;")
	db.Exec("TRUNCATE TABLE LabelsToLedger;")
}

/* End File */
