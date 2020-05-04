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

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
)

//
// AddCreditCardByToken - Add a new credit card.
//
func AddCreditCardByToken(custID string, token string) (string, error) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		return "", errors.New("No STRIPE_SECRET_KEY found in AddCreditCardByToken")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CardParams{
		Customer: stripe.String(custID),
		Token:    stripe.String(token),
	}

	// call to stripe to add card
	c, err := card.New(params)

	if err != nil {
		return "", err
	}

	// Return the new card Id
	return c.ID, nil
}

//
// //
// // List all credit cards on file.
// //
// func StripeListAllCreditCards(custId string) ([]string, error) {
//
// 	cardList := []string{}
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in ListAllCreditCards"))
// 		return nil, errors.New("No STRIPE_SECRET_KEY found in ListAllCreditCards")
// 	}
//
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	params := &stripe.CardListParams{
// 		Customer: stripe.String(custId),
// 	}
//
// 	params.Filters.AddFilter("limit", "", "100")
//
// 	list := card.List(params)
//
// 	for list.Next() {
// 		c := list.Card()
// 		cardList = append(cardList, c.ID)
// 	}
//
// 	// Return the card list
// 	return cardList, nil
//
// }
//
// //
// // Delete credit cards on file.
// //
// func StripeDeleteCreditCard(custId string, cardId string) error {
//
// 	// Make sure we have a STRIPE_SECRET_KEY
// 	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
// 		Critical(errors.New("No STRIPE_SECRET_KEY found in DeleteCreditCard"))
// 		return errors.New("No STRIPE_SECRET_KEY found in DeleteCreditCard")
// 	}
//
// 	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
//
// 	params := &stripe.CardParams{
// 		Customer: stripe.String(custId),
// 	}
//
// 	// Delete card at stripe
// 	_, err := card.Del(cardId, params)
//
// 	if err != nil {
// 		Info(errors.New("DeleteCreditCard : Unable to delete card. (" + err.Error() + ")"))
// 		return err
// 	}
//
// 	// Return happy
// 	return nil
//
// }
