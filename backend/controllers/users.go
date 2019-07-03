//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/response"
)

//
// GetUsers - Return a list of users. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetUsers(c *gin.Context) {
	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Query and get users
	users := t.db.GetUsersByAccount(accountId)

	// Return happy.
	response.Results(c, users, nil)
}

/* End File */
