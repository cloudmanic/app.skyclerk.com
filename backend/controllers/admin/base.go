//
// Date: 11/3/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import "app.skyclerk.com/backend/models"

// Controller struct
type Controller struct {
	db models.Datastore
}

//
// Set the database.
//
func (t *Controller) SetDB(db models.Datastore) {
	t.db = db
}

/* End File */
