//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package airbnb

import (
	"testing"

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/nbio/st"
)

//
// Test CSVImport
//
func TestCSVImport01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Run Import
	importCount := CSVImport(db, 33, "./airbnb_.csv")

	// Get the ledger entries we imported.
	ledgers := []models.Ledger{}
	db.Preload("Category").Preload("Contact").Preload("Labels").Find(&ledgers)

	// Test results
	st.Expect(t, importCount, 62)
	st.Expect(t, len(ledgers), 62)
	st.Expect(t, ledgers[61].Id, uint(62))
	st.Expect(t, ledgers[61].AccountId, uint(33))
	st.Expect(t, ledgers[61].Contact.Name, "Maya Richman")
	st.Expect(t, ledgers[61].Contact.AccountId, uint(33))
	st.Expect(t, ledgers[61].Category.Name, "Rental Income")
	st.Expect(t, ledgers[61].Category.Type, "2")
	st.Expect(t, ledgers[61].Labels[0].Name, "airbnb")
	st.Expect(t, ledgers[61].Labels[0].AccountId, uint(33))
	st.Expect(t, ledgers[61].Amount, 306.00)
	st.Expect(t, ledgers[61].Note, "Mt. Bachelor Village Apartment : P32WRZ")
}

/* End File */
