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
// TestGetLabelsPnL01 - Test Labels
//
func TestGetLabelsPnL01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		r := test.GetRandomLabel(33)
		r.Name = "Label #1"
		l.Labels = []models.Label{r}
		l.Date = helpers.ParseDateNoError("2019-03-01")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		r := test.GetRandomLabel(33)
		r.Name = "Label #2"
		l.Labels = []models.Label{r}
		l.Date = helpers.ParseDateNoError("2019-03-05")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = -100
		r := test.GetRandomLabel(33)
		r.Name = "Label #3"
		l.Labels = []models.Label{r}
		l.Date = helpers.ParseDateNoError("2019-03-05")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = -100
		r := test.GetRandomLabel(33)
		r.Name = "Label #4"
		l.Labels = []models.Label{r}
		l.Date = helpers.ParseDateNoError("2019-03-10")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		r := test.GetRandomLabel(33)
		r.Name = "Label #5"
		l.Labels = []models.Label{r}
		l.Date = helpers.ParseDateNoError("2019-03-15")
		db.LedgerCreate(&l)
	}

	// Set start / end
	start := helpers.ParseDateNoError("2019-03-01")
	end := helpers.ParseDateNoError("2019-06-30")

	// Run test function
	result := GetLabelsPnL(db, 33, start, end, "ASC")

	// Test results
	st.Expect(t, result[0].Name, "Label #1")
	st.Expect(t, result[0].Amount, 500.00)
	st.Expect(t, result[1].Name, "Label #2")
	st.Expect(t, result[1].Amount, 500.00)
	st.Expect(t, result[2].Name, "Label #3")
	st.Expect(t, result[2].Amount, -500.00)
	st.Expect(t, result[3].Name, "Label #4")
	st.Expect(t, result[3].Amount, -500.00)
	st.Expect(t, result[4].Name, "Label #5")
	st.Expect(t, result[4].Amount, 500.00)
}

//
// TestGetCategoriesPnL01 - Test Categories
//
func TestGetCategoriesPnL01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		l.Category.Name = "Category #1"
		l.Date = helpers.ParseDateNoError("2019-03-01")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		l.Category.Name = "Category #2"
		l.Date = helpers.ParseDateNoError("2019-03-05")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = -100
		l.Category.Name = "Category #3"
		l.Date = helpers.ParseDateNoError("2019-03-05")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = -100
		l.Category.Name = "Category #4"
		l.Date = helpers.ParseDateNoError("2019-03-10")
		db.LedgerCreate(&l)
	}

	// Create like 5 ledger entries
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		l.Amount = 100
		l.Category.Name = "Category #5"
		l.Date = helpers.ParseDateNoError("2019-03-15")
		db.LedgerCreate(&l)
	}

	// Set start / end
	start := helpers.ParseDateNoError("2019-03-01")
	end := helpers.ParseDateNoError("2019-06-30")

	// Run test function
	result := GetCategoriesPnL(db, 33, start, end, "ASC")

	// Test results
	st.Expect(t, result[0].Name, "Category #1")
	st.Expect(t, result[0].Amount, 500.00)
	st.Expect(t, result[1].Name, "Category #2")
	st.Expect(t, result[1].Amount, 500.00)
	st.Expect(t, result[2].Name, "Category #3")
	st.Expect(t, result[2].Amount, -500.00)
	st.Expect(t, result[3].Name, "Category #4")
	st.Expect(t, result[3].Amount, -500.00)
	st.Expect(t, result[4].Name, "Category #5")
	st.Expect(t, result[4].Amount, 500.00)

	// -------- Different Order -------------- //

	// Run test function
	result = GetCategoriesPnL(db, 33, start, end, "DESC")

	// Test results
	st.Expect(t, result[0].Name, "Category #5")
	st.Expect(t, result[0].Amount, 500.00)
	st.Expect(t, result[1].Name, "Category #4")
	st.Expect(t, result[1].Amount, -500.00)
	st.Expect(t, result[2].Name, "Category #3")
	st.Expect(t, result[2].Amount, -500.00)
	st.Expect(t, result[3].Name, "Category #2")
	st.Expect(t, result[3].Amount, 500.00)
	st.Expect(t, result[4].Name, "Category #1")
	st.Expect(t, result[4].Amount, 500.00)
}

//
// TestGetIncomeByContact01 - Get income by contact
//
func TestGetIncomeByContact01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
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

	// Create like 10 ledger entries for June
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-06-01")

		// test non-Name options
		if l.Contact.Name == "Home Depot" {
			l.Contact.Name = ""
		}

		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Set start / end
	start := helpers.ParseDateNoError("2019-03-01")
	end := helpers.ParseDateNoError("2019-06-30")

	// Run test function
	result := GetIncomeByContact(db, 33, start, end, "ASC")

	// Sort of a real test here.
	incomeTotal := 0.00

	for key := range dMap {
		if dMap[key].Amount > 0 {
			incomeTotal = incomeTotal + dMap[key].Amount
		}
	}

	// Build total from results
	total := 0.00

	for _, row := range result {
		total = total + row.Amount
	}

	// Test results.
	st.Expect(t, helpers.Round(total, 2), helpers.Round(incomeTotal, 2))
}

//
// TestGetExpenseByContact01 - Get expenses by vendor
//
func TestGetExpenseByContact01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
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

	// Create like 10 ledger entries for June
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-06-01")

		// test non-Name options
		if l.Contact.Name == "Home Depot" {
			l.Contact.Name = ""
		}

		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Set start / end
	start := helpers.ParseDateNoError("2019-03-01")
	end := helpers.ParseDateNoError("2019-06-30")

	// Run test function
	result := GetExpenseByContact(db, 33, start, end, "ASC")

	// Sort of a cheal test here.
	expenseTotal := 0.00

	for key := range dMap {
		if dMap[key].Amount < 0 {
			expenseTotal = expenseTotal + dMap[key].Amount
		}
	}

	// Build total from results
	total := 0.00

	for _, row := range result {
		total = total + row.Amount
	}

	// Test results.
	st.Expect(t, helpers.Round(total, 2), helpers.Round(expenseTotal, 2))
}

//
// TestGetPnL01 - return PnL by group / start / end
//
func TestGetPnL01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
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
		l.Date = helpers.ParseDateNoError("2019-03-01 00:00:00")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for April
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-04-01 00:00:00")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for May
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-05-01 00:00:00")
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries for June
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-06-01 00:00:00")
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

	// ---------- Security Check ------------- //

	// Run test function
	pl3 := GetPnL(db, 33, start, end, "month", "JJJJ")

	// Test results
	st.Expect(t, len(pl), 4)
	st.Expect(t, helpers.Round(pl3[0].Profit, 2), helpers.Round(profit032019, 2))
	st.Expect(t, helpers.Round(pl3[1].Profit, 2), helpers.Round(profit042019, 2))
	st.Expect(t, helpers.Round(pl3[2].Profit, 2), helpers.Round(profit052019, 2))
	st.Expect(t, helpers.Round(pl3[3].Profit, 2), helpers.Round(profit062019, 2))

	st.Expect(t, helpers.Round(pl3[0].Income, 2), helpers.Round(income032019, 2))
	st.Expect(t, helpers.Round(pl3[1].Income, 2), helpers.Round(income042019, 2))
	st.Expect(t, helpers.Round(pl3[2].Income, 2), helpers.Round(income052019, 2))
	st.Expect(t, helpers.Round(pl3[3].Income, 2), helpers.Round(income062019, 2))

	st.Expect(t, helpers.Round(pl3[0].Expense, 2), helpers.Round(expense032019, 2))
	st.Expect(t, helpers.Round(pl3[1].Expense, 2), helpers.Round(expense042019, 2))
	st.Expect(t, helpers.Round(pl3[2].Expense, 2), helpers.Round(expense052019, 2))
	st.Expect(t, helpers.Round(pl3[3].Expense, 2), helpers.Round(expense062019, 2))

	// ---------- Security Check - Group ------------- //

	// Run test function
	pl4 := GetPnL(db, 33, start, end, "blah", "ASC")

	// Test results
	st.Expect(t, len(pl4), 0)
}

//
// TestGetPnL02 - return PnL by group / start / end
//
func TestGetPnL02(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries - 2019-Q1
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-02-01 00:00:00")
		l.Amount = 100.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries - 2018-Q2
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2018-04-01 00:00:00")
		l.Amount = 200.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries - 2018-Q1
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2018-01-01 00:00:00")
		l.Amount = -100.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries - 2018-Q4
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2018-10-01 00:00:00")
		l.Amount = -200.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Set start / end
	start := helpers.ParseDateNoError("2018-01-01")
	end := helpers.ParseDateNoError("2019-12-30")

	// ---------- quarter ------------- //

	// Run test function
	pl := GetPnL(db, 33, start, end, "quarter", "ASC")

	// Test results
	st.Expect(t, len(pl), 4)

	st.Expect(t, pl[0].Profit, -1000.00)
	st.Expect(t, pl[1].Profit, 2000.00)
	st.Expect(t, pl[2].Profit, -2000.00)
	st.Expect(t, pl[3].Profit, 1000.00)

	st.Expect(t, pl[0].Income, 0.00)
	st.Expect(t, pl[1].Income, 2000.00)
	st.Expect(t, pl[2].Income, 0.00)
	st.Expect(t, pl[3].Income, 1000.00)

	st.Expect(t, pl[0].Expense, -1000.00)
	st.Expect(t, pl[1].Expense, 0.00)
	st.Expect(t, pl[2].Expense, -2000.00)
	st.Expect(t, pl[3].Expense, 0.00)
}

//
// TestGetPnL03 - Year
//
func TestGetPnL03(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create like 5 ledger entries. Diffent account.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries - 2019
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2019-02-01 00:00:00")
		l.Amount = 100.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries - 2018
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2018-04-01 00:00:00")
		l.Amount = 200.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries - 2017
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		l.Date = helpers.ParseDateNoError("2017-01-01 00:00:00")
		l.Amount = -100.00
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Set start / end
	start := helpers.ParseDateNoError("2017-01-01")
	end := helpers.ParseDateNoError("2019-12-30")

	// ---------- quarter ------------- //

	// Run test function
	pl := GetPnL(db, 33, start, end, "year", "ASC")

	// Test results
	st.Expect(t, len(pl), 3)

	st.Expect(t, pl[0].Date, "2017")
	st.Expect(t, pl[1].Date, "2018")
	st.Expect(t, pl[2].Date, "2019")

	st.Expect(t, pl[0].Profit, -1000.00)
	st.Expect(t, pl[1].Profit, 2000.00)
	st.Expect(t, pl[2].Profit, 1000.00)

	st.Expect(t, pl[0].Income, 0.00)
	st.Expect(t, pl[1].Income, 2000.00)
	st.Expect(t, pl[2].Income, 1000.00)

	st.Expect(t, pl[0].Expense, -1000.00)
	st.Expect(t, pl[1].Expense, 0.00)
	st.Expect(t, pl[2].Expense, 0.00)
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
	st.Expect(t, helpers.Round(pl.Value, 2), helpers.Round(total, 2))

	// ---------- Test empty year ------------- //

	// Run test function
	pl2 := GetCurrentYearPnL(db, 33, 2005)

	// Test results
	st.Expect(t, pl2.Year, 2005)
	st.Expect(t, pl2.Value, 0.00)
}
