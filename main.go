//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import (
	"github.com/cloudmanic/skyclerk.com/cmd"
	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/cloudmanic/skyclerk.com/services"
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

}

/* End File */
