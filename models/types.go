//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

// Used for date formatting with json conversion.
type Date struct {
	time.Time
}

// Convert JSON string to a date. Format XXXX-XX-XX
//
func (t *Date) UnmarshalJSON(b []byte) error {

	// Remove quotes
	str := strings.Replace(string(b), "\"", "", -1)

	// Parse string
	tt, _ := time.Parse("2006-01-02", str)

	// Return UTC
	*t = Date{tt.UTC()}
	return nil
}

//
// Convert JSON string to a date. Make it so when we create JSON we return this format XXXX-XX-XX
//
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))), nil
}

//
// Format time Date Going into the DB
//
func (t Date) Value() (driver.Value, error) {
	return t.Format("2006-01-02"), nil
}

//
// Convert to type Date Coming out of the DB
//
func (t *Date) Scan(value interface{}) error {

	// Parse string
	tt, _ := time.Parse("2006-01-02 03:04:05 -0700 MST", fmt.Sprintf("%s", value))

	// Return UTC
	*t = Date{tt.UTC()}
	return nil
}

/* End File */
