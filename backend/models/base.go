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
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"
	_ "github.com/mattn/go-sqlite3"
)

// Start up the model.
func init() {
	// Only load .env file if we're not in a test environment
	if flag.Lookup("test.v") == nil {
		// Get the path to the .env file relative to this source file
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		envPath := filepath.Join(basepath, "..", ".env")
		env.ReadEnv(envPath)
	}

	// Set defaults for testing when environment variables are not set
	// Force APP_ENV to test during tests to ensure consistent behavior
	if flag.Lookup("test.v") != nil {
		os.Setenv("APP_ENV", "test")
	} else {
		setDefaultIfEmpty("APP_ENV", "test")
	}
	setDefaultIfEmpty("STRIPE_CLIENT_ID", "ca_test_default")
	setDefaultIfEmpty("STRIPE_SECRET_KEY", "sk_test_default")
	setDefaultIfEmpty("APP_URL", "http://localhost:8080")
	setDefaultIfEmpty("SITE_URL", "http://localhost:4200")
	setDefaultIfEmpty("ENCRYPTION_KEY", "test-encryption-key-32-characters")
	setDefaultIfEmpty("POSTMARK_SERVER_KEY", "test-postmark-server-key")
	setDefaultIfEmpty("POSTMARK_ACCOUNT_KEY", "test-postmark-account-key")
	setDefaultIfEmpty("MAILGUN_DOMAIN", "test.example.com")
	setDefaultIfEmpty("MAILGUN_API_KEY", "test-mailgun-key")
	setDefaultIfEmpty("SENDY_API_KEY", "test-sendy-key")
	setDefaultIfEmpty("SLACK_URL", "")
	setDefaultIfEmpty("OBJECT_BASE_URL", "http://127.0.0.1:9000")
	setDefaultIfEmpty("OBJECT_ACCESS_KEY_ID", "test_access_key")
	setDefaultIfEmpty("OBJECT_SECRET_ACCESS_KEY", "test_secret_key")
	setDefaultIfEmpty("OBJECT_ENDPOINT", "127.0.0.1:9000")
	setDefaultIfEmpty("OBJECT_BUCKET", "test-bucket")
	setDefaultIfEmpty("CACHE_DIR", "/tmp/skyclerk-cache")
	// Set font path relative to this source file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	fontPath := filepath.Join(basepath, "..", "..", "fonts")
	setDefaultIfEmpty("FONT_PATH", fontPath)
}

// setDefaultIfEmpty sets an environment variable to a default value if it's not already set
func setDefaultIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}

// NewDB Setup the db connection.
func NewDB() (*DB, error) {
	var db *gorm.DB
	var err error

	// Set up SQLite database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Create default path in cache/sqlite directory
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		cacheDir := filepath.Join(basepath, "..", "cache", "sqlite")

		// Create directory if it doesn't exist
		os.MkdirAll(cacheDir, 0755)

		dbPath = filepath.Join(cacheDir, "skyclerk.db")

		// Is this a testing run?
		if flag.Lookup("test.v") != nil {
			dbPath = filepath.Join(cacheDir, "skyclerk_testing.db")
		}
	}

	// Connect to SQLite
	db, err = gorm.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	// Set GORM log mode to silent to suppress AutoMigrate warnings
	db.LogMode(false)

	// Ping make sure the server is up.
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	// Set max connections
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Run migrations
	doMigrations(db)

	// Return db connection.
	return &DB{db}, nil
}

// doMigrations - Run our migrations
func doMigrations(db *gorm.DB) {
	db.AutoMigrate(&LabelsToLedger{}) // Must be first.
	db.AutoMigrate(&FilesToLedger{})  // Must be first.
	db.AutoMigrate(&Activity{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Application{})
	db.AutoMigrate(&Label{})
	db.AutoMigrate(&Ledger{})
	db.AutoMigrate(&Contact{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&SnapClerk{})
	db.AutoMigrate(&Invite{})
	db.AutoMigrate(&Billing{})
	db.AutoMigrate(&ForgotPassword{})
	db.AutoMigrate(&ConnectedAccounts{})
	db.AutoMigrate(&Cache{})
}

/* End File */
