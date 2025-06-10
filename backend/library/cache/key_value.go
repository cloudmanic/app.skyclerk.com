//
// Date: 2/27/2017
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cache

import (
	"encoding/json"
	"time"

	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

var (
	db *models.DB
)

//
// StartCache - startup SQLite Cache
//
func StartCache() {
	// Setup database connection
	var err error
	db, err = models.NewDB()

	if err != nil {
		services.Fatal(err)
	}
}

//
// SetDB - Set the database for cache (used in testing)
//
func SetDB(database *models.DB) {
	db = database
}

//
// Delete a stored key. We do not
// return error as if the database is not up we should
// have a massive fail.
//
func Delete(key string) {
	// Delete from cache
	err := db.CacheDelete(key)

	if err != nil {
		services.Fatal(err)
	}
}

//
// Set Store a key value into cache. We do not
// return error as if the database is not up we should
// have a massive fail.
//
func Set(key string, value interface{}) {
	// Pack up into JSON
	b, err := json.Marshal(&value)

	if err != nil {
		services.Fatal(err)
	}

	// Store in cache
	err = db.CacheSet(key, string(b))

	if err != nil {
		services.Fatal(err)
	}
}

//
// Store a key value into cache that expires. We do not
// return error as if the database is not up we should
// have a massive fail.
//
func SetExpire(key string, expire time.Duration, value interface{}) {
	// Pack up into JSON
	b, err := json.Marshal(&value)

	if err != nil {
		services.Fatal(err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(expire)

	// Store in cache
	err = db.CacheSetWithExpiration(key, string(b), expiresAt)

	if err != nil {
		services.Fatal(err)
	}
}

//
// Get key from cache
//
func Get(key string, result interface{}) (bool, error) {
	value, err := db.CacheGet(key)

	// Does not exist
	if err != nil {
		return false, err
	}

	// UnMarshal Result
	err = json.Unmarshal([]byte(value), result)

	if err != nil {
		services.Fatal(err)
		return false, err
	}

	return true, nil
}

/* End File */
