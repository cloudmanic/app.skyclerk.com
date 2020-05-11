//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package cron

import (
	"os"

	"github.com/robfig/cron"

	"app.skyclerk.com/backend/cron/account"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// Start will start our cront job.
//
func Start(db models.Datastore) {
	// Lets get started
	services.InfoMsg("Cron Started: " + os.Getenv("APP_ENV"))

	// Stuff we do on start as well
	account.ExpireTrails(db)
	//user.ClearExpiredSessions(db)

	// New Cron instance
	c := cron.New()

	// User clean up stuff
	c.AddFunc("@every 50m", func() { account.ExpireTrails(db) }) // Some reason 1h does not work.
	//c.AddFunc("@every 6h", func() { user.ClearExpiredSessions(db) })

	// System stuff.
	//c.AddFunc("@every 10s", func() { DatabasePing(db) })

	// Start cron service
	c.Run()
}

// //
// // DatabasePing is use this to keep the database alive.
// //
// func DatabasePing(db *models.Datastore) {
// 	// Just run a query to make sure things are active.
// 	a := []models.Application{}
// 	db.New().Find(&a)
// }

/* End File */
