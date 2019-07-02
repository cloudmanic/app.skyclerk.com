//
// Date: 5/9/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"strings"
	"time"

	"app.skyclerk.com/backend/models"
)

// NameValue struct
type NameValue struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// PnL struct
type PnL struct {
	Date    string  `json:"date"`
	Profit  float64 `json:"profit"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// YearPnL struct
type YearPnL struct {
	Year  int     `json:"year"`
	Value float64 `json:"value"`
}

//
// GetCategoriesPnL returns catories and the total for the time period.
//
func GetCategoriesPnL(db models.Datastore, accountId uint, start time.Time, end time.Time, sort string) []NameValue {
	// SQL String
	sql := "SELECT CategoriesName as name, SUM(LedgerAmount) as amount "
	sql = sql + "FROM Ledger JOIN Categories ON Categories.CategoriesId = Ledger.LedgerCategoryId "
	sql = sql + "WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ? "
	sql = sql + "GROUP BY CategoriesName ORDER BY name "

	// Struct we return
	rt := []NameValue{}

	// Quick security check
	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}

	// Add in sort
	sql = sql + sort

	// Run query
	db.New().Raw(sql, accountId, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&rt)

	// Return happy.
	return rt
}

//
// GetIncomeByContact - Get income by contact
//
func GetIncomeByContact(db models.Datastore, accountId uint, start time.Time, end time.Time, sort string) []NameValue {
	// SQL String
	sql := "SELECT IF(LENGTH(ContactsName)>0,ContactsName, CONCAT(ContactsFirstName, \" \", ContactsLastName)) AS name, "
	sql = sql + "sum(LedgerAmount) AS amount "
	sql = sql + "FROM Ledger "
	sql = sql + "JOIN Contacts ON Contacts.ContactsId = Ledger.LedgerContactId "
	sql = sql + "WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ? "
	sql = sql + "AND LedgerAmount > 0 GROUP BY name ORDER BY name "

	// Struct we return
	rt := []NameValue{}

	// Quick security check
	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}

	// Add in sort
	sql = sql + sort

	// Run query
	db.New().Raw(sql, accountId, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&rt)

	// Return happy.
	return rt
}

//
// GetExpenseByContact - Get expense by contact
//
func GetExpenseByContact(db models.Datastore, accountId uint, start time.Time, end time.Time, sort string) []NameValue {
	// SQL String
	sql := "SELECT IF(LENGTH(ContactsName)>0,ContactsName, CONCAT(ContactsFirstName, \" \", ContactsLastName)) AS name, "
	sql = sql + "sum(LedgerAmount) AS amount "
	sql = sql + "FROM Ledger "
	sql = sql + "JOIN Contacts ON Contacts.ContactsId = Ledger.LedgerContactId "
	sql = sql + "WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ? "
	sql = sql + "AND LedgerAmount < 0 GROUP BY name ORDER BY name "

	// Struct we return
	rt := []NameValue{}

	// Quick security check
	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}

	// Add in sort
	sql = sql + sort

	// Run query
	db.New().Raw(sql, accountId, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&rt)

	// Return happy.
	return rt
}

//
// GetPnL - Profit / Loss based on start / stop  group
//
func GetPnL(db models.Datastore, accountId uint, start time.Time, end time.Time, group string, sort string) []PnL {
	// SQL String
	sql := ""

	// Struct we return
	rt := []PnL{}

	// Quick security check
	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}

	// Build sql
	switch group {
	case "month":
		sql = "SELECT date_format(LedgerDate, '%Y-%m') AS date, SUM(LedgerAmount) AS profit, SUM(CASE WHEN LedgerAmount>0 THEN LedgerAmount ELSE 0 END) AS income, SUM(CASE WHEN LedgerAmount<0 THEN LedgerAmount ELSE 0 END) AS expense FROM Ledger WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ? GROUP BY date_format(LedgerDate, '%Y-%m') ORDER BY date " + sort

	case "quarter":
		sql = `SELECT CONCAT(YEAR(LedgerDate), '-Q',quarter(LedgerDate)) AS date,
		SUM(LedgerAmount) AS profit,
		SUM(CASE WHEN LedgerAmount>0 THEN LedgerAmount ELSE 0 END) AS income,
		SUM(CASE WHEN LedgerAmount<0 THEN LedgerAmount ELSE 0 END) AS expense

		FROM Ledger

		WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ?

		GROUP BY date ORDER BY date ` + sort

	case "year":
		sql = `SELECT
		date_format(LedgerDate, '%Y') AS date,
		SUM(LedgerAmount) AS profit,
		SUM(CASE WHEN LedgerAmount>0 THEN LedgerAmount ELSE 0 END) AS income,
		SUM(CASE WHEN LedgerAmount<0 THEN LedgerAmount ELSE 0 END) AS expense
		FROM Ledger
		WHERE LedgerAccountId = ? AND LedgerDate >= ? AND LedgerDate <= ?
		GROUP BY date ORDER BY date ` + sort

	default:
		return rt
	}

	// Run query
	db.New().Raw(sql, accountId, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(&rt)

	// Return happy.
	return rt
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

	// If we have no values we just add in this year.
	if rt.Year == 0 {
		rt.Year = year
	}

	// Return happy.
	return rt
}

/* End File */
