//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"github.com/cloudmanic/app.skyclerk.com/backend/cmd"
	"github.com/cloudmanic/app.skyclerk.com/backend/controllers"
	"github.com/cloudmanic/app.skyclerk.com/backend/models"
	"github.com/cloudmanic/app.skyclerk.com/backend/services"
	_ "github.com/jpfuentes2/go-env/autoload"
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

	// Start webserver & controllers
	c.StartWebServer()
}

/* End File */
