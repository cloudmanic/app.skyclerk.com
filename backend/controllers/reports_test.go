//
// Date: 5/9/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/reports"
	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestReportsCurrentPnl01 - test getting the current year P&L
//
func TestReportsCurrentPnl01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Figure out our own P&L
	total := 0.00

	for key := range dMap {
		if dMap[key].Date.Format("2006") == "2019" {
			total = total + dMap[key].Amount
		}
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/pnl-current-year", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/pnl-current-year", c.ReportsCurrentPnl)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := reports.YearPnL{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, results.Year, 2019)
	st.Expect(t, results.Value, helpers.Round(total, 2))
}

/* End File */
