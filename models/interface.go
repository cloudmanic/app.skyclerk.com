//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

// Database interface
type Datastore interface {

	// Ledger
	LedgerCreate(ledger *Ledger) error
}

/* End File */
