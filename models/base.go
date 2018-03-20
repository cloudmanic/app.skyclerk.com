//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
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
	// Helpful for testing
	if flag.Lookup("test.v") != nil {
		env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/skyclerk.com/.env")
	}
}

//
// Setup the db connection.
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

	// Enable
	//db.LogMode(true)
	//db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	// Run migrations
	//db.AutoMigrate(&Unit{})

	// Is this a testing run? If so load testing data.
	if flag.Lookup("test.v") != nil {
		LoadTestingData(db)
	}

	// Return db connection.
	return &DB{db}, nil
}

//
// Load testing data.
//
func LoadTestingData(db *gorm.DB) {

	// Shared time we use.
	//ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
	//ds := Date{time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)}

}

/* End File */
