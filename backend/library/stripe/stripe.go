//
// Date: 9/2/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//
// Notes: the reason this is in the services package is we don't
// want to conflict with the stripe name space
//

package stripe

// //
// // Get an invoice by Id
// //
// func StripeGetInvoice(id string) (*stripe.Invoice, error) {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeGetInvoice"))
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetInvoice")
// 	}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Get the invoice
// 	iv, err := invoice.Get(id, nil)
//
// 	if err != nil {
// 		Info(errors.New("StripeGetInvoice : Unable to get an invoice. " + id + " (" + err.Error() + ")"))
// 		return nil, err
// 	}
//
// 	// Return happy
// 	return iv, nil
//
// }
//
// //
// // Get charges by customer.
// //
// func StripeGetChargesByCustomer(id string) ([]*stripe.Charge, error) {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction"))
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetChargesByCustomer")
// 	}
//
// 	charges := []*stripe.Charge{}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Get the charges
// 	params := &stripe.ChargeListParams{}
// 	params.Filters.AddFilter("limit", "", "100")
// 	params.Filters.AddFilter("customer", "", id)
//
// 	i := charge.List(params)
//
// 	for i.Next() {
// 		c := i.Charge()
// 		charges = append(charges, c)
// 	}
//
// 	// Return happy
// 	return charges, nil
// }
//
// //
// // Get one transaction balance.
// //
// func StripeGetBalanceTransaction(id string) (*stripe.BalanceTransaction, error) {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction"))
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction")
// 	}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Get the transaction
// 	bt, err := balance.GetBalanceTransaction(id, nil)
//
// 	if err != nil {
// 		Info(errors.New("StripeGetBalanceTransaction : Unable to get a transaction balance. " + id + " (" + err.Error() + ")"))
// 		return nil, err
// 	}
//
// 	// Return happy
// 	return bt, nil
// }
//
// //
// // Apply a coupon to a Subscription
// //
// func StripeApplyCoupon(subId string, couponId string) error {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCoupon"))
// 		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteNewCoupon")
// 	}
//
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	params := &stripe.SubscriptionParams{
// 		Coupon: stripe.String(couponId),
// 	}
//
// 	// Send request to stripe
// 	_, err := sub.Update(subId, params)
//
// 	if err != nil {
// 		Info(errors.New("StripeApplyCoupon : Unable to apply a coupon. " + couponId + " (" + err.Error() + ")"))
// 		return err
// 	}
//
// 	// Return happy
// 	return nil
//
// }
//
// //
// // Get Coupon.
// //
// func StripeGetCoupon(id string) (*stripe.Coupon, error) {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeAddCustomer"))
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetCoupon")
// 	}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Setup the customer object and Delete
// 	c, err := coupon.Get(id, nil)
//
// 	if err != nil {
// 		Info(errors.New("StripeGetCoupon : Unable to get a coupon. " + id + " (" + err.Error() + ")"))
// 		return nil, err
// 	}
//
// 	// Return happy
// 	return c, nil
//
// }
//
// //
// // Delete a coupon.
// //
// func StripeDeleteCoupon(id string) error {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in StripeDeleteCoupon"))
// 		return errors.New("No STRIPE_SECRET_KEY found in StripeDeleteNewCoupon")
// 	}
//
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Delete coupon.
// 	_, err := coupon.Del(id, nil)
//
// 	if err != nil {
// 		Info(errors.New("StripeDeleteCoupon : Unable to create a new coupon. (" + err.Error() + ")"))
// 		return err
// 	}
//
// 	// Return the new coupon Id
// 	return nil
// }
//
// //
// // Create a new coupon.
// //
// func StripeCreateNewCoupon(name string, percentOff float64) (string, error) {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in CreateNewCoupon"))
// 		return "", errors.New("No STRIPE_SECRET_KEY found in StripeCreateNewCoupon")
// 	}
//
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Setup the coupon object
// 	params := &stripe.CouponParams{
// 		PercentOff: stripe.Float64(percentOff),
// 		Duration:   stripe.String(string(stripe.CouponDurationForever)),
// 		Name:       stripe.String(name),
// 	}
//
// 	// Create new coupon.
// 	couponObj, err := coupon.New(params)
//
// 	if err != nil {
// 		Info(errors.New("StripeCreateNewCoupon : Unable to create a new coupon. (" + err.Error() + ")"))
// 		return "", err
// 	}
//
// 	// Return the new coupon Id
// 	return couponObj.ID, nil
// }

/* End File */
