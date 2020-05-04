//
// Date: 2019-03-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package test

import (
	"fmt"
	"math/rand"
	"time"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/models"
)

//
// GetRandomAccount returns a random account.
//
func GetRandomAccount(accountId int64) models.Account {
	rand.Seed(time.Now().UnixNano())

	dates := []time.Time{
		time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC),
		time.Date(2018, 8, 19, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 1, 10, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 10, 29, 17, 20, 01, 507451, time.UTC),
	}

	zip := []string{"13601", "97222", "12384", "97345", "97701"}
	state := []string{"NY", "OR", "TX", "NC", "SC", "VT", "CA"}
	city := []string{"Newberg", "Watertown", "New York", "Clayton", "Portland", "Road Town", "Seattle"}

	acc := models.Account{
		Id:           uint(accountId),
		OwnerId:      uint(1),
		Name:         "Name " + helpers.RandStr(16),
		Address:      "Address " + helpers.RandStr(16),
		City:         city[rand.Intn(len(city))],
		State:        state[rand.Intn(len(state))],
		Zip:          zip[rand.Intn(len(zip))],
		Country:      "USA",
		Currency:     "USD",
		Locale:       "en-US",
		LastActivity: dates[rand.Intn(len(dates))],
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	return acc
}

//
// GetRandomBilling returns a random billing entry.
//
func GetRandomBilling(billingID int64, accountID int64) models.Billing {
	// Get billing model
	b := models.Billing{
		Id:                 uint(billingID),
		UpdatedAt:          time.Now(),
		CreatedAt:          time.Now(),
		StripeCustomer:     "cus_HDNy0kYb6Q8Zh8",
		StripeSubscription: "",
		Status:             "Active",
		TrialExpire:        time.Now().AddDate(0, 1, 0),
	}

	return b
}

//
// GetRandomApplication returns a random application.
//
func GetRandomApplication() models.Application {
	app := models.Application{
		Name:      "Unit Test Application - " + helpers.RandStr(16),
		ClientId:  helpers.RandStr(16),
		Secret:    helpers.RandStr(16),
		GrantType: "password",
	}

	return app
}

//
// GetRandomUser returns a random user.
//
func GetRandomUser(accountId int64) models.User {
	rand.Seed(time.Now().UnixNano())

	first := []string{"Bob", "Sue", "Jack", "AH", "Joe", "Spicer", "Steve"}
	last := []string{"Smith", "Doe", "Johnson", "Kaufmann", "Matthews", "Jobs"}

	salt := "salt123"
	pass := "F00bAr123"
	passMd5 := helpers.GetMd5(pass + salt)

	user := models.User{
		FirstName:   first[rand.Intn(len(first))],
		LastName:    last[rand.Intn(len(last))],
		Email:       helpers.RandStr(16) + "@example.com",
		Md5Password: passMd5,
		Md5Salt:     salt,
		Status:      "Active",
	}

	return user
}

//
// GetRandomLedger returns a random ledger.
//
func GetRandomLedger(accountId int64) models.Ledger {
	rand.Seed(time.Now().UnixNano())

	dates := []time.Time{
		time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC),
		time.Date(2018, 8, 19, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 1, 10, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 10, 29, 17, 20, 01, 507451, time.UTC),
	}
	amounts := []float64{1234.56, 33.44, 99.00, 555.32, 4583.01, 3.01, 0.20, 3429.34, 823.19, -44.34, -1234.53, -10.66, -3453.12, -192.33}

	ledger := models.Ledger{
		AccountId: uint(accountId),
		Date:      dates[rand.Intn(len(dates))],
		Amount:    amounts[rand.Intn(len(amounts))],
		Note:      "Test Note - " + helpers.RandStr(16),
		Contact:   GetRandomContact(accountId),
		Category:  GetRandomCategory(accountId),
		Labels:    []models.Label{GetRandomLabel(accountId), GetRandomLabel(accountId), GetRandomLabel(accountId)},
	}

	return ledger
}

//
// GetRandomContact returns a random contact.
//
func GetRandomContact(accountId int64) models.Contact {
	rand.Seed(time.Now().UnixNano())

	first := []string{"Bob", "Sue", "Jack", "AH", "Joe", "Spicer", "Steve"}
	last := []string{"Smith", "Doe", "Johnson", "Kaufmann", "Matthews", "Jobs"}
	company := []string{"Cloudmanic Labs, LLC", "Options Cafe", "Skyckerk", "Apple Inc.", "Home Depot"}

	contact := models.Contact{
		AccountId: uint(accountId),
		Name:      company[rand.Intn(len(company))],
		FirstName: first[rand.Intn(len(first))],
		LastName:  last[rand.Intn(len(last))],
		Email:     helpers.RandStr(16) + "@example.com",
		Avatar:    fmt.Sprintf("accounts/%d/avatars/5.png", accountId),
	}

	return contact
}

//
// GetRandomCategory returns a random category.
//
func GetRandomCategory(accountId int64) models.Category {
	rand.Seed(time.Now().UnixNano())

	catType := []string{"1", "2"}
	name := []string{"Sales", "Mailing", "Marketing", "Taxes", "Computers"}

	category := models.Category{
		AccountId: uint(accountId),
		Name:      name[rand.Intn(len(name))],
		Type:      catType[rand.Intn(len(catType))],
	}

	return category
}

//
// GetRandomLabel returns a random label.
//
func GetRandomLabel(accountId int64) models.Label {
	rand.Seed(time.Now().UnixNano())

	name := []string{"Options Cafe", "Skyclerk", "Clients", "Refund", "Non-Taxable", "Marketing", "1099", "Deposit", "125 Main Street"}

	label := models.Label{
		AccountId: uint(accountId),
		Name:      name[rand.Intn(len(name))],
	}

	return label
}

//
// GetRandomSnapClerk returns a random SnapClerk.
//
func GetRandomSnapClerk(accountId int64) models.SnapClerk {
	rand.Seed(time.Now().UnixNano())

	amounts := []float64{1234.56, 33.44, 99.00, 555.32, 4583.01, 3.01, 0.20, 3429.34, 823.19}

	sc := models.SnapClerk{
		AccountId: uint(accountId),
		Status:    "Pending",
		AddedById: 1,
		FileId:    0,
		File:      GetRandomFile(accountId),
		Amount:    amounts[rand.Intn(len(amounts))],
		Contact:   helpers.RandStr(10),
		Category:  helpers.RandStr(5),
		Labels:    fmt.Sprintf("%s,%s,%s", helpers.RandStr(5), helpers.RandStr(5), helpers.RandStr(5)),
		Note:      "Test Note - " + helpers.RandStr(16),
		Lat:       helpers.RandStr(6), // TODO(spicer): Make these real random lat / lons
		Lon:       helpers.RandStr(6),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Return Happy
	return sc
}

//
// GetRandomFile - We do not really upload a file to S3
//
func GetRandomFile(accountId int64) models.File {
	rand.Seed(time.Now().UnixNano())

	name := helpers.RandStr(10) + ".jpeg"

	f := models.File{
		AccountId:        uint(accountId),
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
		Host:             "amazon-s3",
		Name:             name,
		Path:             fmt.Sprintf("accounts/%d/55_%s", accountId, name),
		ThumbPath:        fmt.Sprintf("accounts/%d/55_thumb_%s", accountId, name),
		Type:             "image/jpeg",
		Hash:             "6f27495962c7e17bcf5352cdc142b26a",
		Size:             1368944,
		Url:              "",
		Thumb600By600Url: "",
	}

	// Return Happy
	return f
}

/* End File */
