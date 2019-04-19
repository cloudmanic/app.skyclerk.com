//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cmd

import (
	"flag"
	"fmt"

	"github.com/cloudmanic/app.skyclerk.com/backend/cmd/actions"
	"github.com/cloudmanic/app.skyclerk.com/backend/models"
)

//
// Run this and see if we have any commands to run.
//
func Run(db *models.DB) bool {

	// Grab flags
	action := flag.String("cmd", "none", "")
	file := flag.String("file", "", "")
	accountId := flag.Int("account_id", 0, "An account id.")
	flag.Parse()

	switch *action {

	// Import a CSV from AirBnb
	case "airbnb-import":
		actions.AirBnbImport(db, uint(*accountId), *file)
		return true
		break

	// Just a test
	case "test":
		fmt.Println("CMD Works....")
		return true
		break

	}

	return false
}

/* End File */
