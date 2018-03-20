//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

// Database interface
type Datastore interface {

	// Ledger
	CreateLedger(ledger *Ledger) error
}

/* End File */
