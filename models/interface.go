//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

// Database interface
type Datastore interface {
	// Generic database functions
	Count(model interface{}, params QueryParam) (uint, error)
	Query(model interface{}, params QueryParam) error
	QueryMeta(model interface{}, params QueryParam) (QueryMetaData, error)
	QueryWithNoFilterCount(model interface{}, params QueryParam) (int, error)
	GetQueryMetaData(limitCount int, noLimitCount int, params QueryParam) QueryMetaData

	// Ledger
	LedgerCreate(ledger *Ledger) error
}

/* End File */
