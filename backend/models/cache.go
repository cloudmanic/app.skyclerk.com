//
// Date: 2025-01-06
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2025 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

// Cache struct for SQLite-based caching
type Cache struct {
	Id        uint   `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string    `gorm:"unique;not null;index"`
	Value     string    `gorm:"type:text;not null"`
	ExpiresAt *time.Time `gorm:"index"`
}

//
// Set - Store a key-value pair in cache without expiration
//
func (db *DB) CacheSet(key string, value string) error {
	cache := Cache{
		Key:   key,
		Value: value,
	}
	
	// Use Create or Update
	if err := db.Where("key = ?", key).Assign(&cache).FirstOrCreate(&cache).Error; err != nil {
		return err
	}
	
	return nil
}

//
// SetWithExpiration - Store a key-value pair in cache with expiration
//
func (db *DB) CacheSetWithExpiration(key string, value string, expiresAt time.Time) error {
	cache := Cache{
		Key:       key,
		Value:     value,
		ExpiresAt: &expiresAt,
	}
	
	// Use Create or Update
	if err := db.Where("key = ?", key).Assign(&cache).FirstOrCreate(&cache).Error; err != nil {
		return err
	}
	
	return nil
}

//
// Get - Retrieve a value from cache by key
//
func (db *DB) CacheGet(key string) (string, error) {
	var cache Cache
	
	// First, clean up expired entries for this key
	db.cleanupExpiredCache(key)
	
	// Try to find the cache entry
	if db.Where("key = ?", key).First(&cache).RecordNotFound() {
		return "", errors.New("cache key not found")
	}
	
	// Check if expired (double check in case cleanup didn't catch it)
	if cache.ExpiresAt != nil && cache.ExpiresAt.Before(time.Now()) {
		db.CacheDelete(key)
		return "", errors.New("cache key expired")
	}
	
	return cache.Value, nil
}

//
// Delete - Remove a cache entry by key
//
func (db *DB) CacheDelete(key string) error {
	return db.Where("key = ?", key).Delete(&Cache{}).Error
}

//
// CleanupExpired - Remove all expired cache entries
//
func (db *DB) CacheCleanupExpired() error {
	return db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&Cache{}).Error
}

//
// cleanupExpiredCache - Internal method to cleanup expired entries for a specific key
//
func (db *DB) cleanupExpiredCache(key string) {
	db.Where("key = ? AND expires_at IS NOT NULL AND expires_at < ?", key, time.Now()).Delete(&Cache{})
}

/* End File */