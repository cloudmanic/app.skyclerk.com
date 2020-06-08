//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"fmt"

	"app.skyclerk.com/backend/models"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balancetransaction"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

//
// Sync will bring in all transactions for a stripe account.
//
func Sync(connectedAccount models.ConnectedAccounts) {

	stripe.Key = connectedAccount.StripeAccessToken

	params := &stripe.ChargeListParams{
		CreatedRange: &stripe.RangeQueryParams{
			GreaterThan: connectedAccount.StripeLastItem,
		},
	}

	params.Filters.AddFilter("limit", "", "100")

	i := charge.List(params)
	for i.Next() {
		c := i.Charge()

		// Only collect succeeded charges
		if c.Status != "succeeded" {
			continue
		}

		bt, _ := balancetransaction.Get(c.BalanceTransaction.ID, nil)

		cust, _ := customer.Get(c.Customer.ID, nil)

		fmt.Println(c.ID, c.Amount, c.BalanceTransaction.ID, bt.Fee, c.Created, c.Customer.ID, cust.Email, cust.Name, cust.Description)
	}

}

/* End File */
