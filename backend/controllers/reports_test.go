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
// TestReportsPnlLabel01	 - Get pnl by label
//
func TestReportsPnlLabel01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup controllers
	c := &Controller{}
	c.SetDB(db)

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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/pnl-label?start=2019-03-01&end=2019-06-30&sort=asc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/pnl-label", c.ReportsPnlLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []reports.NameValue{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
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
// TestReportsPnlCategory01	 - Get pnl by category
//
func TestReportsPnlCategory01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup controllers
	c := &Controller{}
	c.SetDB(db)

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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/pnl-category?start=2019-03-01&end=2019-06-30&sort=asc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/pnl-category", c.ReportsPnlCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []reports.NameValue{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
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
}

//
// TestReportsIncomeByContact01	 - Get income by contact
//
func TestReportsIncomeByContact01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/income-by-contact?start=2019-03-01&end=2019-06-30&sort=asc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/income-by-contact", c.ReportsIncomeByContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []reports.NameValue{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Sort of a cheal test here.
	incomeTotal := 0.00

	for key := range dMap {
		if dMap[key].Amount > 0 {
			incomeTotal = incomeTotal + dMap[key].Amount
		}
	}

	// Build total from results
	total := 0.00

	for _, row := range results {
		total = total + row.Amount
	}

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, helpers.Round(total, 2), helpers.Round(incomeTotal, 2))
}

//
// TestReportsExpensesByContact01	 - Get expenses by contact
//
func TestReportsExpensesByContact01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/expenses-by-contact?start=2019-03-01&end=2019-06-30&sort=asc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/expenses-by-contact", c.ReportsExpensesByContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []reports.NameValue{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Sort of a cheal test here.
	expenseTotal := 0.00

	for key := range dMap {
		if dMap[key].Amount < 0 {
			expenseTotal = expenseTotal + dMap[key].Amount
		}
	}

	// Build total from results
	total := 0.00

	for _, row := range results {
		total = total + row.Amount
	}

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, helpers.Round(total, 2), helpers.Round(expenseTotal, 2))
}

//
// TestReportsPnl01	 - test getting the current year P&L
//
func TestReportsPnl01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/reports/pnl?start=2019-03-01&end=2019-06-30&sort=asc&group=month", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/reports/pnl", c.ReportsPnl)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	pl := []reports.PnL{}
	err := json.Unmarshal([]byte(w.Body.String()), &pl)

	// Test results
	st.Expect(t, err, nil)
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
}

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
