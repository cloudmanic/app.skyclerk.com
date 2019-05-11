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
	"strings"

	"app.skyclerk.com/backend/cmd/actions"
	"app.skyclerk.com/backend/models"
)

//
// Run this and see if we have any commands to run.
//
func Run(db models.Datastore) bool {

	// Grab flags
	action := flag.String("cmd", "none", "")
	file := flag.String("file", "", "")
	accountId := flag.Int("account_id", 0, "An account id.")
	name := flag.String("name", "", "")
	flag.Parse()

	switch *action {

	// Import a CSV from AirBnb
	case "airbnb-import":
		actions.AirBnbImport(db, uint(*accountId), *file)
		return true
		break

	// Create a new application from the CLI
	case "create-application":
		actions.CreateApplication(db, *name)
		return true
		break

	// Loop through the accounts table and append "accounts" to the file name.
	case "files-add-account-prefix":
		FileAddAccountPrefix(db)
		return true
		break

	// Loop through the contacts table and build an avatar for every contact
	case "contacts-build-missing-avatars":
		err := db.GenerateAvatarsForAllMissing()
		if err != nil {
			panic(err)
		}
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

//
// FileAddAccountPrefix - Once we deploy GO based skyclerk we can deleted this function.
//
func FileAddAccountPrefix(db models.Datastore) {
	// Query and get files.
	files := []models.File{}
	db.New().Find(&files)

	// Loop through files and append
	for key, row := range files {
		if strings.Contains(row.Path, "accounts/") {
			continue
		}

		if row.Host != "amazon-s3" {
			continue
		}

		// Append accounts and save to DB.
		fmt.Println(row.Id, " - ", key, " - ", row.Path)

		row.Path = "accounts/" + row.Path
		row.ThumbPath = "accounts/" + row.ThumbPath
		db.New().Save(&row)
	}

}

/* End File */
