//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"app.skyclerk.com/backend/cmd"
	"app.skyclerk.com/backend/controllers"
	"app.skyclerk.com/backend/cron"
	"app.skyclerk.com/backend/library/cache"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// Main...
//
func main() {
	// Start the db connection.
	db, err := models.NewDB()

	if err != nil {
		services.LogFatal(err)
	}

	// Close db when this app dies. (This might be useless)
	defer db.Close()

	// Start cache
	cache.StartCache()

	// See if this a command. If so run the command and do not start the app.
	status := cmd.Run(db)

	if status == true {
		return
	}

	// ----------- Start Web Server ------------- //

	// Startup controller
	c := &controllers.Controller{}

	// Set the database the controller uses.
	c.SetDB(db)

	// Start the cron process
	go cron.Start(db)

	// Start webserver & controllers
	c.StartWebServer()
}

/* End File */
