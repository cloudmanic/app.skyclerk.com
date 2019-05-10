//
// Date: 5/9/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"testing"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/nbio/st"
)

//
// TestGetCurrentYearPnL01 - return the current year and the profit and lost of that year.
//
func TestGetCurrentYearPnL01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

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

	// Run test function
	pl := GetCurrentYearPnL(db, 33, 2019)

	// Figure out our own P&L
	total := 0.00

	for key := range dMap {
		if dMap[key].Date.Format("2006") == "2019" {
			total = total + dMap[key].Amount
		}
	}

	// Test results
	st.Expect(t, pl.Year, 2019)
	st.Expect(t, pl.Value, helpers.Round(total, 2))

	// ---------- Test empty year ------------- //

	// Run test function
	pl2 := GetCurrentYearPnL(db, 33, 2005)

	// Test results
	st.Expect(t, pl2.Year, 2005)
	st.Expect(t, pl2.Value, 0.00)
}
