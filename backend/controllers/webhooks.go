//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"app.skyclerk.com/backend/services"
)

const foveaSecret = "01e29cbc-f132-442e-a4f1-bcf4f95ce49f"

//
// DoFoveaWebhook  will process webhooks from fovea.cc
//
func (t *Controller) DoFoveaWebhook(c *gin.Context) {

	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	password := gjson.Get(string(body), "password").String()
	accountString := gjson.Get(string(body), "applicationUsername").String()
	// currency := gjson.Get(string(body), "currency").String()
	// locale := gjson.Get(string(body), "locale").String()
	// ownerId := gjson.Get(string(body), "owner_id").Int()

	// First verify the password
	if password != foveaSecret {
		services.Critical(fmt.Errorf("DoFoveaWebhook: Password not found for Fovea webhook"))
		c.JSON(http.StatusNoContent, nil)
		return
	}

	// Webhook received from fovea
	services.InfoMsg(string(body))

	// Parse application string to figure out which appplication we are dealing with.
	fmt.Println(accountString)

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

// TODO(spicer): unit tests for DoFoveaWebhook
// {"type":"purchases.updated","applicationUsername":"5127 - spicer@cloudmanic.com","purchases":{"apple:monthly_6":{"productId":"apple:monthly_6","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665621214","purchaseDate":"2020-05-15T20:39:10.000Z","lastRenewalDate":"2020-05-16T02:30:56.000Z","expirationDate":"2020-05-16T02:35:56.000Z","isTrialPeriod":false,"isIntroPeriod":false,"isBillingRetryPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_CHANGE_RENEWAL_PREF","isExpired":true},"apple:yearly_60":{"productId":"apple:yearly_60","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665621214","purchaseDate":"2020-05-15T20:39:10.000Z","lastRenewalDate":"2020-05-16T02:39:08.000Z","expirationDate":"2020-05-16T03:39:08.000Z","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_CHANGE_RENEWAL_STATUS","renewalIntentChangeDate":"2020-05-16T02:39:08.000Z","isExpired":false}},"password":"01e29cbc-f132-442e-a4f1-bcf4f95ce49f"}

// {"type":"purchases.updated","applicationUsername":"5127:::spicer@cloudmanic.com","purchases":{"apple:monthly_6":{"productId":"apple:monthly_6","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665704156","purchaseDate":"2020-05-16T03:24:13.000Z","lastRenewalDate":"2020-05-16T03:52:15.000Z","expirationDate":"2020-05-16T03:57:15.000Z","cancelationReason":"System.Replaced","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_RECOVER","isExpired":false},"apple:yearly_60":{"productId":"apple:yearly_60","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665704156","purchaseDate":"2020-05-16T03:24:13.000Z","lastRenewalDate":"2020-05-16T03:56:37.000Z","expirationDate":"2020-05-16T04:56:37.000Z","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","isExpired":false}},"password":"01e29cbc-f132-442e-a4f1-bcf4f95ce49f"}

/* End File */