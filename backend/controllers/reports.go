//
// Date: 4/14/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/reports"
)

//
// ReportsIncomeByContact - Get income by contact
//
func (t *Controller) ReportsIncomeByContact(c *gin.Context) {
	// Set start / end big range default
	start := helpers.ParseDateNoError(c.DefaultQuery("start", "1800-01-01"))
	end := helpers.ParseDateNoError(c.DefaultQuery("end", "3000-01-01"))

	// Run function
	result := reports.GetIncomeByContact(t.db, uint(c.MustGet("accountId").(int)), start, end, c.DefaultQuery("sort", "desc"))

	// Return happy JSON
	c.JSON(200, result)
}

//
// ReportsExpensesByContact - Get expenses by contact
//
func (t *Controller) ReportsExpensesByContact(c *gin.Context) {
	// Set start / end big range default
	start := helpers.ParseDateNoError(c.DefaultQuery("start", "1800-01-01"))
	end := helpers.ParseDateNoError(c.DefaultQuery("end", "3000-01-01"))

	// Run function
	result := reports.GetExpenseByContact(t.db, uint(c.MustGet("accountId").(int)), start, end, c.DefaultQuery("sort", "desc"))

	// Return happy JSON
	c.JSON(200, result)
}

//
// ReportsPnl returns income, expense, profit by date range grouping
//
func (t *Controller) ReportsPnl(c *gin.Context) {
	// Set start / end big range default
	start := helpers.ParseDateNoError(c.DefaultQuery("start", "1800-01-01"))
	end := helpers.ParseDateNoError(c.DefaultQuery("end", "3000-01-01"))

	// Run function
	pl := reports.GetPnL(t.db, uint(c.MustGet("accountId").(int)), start, end, c.DefaultQuery("group", "month"), c.DefaultQuery("sort", "desc"))

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
