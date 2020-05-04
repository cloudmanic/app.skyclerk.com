//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//
// Notes: the reason this is in the services package is we don't
// want to conflict with the stripe name space
//

package stripe

import (
	"errors"
	"os"
	"strconv"

	stripe "github.com/stripe/stripe-go"

	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

//
// AddCustomer will add a new customer
//
func AddCustomer(first string, last string, email string, accountID int) (string, error) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddCustomer")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object
	customerParams := &stripe.CustomerParams{Email: &email}
	customerParams.AddMetadata("FirstName", first)
	customerParams.AddMetadata("LastName", last)
	customerParams.AddMetadata("AccountId", strconv.Itoa(accountID))

	// Create new customer.
	customer, err := customer.New(customerParams)

	if err != nil {
		return "", err
	}

	// Return the new customer Id
	return customer.ID, nil
}

//
// AddSubscription - Add a customer subscription.
//
func AddSubscription(custId string, plan string, coupon string) (string, error) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		//Critical(errors.New("No STRIPE_SECRET_KEY found in StripeAddSubscription"))
		return "", errors.New("No STRIPE_SECRET_KEY found in StripeAddSubscription")
	}

	// Setup Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(custId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(plan),
			},
		},
	}

	// Do we have a coupon code?
	if len(coupon) > 0 {
		subParams.Coupon = stripe.String(coupon)
		//Info(errors.New("Coupon code passed with subscribe token: " + coupon + " - " + custId))
	}

	// Create new subscription.
	subscription, err := sub.New(subParams)

	if err != nil {
		//Info(errors.New("StripeAddSubscription : Unable to create a new subscription. (" + err.Error() + ")"))
		return "", err
	}

	// Return the new subscription Id
	return subscription.ID, nil
}

//
// DeleteCustomer - Will delete a  customer.
//
func DeleteCustomer(custToken string) error {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCustomer")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object and Delete
	params := &stripe.CustomerParams{}
	_, err := customer.Del(custToken, params)

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

//
// GetCustomer get customer from totken.
//
func GetCustomer(custToken string) (*stripe.Customer, error) {

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCustomer")
	}

	// Add Stripe Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Setup the customer object and Delete
	cust, err := customer.Get(custToken, nil)

	if err != nil {
		return nil, err
	}

	// Return happy
	return cust, nil
}

/* End File */
