//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/models"
)

//
// GetSnapClerks - Get a list of snap clerks.
//
func (t *Controller) GetSnapClerks(c *gin.Context) {
	// SnapClerk to return
	results := []models.SnapClerk{}

	// Make query
	t.db.New().Preload("File").Where("SnapClerkStatus = ?", "Pending").Find(&results)

	// Loop through and add signed urls to files TODO(spicer): clean this up. Maybe move all this into the model.
	for key, row := range results {
		results[key].File.Url = t.db.GetSignedFileUrl(row.File.Path)
		results[key].File.Thumb600By600Url = t.db.GetSignedFileUrl(row.File.ThumbPath)
	}

	// Return happy JSON
	c.JSON(200, results)
}

/* End File */
