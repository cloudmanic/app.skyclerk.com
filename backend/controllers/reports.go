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
// ReportsPnl returns income, expense, profit by date range grouping
//
func (t *Controller) ReportsPnl(c *gin.Context) {
	// Set start / end
	start := dates.ParseDateNoError("2019-03-01")
	end := dates.ParseDateNoError("2019-06-30")

	// Run function
	pl := reports.GetPnL(t.db, uint(c.MustGet("accountId").(int)), start, end, "month")

	// Return happy JSON
	c.JSON(200, pl)
}

//
// ReportsCurrentPnl returns current year and the P&L for that year.
//
func (t *Controller) ReportsCurrentPnl(c *gin.Context) {
	// Run function
	pl := reports.GetCurrentYearPnL(t.db, uint(c.MustGet("accountId").(int)), time.Now().Year())

	// Return happy JSON
	c.JSON(200, pl)
}

/* End File */
