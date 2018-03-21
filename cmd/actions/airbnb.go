//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"fmt"

	"github.com/cloudmanic/skyclerk.com/library/airbnb"
	"github.com/cloudmanic/skyclerk.com/models"
)

//
// Take a CSV files and import it from
//
// go run main.go -cmd=airbnb-import -file=/Users/spicer/Downloads/airbnb_.csv -account_id=4992
//
func AirBnbImport(db *models.DB, accountId uint, file string) {

	importCount := airbnb.CSVImport(db, accountId, file)
	fmt.Println(importCount, "New Ledger Entries Successfully Imported")

}

/* End File */
