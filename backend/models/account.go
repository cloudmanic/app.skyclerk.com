//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	"app.skyclerk.com/backend/library/sendy"
	"app.skyclerk.com/backend/library/stripe"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Account struct
type Account struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `sql:"not null" json:"-"`
	UpdatedAt    time.Time `sql:"not null" json:"-"`
	OwnerId      uint      `sql:"not null" json:"owner_id"`
	BillingId    uint      `sql:"not null" json:"-"`
	Name         string    `sql:"not null" json:"name"`
	Address      string    `sql:"not null;type:TEXT" json:"-"`
	City         string    `sql:"not null" json:"-"`
	State        string    `sql:"not null" json:"-"`
	Zip          string    `sql:"not null" json:"-"`
	Country      string    `sql:"not null" json:"-"`
	Locale       string    `sql:"not null;default:'en-US'" json:"locale"` // BCP 47 language tag
	Currency     string    `sql:"not null;default:'USD'" json:"currency"` // The ISO 4217 currency code, such as USD for the US dollar and EUR for the euro.
	LastActivity time.Time `sql:"not null" json:"-"`
}

//
// Validate for this model.
//
func (a Account) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {
	return validation.ValidateStruct(&a,
		// Name
		validation.Field(&a.Name,
			validation.Required.Error("The name field is required."),
			validation.Length(1, 50).Error("The name is too long."),
		),

		// Locale
		validation.Field(&a.Locale,
			validation.Required.Error("The locale field is required."),
		),

		// Currency
		validation.Field(&a.Currency,
			validation.Required.Error("The currency field is required."),
		),

		// OwnerId
		validation.Field(&a.OwnerId,
			validation.Required.Error("The owner_id field is required."),
			validation.By(func(value interface{}) error { return db.ValidateOwnerId(a, accountId, objId, action) }),
		),
	)
}

//
// ValidateOwnerId - is this a valid owner id?
//
func (db *DB) ValidateOwnerId(acct Account, accountId uint, objId uint, action string) error {
	// Default error message
	const errMsg = "Invalid owner_id was posted."

	// Make sure this user is part of the account.
	c := AcctToUsers{}
	if db.New().Where("account_id = ? AND user_id = ?", accountId, acct.OwnerId).First(&c).RecordNotFound() {
		return errors.New(errMsg)
	}

	// All good in the hood
	return nil
}

//
// GetAccountById - Get a account by Id.
//
func (t *DB) GetAccountById(id uint) (Account, error) {
	var u Account

	if t.Where("id = ?", id).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil
}

//
// ClearAccount
//
func (t *DB) ClearAccount(accountId uint) {
	// Ledger Ids.
	var ledgerIds []uint

	// Clear labels to ledger.
	l := []Ledger{}
	t.New().Select("LedgerId").Where("LedgerAccountId = ?", accountId).Find(&l)

	// Loop through and build ledger id list
	for _, row := range l {
		ledgerIds = append(ledgerIds, row.Id)
	}

	// Delete look up tables.
	t.New().Where("FilesToLedgerLedgerId IN (?)", ledgerIds).Delete(FilesToLedger{})
	t.New().Where("LabelsToLedgerLedgerId IN (?)", ledgerIds).Delete(LabelsToLedger{})

	// Clear database tables.
	t.New().Exec("DELETE FROM activities WHERE account_id = ?", accountId)
	t.New().Exec("DELETE FROM invites WHERE account_id = ?", accountId)
	t.New().Exec("DELETE FROM Labels WHERE LabelsAccountId = ?", accountId)
	t.New().Exec("DELETE FROM Ledger WHERE LedgerAccountId = ?", accountId)
	t.New().Exec("DELETE FROM Files WHERE FilesAccountId = ?", accountId)
	t.New().Exec("DELETE FROM Contacts WHERE ContactsAccountId = ?", accountId)
	t.New().Exec("DELETE FROM Categories WHERE CategoriesAccountId = ?", accountId)
	t.New().Exec("DELETE FROM SnapClerk WHERE SnapClerkAccountId = ?", accountId)

	// TODO(spicer): delete files at AWS too.
}

//
// DeleteAccount will delete an account
//
func (t *DB) DeleteAccount(accountID uint) {
	// Get the account
	account := Account{}
	t.New().Find(&account, accountID)

	if account.Id != accountID {
		return
	}

	// Get the owner
	owner := User{}
	t.New().Find(&owner, account.OwnerId)

	// Get the billing profile.
	billing := Billing{}
	t.New().Find(&billing, account.BillingId)

	// See if we have any other acccounts with this billing id
	ba := []Account{}
	t.New().Where("billing_id = ?", billing.Id).Find(&ba)

	if len(ba) <= 1 {
		t.New().Exec("DELETE FROM billings WHERE id = ?", billing.Id)

		// Remove account at Stripe.
		if len(billing.StripeCustomer) > 0 {
			stripe.DeleteCustomer(billing.StripeCustomer)
		}
	}

	// Clear users not used else where.
	a2u := []AcctToUsers{}
	t.New().Where("account_id = ?", accountID).Find(&a2u)

	// Loop through users delete if needed
	for _, row := range a2u {
		u := []AcctToUsers{}
		t.New().Where("user_id = ?", row.UserId).Find(&u)

		if len(u) == 1 {
			t.New().Exec("DELETE FROM users WHERE id = ?", row.UserId)
			t.New().Exec("DELETE FROM sessions WHERE user_id = ?", row.UserId)
		}
	}

	// Clear database tables.
	t.New().Exec("DELETE FROM acct_to_users WHERE account_id = ?", accountID)
	t.New().Exec("DELETE FROM accounts WHERE id = ?", accountID)

	// Clear at sendy users
	go sendy.Unsubscribe("trial", owner.Email)
	go sendy.Unsubscribe("expired", owner.Email)
}

/* End File */
