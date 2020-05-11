//
// Date: 2020-05-11
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package account

import (
	"time"

	"app.skyclerk.com/backend/library/sendy"
	"app.skyclerk.com/backend/library/slack"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

// //
// // Clear expired sessions (access tokens)
// //
// func ClearExpiredSessions(db *models.DB) {
//
// 	// Find the Centcom app
// 	centcomApp := models.Application{}
// 	db.New().Where("name <= ?", "Centcom").Find(&centcomApp)
//
// 	if centcomApp.Id > 0 {
//
// 		// Delete expired centcom sessions
// 		db.New().Where("last_activity <= ? AND application_id = ?", time.Now().AddDate(0, 0, -1), centcomApp.Id).Delete(&models.Session{})
//
// 		// Just cleared Centcom sessions
// 		services.InfoMsg("Centcom sessions cleared.")
//
// 	}
//
// 	// Find the Personal app
// 	personalApp := models.Application{}
// 	db.New().Where("name <= ?", "Personal Access Token").Find(&personalApp)
//
// 	if personalApp.Id > 0 {
//
// 		// Clear all sessions that have not had activity in the last 14 days (2 weeks)
// 		db.New().Where("last_activity <= ? AND application_id != ?", time.Now().AddDate(0, 0, -14), personalApp.Id).Delete(&models.Session{})
//
// 		// Log cleared sessions.
// 		services.InfoMsg("All expired sessions cleared.")
//
// 	}
//
// }

//
// ExpireTrails will update expired accounts
//
func ExpireTrails(db models.Datastore) {
	// Get all billing entries.
	billings := []models.Billing{}
	db.New().Where("trial_expire <= ? AND status = ? AND stripe_subscription = ?", time.Now(), "Trial", "").Find(&billings)

	// Loop through all the billing entries.
	for _, row := range billings {
		// Must have an ID
		if row.Id <= 0 {
			continue
		}

		// Get the account by billing
		account := models.Account{}
		db.New().Where("billing_id = ?", row.Id).Find(&account)

		// Get account owner
		owner := models.User{}
		db.New().Find(&owner, account.OwnerId)

		// Expire account.
		row.Status = "Expired"
		db.New().Save(&row)
		services.InfoMsg("Free trial has just expired : " + owner.Email)

		if len(owner.Email) > 0 {
			go slack.Notify("#events", "Skyclerk User Free Trial Expired : "+owner.Email)
			go sendy.Subscribe("expired", owner.Email, owner.FirstName, owner.LastName, "")
			go sendy.Unsubscribe("trial", owner.Email)
		}
	}

	services.InfoMsg("All expire trails set to expired.")
}

/* End File */
