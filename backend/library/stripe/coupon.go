//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

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
