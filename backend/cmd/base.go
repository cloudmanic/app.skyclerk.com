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
	"time"

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

	// Create a new application from the CLI
	case "create-application":
		actions.CreateApplication(db, *name)
		return true

	// // Loop through the accounts table and append "accounts" to the file name.
	// case "files-add-account-prefix":
	// 	FileAddAccountPrefix(db)
	// 	return true

	// Loop through the contacts table and build an avatar for every contact
	case "contacts-build-missing-avatars":
		err := db.GenerateAvatarsForAllMissing()
		if err != nil {
			panic(err)
		}
		return true

	// // Copy app log
	// case "copy-app-log":
	// 	CopyAppLog(db)
	// 	return true

	// Update billing table.
	case "billing-update-entries":
		UpdateBillingEntries(db)
		return true

	// Just a test
	case "test":
		fmt.Println("CMD Works....")
		return true

	}

	return false
}

//
// UpdateBillingEntries - Build script to make sure everyone has a billing account entry.
//
// go run main.go -cmd=billing-update-entries
//
func UpdateBillingEntries(db models.Datastore) {
	// Get all accounts
	accounts := []models.Account{}
	db.New().Where("billing_id = 0").Find(&accounts)

	length := len(accounts)

	for key, row := range accounts {

		fmt.Println(key, " / ", length)

		// See if we already have a billing profile
		g := models.Account{}
		db.New().Where("id = ? AND billing_id > 0", row.Id).Find(&g)
		if row.Id == g.Id {
			continue
		}

		now := time.Now()
		tExpire := now.Add(time.Hour * 24 * 85)

		// Setup the billing profile for this account.
		bp := models.Billing{
			Status:      "Trial",
			TrialExpire: tExpire,
		}
		db.New().Save(&bp)

		// Get all accounts this user owns
		acts := []models.Account{}
		db.New().Where("owner_id = ?", row.OwnerId).Find(&acts)

		for _, row2 := range acts {
			row2.BillingId = bp.Id
			db.New().Save(&row2)
		}

	}

}

// //
// // CopyAppLog - Copy applog stuff over to our new activities table - Delete one we run in production.
// //
// // go run main.go -cmd=copy-app-log
// //
// func CopyAppLog(db models.Datastore) {
//
// 	type Applog struct {
// 		ApplogId        uint      `gorm:"column:ApplogId"`
// 		ApplogAccountId uint      `gorm:"column:ApplogAccountId"`
// 		ApplogAction    string    `gorm:"column:ApplogAction"`
// 		ApplogText      string    `gorm:"column:ApplogText"`
// 		ApplogPoster    int       `gorm:"column:ApplogPoster"`
// 		ApplogCreatedAt time.Time `gorm:"column:ApplogCreatedAt"`
// 	}
//
// 	a := []Applog{}
//
// 	db.New().Table("Applog").Find(&a)
//
// 	for key, row := range a {
// 		action := row.ApplogAction
//
// 		if action == "Receipt" {
// 			action = "snapclerk"
// 		}
//
// 		if action == "Income" {
// 			action = "income"
// 		}
//
// 		if action == "Expense" {
// 			action = "expense"
// 		}
//
// 		t := models.Activity{
// 			CreatedAt: row.ApplogCreatedAt,
// 			AccountId: row.ApplogAccountId,
// 			Action:    action,
// 			SubAction: "other",
// 			Name:      row.ApplogText,
// 		}
//
// 		db.New().Save(&t)
//
// 		// Append accounts and save to DB.
// 		fmt.Println(key, " of ", len(a))
// 	}
//
// }
//
// //
// // FileAddAccountPrefix - Once we deploy GO based skyclerk we can deleted this function.
// //
// func FileAddAccountPrefix(db models.Datastore) {
// 	// Query and get files.
// 	files := []models.File{}
// 	db.New().Find(&files)
//
// 	// Loop through files and append
// 	for key, row := range files {
// 		if strings.Contains(row.Path, "accounts/") {
// 			continue
// 		}
//
// 		if row.Host != "amazon-s3" {
// 			continue
// 		}
//
// 		// Append accounts and save to DB.
// 		fmt.Println(row.Id, " - ", key, " - ", row.Path)
//
// 		row.Path = "accounts/" + row.Path
// 		row.ThumbPath = "accounts/" + row.ThumbPath
// 		db.New().Save(&row)
// 	}
//
// }

/* End File */
