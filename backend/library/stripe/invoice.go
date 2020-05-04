//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"errors"
	"os"

	"github.com/stripe/stripe-go/invoice"

	stripe "github.com/stripe/stripe-go"
)

//
// GetInvoice will get an invoice by Id
//
func GetInvoice(ID string) (*stripe.Invoice, error) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetInvoice")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get the invoice
	iv, err := invoice.Get(ID, nil)

	if err != nil {
		return nil, err
	}

	// Return happy
	return iv, nil
}

/* End File */
