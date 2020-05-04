//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"errors"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

//
// GetChargesByCustomer will get charges by customer.
//
func GetChargesByCustomer(ID string) ([]*stripe.Charge, error) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetChargesByCustomer")
	}

	charges := []*stripe.Charge{}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get the charges
	params := &stripe.ChargeListParams{}
	params.Filters.AddFilter("limit", "", "100")
	params.Filters.AddFilter("customer", "", ID)

	i := charge.List(params)

	for i.Next() {
		c := i.Charge()
		charges = append(charges, c)
	}

	// Return happy
	return charges, nil
}

/* End File */
