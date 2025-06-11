//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"errors"
	"flag"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

// init sets default environment variables for stripe functionality during tests
func init() {
	// Only set defaults during tests
	if flag.Lookup("test.v") != nil {
		setDefaultIfEmpty("STRIPE_SECRET_KEY", "sk_test_default")
	}
}

// setDefaultIfEmpty sets an environment variable to a default value if it's not already set
func setDefaultIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}

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
