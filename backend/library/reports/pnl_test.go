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
// TestGetPnL01 - return PnL by group / start / end
//
func TestGetPnL01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 50 ledger entries.
	for i := 0; i < 50; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for March
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-03-01")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for April
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-04-01")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for May
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-05-01")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for May
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-06-01")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Set start / end
	start := helpers.ParseDateNoError("2019-03-01")
	end := helpers.ParseDateNoError("2019-06-30")

	// Run test function
	pl := GetPnL(db, 33, start, end, "month", "ASC")

	// Figure out our own P&L
	profit032019 := 0.00
	profit042019 := 0.00
	profit052019 := 0.00
	profit062019 := 0.00

	income032019 := 0.00
	income042019 := 0.00
	income052019 := 0.00
	income062019 := 0.00

	expense032019 := 0.00
	expense042019 := 0.00
	expense052019 := 0.00
	expense062019 := 0.00

	for key := range dMap {
		if dMap[key].Date.Format("2006-01") == "2019-03" {
			profit032019 = profit032019 + dMap[key].Amount

			if dMap[key].Amount > 0 {
				income032019 = income032019 + dMap[key].Amount
			}

			if dMap[key].Amount < 0 {
				expense032019 = expense032019 + dMap[key].Amount
			}
		}

		if dMap[key].Date.Format("2006-01") == "2019-04" {
			profit042019 = profit042019 + dMap[key].Amount

			if dMap[key].Amount > 0 {
				income042019 = income042019 + dMap[key].Amount
			}

			if dMap[key].Amount < 0 {
				expense042019 = expense042019 + dMap[key].Amount
			}
		}

		if dMap[key].Date.Format("2006-01") == "2019-05" {
			profit052019 = profit052019 + dMap[key].Amount

			if dMap[key].Amount > 0 {
				income052019 = income052019 + dMap[key].Amount
			}

			if dMap[key].Amount < 0 {
				expense052019 = expense052019 + dMap[key].Amount
			}
		}

		if dMap[key].Date.Format("2006-01") == "2019-06" {
			profit062019 = profit062019 + dMap[key].Amount

			if dMap[key].Amount > 0 {
				income062019 = income062019 + dMap[key].Amount
			}

			if dMap[key].Amount < 0 {
				expense062019 = expense062019 + dMap[key].Amount
			}
		}
	}

	// Test results
	st.Expect(t, len(pl), 4)
	st.Expect(t, helpers.Round(pl[0].Profit, 2), helpers.Round(profit032019, 2))
	st.Expect(t, helpers.Round(pl[1].Profit, 2), helpers.Round(profit042019, 2))
	st.Expect(t, helpers.Round(pl[2].Profit, 2), helpers.Round(profit052019, 2))
	st.Expect(t, helpers.Round(pl[3].Profit, 2), helpers.Round(profit062019, 2))

	st.Expect(t, helpers.Round(pl[0].Income, 2), helpers.Round(income032019, 2))
	st.Expect(t, helpers.Round(pl[1].Income, 2), helpers.Round(income042019, 2))
	st.Expect(t, helpers.Round(pl[2].Income, 2), helpers.Round(income052019, 2))
	st.Expect(t, helpers.Round(pl[3].Income, 2), helpers.Round(income062019, 2))

	st.Expect(t, helpers.Round(pl[0].Expense, 2), helpers.Round(expense032019, 2))
	st.Expect(t, helpers.Round(pl[1].Expense, 2), helpers.Round(expense042019, 2))
	st.Expect(t, helpers.Round(pl[2].Expense, 2), helpers.Round(expense052019, 2))
	st.Expect(t, helpers.Round(pl[3].Expense, 2), helpers.Round(expense062019, 2))

	// ---------- Different sort ------------- //

	// Run test function
	pl2 := GetPnL(db, 33, start, end, "month", "DESC")

	// Test results
	st.Expect(t, len(pl), 4)
	st.Expect(t, helpers.Round(pl2[3].Profit, 2), helpers.Round(profit032019, 2))
	st.Expect(t, helpers.Round(pl2[2].Profit, 2), helpers.Round(profit042019, 2))
	st.Expect(t, helpers.Round(pl2[1].Profit, 2), helpers.Round(profit052019, 2))
	st.Expect(t, helpers.Round(pl2[0].Profit, 2), helpers.Round(profit062019, 2))

	st.Expect(t, helpers.Round(pl2[3].Income, 2), helpers.Round(income032019, 2))
	st.Expect(t, helpers.Round(pl2[2].Income, 2), helpers.Round(income042019, 2))
	st.Expect(t, helpers.Round(pl2[1].Income, 2), helpers.Round(income052019, 2))
	st.Expect(t, helpers.Round(pl2[0].Income, 2), helpers.Round(income062019, 2))

	st.Expect(t, helpers.Round(pl2[3].Expense, 2), helpers.Round(expense032019, 2))
	st.Expect(t, helpers.Round(pl2[2].Expense, 2), helpers.Round(expense042019, 2))
	st.Expect(t, helpers.Round(pl2[1].Expense, 2), helpers.Round(expense052019, 2))
	st.Expect(t, helpers.Round(pl2[0].Expense, 2), helpers.Round(expense062019, 2))
}

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
