//
// Date: 5/2/2020
// Author(s): Spicer Matthews (spicer@skyclerk.// COMBAK: )
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

// //
// // GetBalanceTransaction will return one transaction balance.
// //
// func GetBalanceTransaction(ID string) (*stripe.BalanceTransaction, error) {
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction")
// 	}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Get the transaction
// 	bt, err := balance.GetBalanceTransaction(ID, nil)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Return happy
// 	return bt, nil
// }
//
// //
// // GetBalanceTransactions will return all transaction balance.
// //
// func GetBalanceTransactions() ([]*stripe.BalanceTransaction, error) {
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in StripeGetBalanceTransaction")
// 	}
//
// 	list := []*stripe.BalanceTransaction{}
//
// 	// Add Stripe Key
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	// Get the charges
// 	params := &stripe.BalanceTransactionListParams{}
// 	params.Filters.AddFilter("limit", "", "100")
//
// 	i := balancetransaction.List(params)
//
// 	for i.Next() {
// 		c := i.Charge()
// 		list = append(charges, c)
// 	}
//
// 	// Return happy
// 	return list, nil
// }

/* End File */
