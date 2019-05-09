//
// Date: 4/14/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/reports"
)

//
// ReportsCurrentPnl returns current year and the P&L for that year.
//
func (t *Controller) ReportsCurrentPnl(c *gin.Context) {
	// Run test function
	pl := reports.GetCurrentYearPnL(t.db, uint(c.MustGet("accountId").(int)), time.Now().Year())

	// Return happy JSON
	c.JSON(200, pl)
}
