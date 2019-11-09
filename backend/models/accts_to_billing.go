//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

// AcctToBilling struct NOTE: We do not run this in our migrations. It is auto creeated.
type AcctToBilling struct {
	AccountId uint `sql:"not null"  json:"-"`
	BillingId uint `sql:"not null"  json:"-"`
}

/* End File */
