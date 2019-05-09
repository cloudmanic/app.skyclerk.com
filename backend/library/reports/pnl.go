//
// Date: 5/9/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"app.skyclerk.com/backend/models"
)

// YearPnL struct
type YearPnL struct {
	Year  int     `json:"year"`
	Value float64 `json:"value"`
}

//
// GetCurrentYearPnL - return the current year and the profit and lost of that year.
//
func GetCurrentYearPnL(db models.Datastore, accountId uint, year int) YearPnL {
	// Struct we return
	rt := YearPnL{}

	// Build sql
	sql := "SELECT SUM(LedgerAmount) AS value, Year(LedgerDate) AS year FROM Ledger WHERE LedgerAccountId = ? AND Year(LedgerDate) = ?"

	// Run query
	db.New().Raw(sql, accountId, year).Scan(&rt)

	// Return happy.
	return rt
}

/* End File */
