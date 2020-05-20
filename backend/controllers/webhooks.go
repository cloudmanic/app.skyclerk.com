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

	"app.skyclerk.com/backend/library/slack"
	"app.skyclerk.com/backend/services"
)

const foveaSecret = "01e29cbc-f132-442e-a4f1-bcf4f95ce49f"

//
// DoPostmarkWebhook will process a webhook from postmarkapp.com
//
func (t *Controller) DoPostmarkWebhook(c *gin.Context) {

	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	to := gjson.Get(string(body), "To").String()

	// TODO(spicer): Use the postmark API to verify this is a real email posted from them. From message ID

	// Webhook received from fovea
	services.InfoMsg(string(body))

	// Send Slack hook.
	go slack.Notify("#events", fmt.Sprintf("Skyclerk DoPostmarkWebhook Webhook: To Email: %s", to))

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

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

	// Send Slack hook.
	go slack.Notify("#events", fmt.Sprintf("Skyclerk DoFoveaWebhook Webhook: Email: %s", accountString))

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/*
{
  "type": "purchases.updated",
  "applicationUsername": "ios@skyclerk.com",
  "purchases": {
    "apple:monthly_6": {
      "productId": "apple:monthly_6",
      "platform": "apple",
      "sandbox": false,
      "purchaseId": "apple:70000786633723",
      "purchaseDate": "2020-05-19T21:05:40.000Z",
      "lastRenewalDate": "2020-05-19T21:05:37.000Z",
      "expirationDate": "2020-06-19T21:05:37.000Z",
      "isTrialPeriod": false,
      "isIntroPeriod": false,
      "renewalIntent": "Renew",
      "isExpired": false
    }
  },
  "password": "01e29cbc-f132-442e-a4f1-bcf4f95ce49f"
}
*/

// TODO(spicer): unit tests for DoFoveaWebhook
// {"type":"purchases.updated","applicationUsername":"5127 - spicer@cloudmanic.com","purchases":{"apple:monthly_6":{"productId":"apple:monthly_6","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665621214","purchaseDate":"2020-05-15T20:39:10.000Z","lastRenewalDate":"2020-05-16T02:30:56.000Z","expirationDate":"2020-05-16T02:35:56.000Z","isTrialPeriod":false,"isIntroPeriod":false,"isBillingRetryPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_CHANGE_RENEWAL_PREF","isExpired":true},"apple:yearly_60":{"productId":"apple:yearly_60","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665621214","purchaseDate":"2020-05-15T20:39:10.000Z","lastRenewalDate":"2020-05-16T02:39:08.000Z","expirationDate":"2020-05-16T03:39:08.000Z","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_CHANGE_RENEWAL_STATUS","renewalIntentChangeDate":"2020-05-16T02:39:08.000Z","isExpired":false}},"password":"01e29cbc-f132-442e-a4f1-bcf4f95ce49f"}

// {"type":"purchases.updated","applicationUsername":"5127:::spicer@cloudmanic.com","purchases":{"apple:monthly_6":{"productId":"apple:monthly_6","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665704156","purchaseDate":"2020-05-16T03:24:13.000Z","lastRenewalDate":"2020-05-16T03:52:15.000Z","expirationDate":"2020-05-16T03:57:15.000Z","cancelationReason":"System.Replaced","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","lastNotification":"DID_RECOVER","isExpired":false},"apple:yearly_60":{"productId":"apple:yearly_60","platform":"apple","sandbox":true,"purchaseId":"apple:1000000665704156","purchaseDate":"2020-05-16T03:24:13.000Z","lastRenewalDate":"2020-05-16T03:56:37.000Z","expirationDate":"2020-05-16T04:56:37.000Z","isTrialPeriod":false,"isIntroPeriod":false,"renewalIntent":"Renew","isExpired":false}},"password":"01e29cbc-f132-442e-a4f1-bcf4f95ce49f"}

/*
{
    "FromName": "Postmarkapp Support",
    "MessageStream": "inbound",
    "From": "support@postmarkapp.com",
    "FromFull": {
        "Email": "support@postmarkapp.com",
        "Name": "Postmarkapp Support",
        "MailboxHash": ""
    },
    "To": "\"Firstname Lastname\" <mailbox+SampleHash@inbound.postmarkapp.com>",
    "ToFull": [
        {
            "Email": "mailbox+SampleHash@inbound.postmarkapp.com",
            "Name": "Firstname Lastname",
            "MailboxHash": "SampleHash"
        }
    ],
    "Cc": "\"First Cc\" <firstcc@postmarkapp.com>, secondCc@postmarkapp.com",
    "CcFull": [
        {
            "Email": "firstcc@postmarkapp.com",
            "Name": "First Cc",
            "MailboxHash": ""
        },
        {
            "Email": "secondCc@postmarkapp.com",
            "Name": "",
            "MailboxHash": ""
        }
    ],
    "Bcc": "\"First Bcc\" <firstbcc@postmarkapp.com>",
    "BccFull": [
        {
            "Email": "firstbcc@postmarkapp.com",
            "Name": "First Bcc",
            "MailboxHash": ""
        }
    ],
    "OriginalRecipient": "mailbox+SampleHash@inbound.postmarkapp.com",
    "Subject": "Test subject",
    "MessageID": "00000000-0000-0000-0000-000000000000",
    "ReplyTo": "replyto@example.com",
    "MailboxHash": "SampleHash",
    "Date": "Wed, 20 May 2020 18:33:57 -0400",
    "TextBody": "This is a test text body.",
    "HtmlBody": "<html><body><p>This is a test html body.</p></body></html>",
    "StrippedTextReply": "This is the reply text",
    "Tag": "TestTag",
    "Headers": [
        {
            "Name": "X-Header-Test",
            "Value": ""
        }
    ],
    "Attachments": [
        {
            "Name": "test.txt",
            "Content": "VGhpcyBpcyBhdHRhY2htZW50IGNvbnRlbnRzLCBiYXNlLTY0IGVuY29kZWQu",
            "ContentType": "text/plain",
            "ContentLength": 45
        }
    ]
}
*/

/* End File */
