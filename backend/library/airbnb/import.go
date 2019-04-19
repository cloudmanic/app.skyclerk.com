//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package airbnb

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/cloudmanic/app.skyclerk.com/backend/models"
	"github.com/cloudmanic/app.skyclerk.com/backend/services"
)

//
// Import CSV Ledger Entries
//
func CSVImport(db *models.DB, accountId uint, file string) int {

	var count int = 0

	// Open CSV file and read it.
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// Lop through each line in this CSV
	for {

		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			services.LogFatal(err)
		}

		// We are only interested in payouts.
		if line[1] != "Reservation" {
			continue
		}

		// Check to see if we already have this confirmation code in our system.
		l := models.Ledger{}
		db.Where("LedgerAirBnbHash = ?", line[2]).First(&l)

		if l.Id > 0 {
			continue
		}

		// Get the amount.
		amount, err := strconv.ParseFloat(line[10], 64)

		if err != nil {
			services.Error(err)
		}

		// Get the date
		date, err := dateparse.ParseAny(line[0])

		if err != nil {
			services.Error(err)
		}

		// Build the contact
		contact := models.Contact{
			AccountId: accountId,
			Name:      line[5],
		}

		// Setup the label
		labels := []models.Label{
			{AccountId: accountId, Name: "airbnb"},
		}

		// Setup the category
		category := models.Category{
			Type:      "2",
			Irs:       "1",
			Show:      "1",
			AccountId: accountId,
			Name:      "Rental Income",
		}

		// Setup the Ledger entry.
		ledger := models.Ledger{
			AccountId:  accountId,
			Date:       date,
			Amount:     amount,
			Contact:    contact,
			Category:   category,
			Labels:     labels,
			Note:       line[6] + " : " + line[2],
			AirBnbHash: line[2],
		}

		// Store the entry in the database
		err = db.LedgerCreate(&ledger)

		if err != nil {
			services.Error(err)
		} else {
			count++
		}
	}

	// Return the total imported
	return count
}

/* End File */
