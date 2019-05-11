//
// Date: 3/1/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"fmt"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/models"
)

//
// Create a new application.
//
// go run main.go -cmd=create-application -name="Ionic App"
//
func CreateApplication(db models.Datastore, name string) {

	// Generate a random string for the client id.
	clientId, err := helpers.GenerateRandomString(15)

	if err != nil {
		panic(err)
	}

	// Setup the application
	app := models.Application{Name: name, ClientId: clientId, GrantType: "password"}

	// Create new application
	db.New().Save(&app)

	fmt.Println("Success Application Id: ", app.Id, " ClientId: "+app.ClientId)
}

/* End File */
