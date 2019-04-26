//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"

	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

var (
	defaultLimit int = 100
)

type Controller struct {
	db models.Datastore
}

type ValidateRequest interface {
	Validate(models.Datastore, string, uint, uint, uint) error
}

//
// Set the database.
//
func (t *Controller) SetDB(db models.Datastore) {
	t.db = db
}

//
// Validate and Create object.
//
func (t *Controller) ValidateRequest(c *gin.Context, obj ValidateRequest, action string) error {
	// Bind the JSON that got sent into an object and validate.
	if err := c.ShouldBindJSON(obj); err == nil {

		// Get user id.
		userId := uint(c.MustGet("userId").(int))

		// AccountId.
		accountId := uint(c.MustGet("accountId").(int))

		// If the action is update add the id
		var id uint = 0

		if action == "update" {
			id2, err := strconv.ParseInt(c.Param("id"), 10, 32)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"errors": err})
				return err
			}

			id = uint(id2)
		}

		// Run validation
		err := obj.Validate(t.db, action, userId, accountId, id)

		// If we had validation errors return them and do no more.
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return err
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in body. There is a chance the JSON maybe valid but does not match the data type requirements. For example maybe you passed a string in for an integer."})
		return err
	}

	return nil
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
			if (origin == os.Getenv("SITE_URL")) || strings.Contains(origin, "localhost") || strings.Contains(origin, os.Getenv("LOCAL_IP")) {
				return true
			}
			services.LogInfo("Failed CORS request - origin: " + origin)
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	// Add custom parms validating
	router.Use(t.ParamValidateMiddleware())

	// Register Routes
	t.DoRoutes(router)

	// Setup http server
	if os.Getenv("APP_ENV") == "local" {
		srv := &http.Server{
			Handler:      router,
			Addr:         ":" + os.Getenv("HTTP_PORT"),
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
		}

		// Log this start.
		services.LogInfo("Starting web server at " + os.Getenv("SITE_URL"))

		// Start server and log if fails
		log.Fatal(srv.ListenAndServe())
	} else {
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/letsencrypt/"),
			Email:      "help@skyclerk.com",
			HostPolicy: autocert.HostWhitelist(os.Getenv("SITE_DOMAIN")),
		}

		log.Printf("Starting secure server at " + os.Getenv("SITE_URL"))

		log.Fatal(autotls.RunWithManager(router, m))
	}

}

/* End File */
