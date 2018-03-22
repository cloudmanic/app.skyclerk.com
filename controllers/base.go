//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/cloudmanic/skyclerk.com/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	defaultLimit int = 100
)

type Controller struct {
	db models.Datastore
}

//
// Set the database.
//
func (t *Controller) SetDB(db models.Datastore) {
	t.db = db
}

//
// Start the webserver
//
func (t *Controller) StartWebServer() {

	// Set GIN Settings
	gin.SetMode("release")
	gin.DisableConsoleColor()

	// Set Router
	router := gin.New()

	// Logger - Global middleware
	if os.Getenv("HTTP_LOG_REQUESTS") == "true" {
		router.Use(gin.Logger())
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// CORS Middleware - Global middleware
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range,Range"},
		ExposeHeaders:    []string{"Content-Length", "X-Last-Page", "X-Offset", "X-Limit", "X-No-Limit-Count"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return (origin == os.Getenv("SITE_URL")) || strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	// Add custom parms validating
	router.Use(t.ParamValidateMiddleware())

	// Register Routes
	t.DoRoutes(router)

	// Setup http server
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + os.Getenv("HTTP_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Log this start.
	services.LogInfo("Starting web server at http://localhost:" + os.Getenv("HTTP_PORT"))

	// Start server and log if fails
	log.Fatal(srv.ListenAndServe())
}

/* End File */
