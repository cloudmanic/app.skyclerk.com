//
// Date: 2019-03-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package test

import (
	"math/rand"
	"time"

	"github.com/cloudmanic/skyclerk.com/library/helpers"
	"github.com/cloudmanic/skyclerk.com/models"
)

//
// GetRandomAccount returns a random account.
//
func GetRandomAccount(accountId int64) models.Account {
	rand.Seed(time.Now().UnixNano())

	dates := []time.Time{
		time.Now(),
		time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC),
		time.Date(2018, 8, 19, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 1, 10, 17, 20, 01, 507451, time.UTC),
	}

	zip := []string{"13601", "97222", "12384", "97345", "97701"}
	state := []string{"NY", "OR", "TX", "NC", "SC", "VT", "CA"}
	city := []string{"Newberg", "Watertown", "New York", "Clayton", "Portland", "Road Town", "Seattle"}

	acc := models.Account{
		Id:           uint(accountId),
		OwnerId:      uint(1),
		Name:         "Name " + helpers.RandStr(16),
		PlanId:       0,
		Address:      "Address " + helpers.RandStr(16),
		City:         city[rand.Intn(len(city))],
		State:        state[rand.Intn(len(state))],
		Zip:          zip[rand.Intn(len(zip))],
		Country:      "USA",
		LastActivity: dates[rand.Intn(len(dates))],
		StripeId:     "",
		CardType:     "",
		CardLast4:    "",
		CardExpMonth: "",
		CardExpYear:  "",
		SignupIp:     "127.0.0.1",
	}

	return acc
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
	pass := "foobar"
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
		time.Now(),
		time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC),
		time.Date(2018, 8, 19, 17, 20, 01, 507451, time.UTC),
		time.Date(2019, 1, 10, 17, 20, 01, 507451, time.UTC),
	}
	amounts := []float64{1234.56, 33.44, 99.00, 555.32, 4583.01, 3.01, 0.20, 3429.34, 823.19}

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

	name := []string{"Options Cafe", "Skyclerk", "Clients", "Refund", "Non-Taxable"}

	label := models.Label{
		AccountId: uint(accountId),
		Name:      name[rand.Intn(len(name))],
	}

	return label
}

/* End File */
