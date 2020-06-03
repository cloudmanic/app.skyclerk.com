//
// Date: 8/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package sendy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"app.skyclerk.com/backend/services"
)

//
// IsUnSubscribed Get a subscriber's status
//
func IsUnSubscribed(listId string, email string) (bool, error) {

	listIdString := GetListId(listId)

	// Build form request
	form := url.Values{
		"api_key": {os.Getenv("SENDY_API_KEY")},
		"list_id": {listIdString},
		"email":   {email},
	}

	// Send request.
	resp, err := http.PostForm("https://sendy.cloudmanic.com/api/subscribers/subscription-status.php", form)

	if err != nil {
		services.Info(errors.New("SendySubscriberStatus - Unable to get subscribe status " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		services.Info(errors.New("SendySubscriberStatus (no 200) - Unable to get subscribe status " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
		return false, err
	}

	defer resp.Body.Close()

	// Read the data we got.
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	if string(body) == "Unsubscribed" {
		return true, nil
	}

	// Return happy
	return false, nil
}

//
// Subscribe to a sendy newsletter list
//
func Subscribe(listId string, email string, first string, last string, ip string, paid string) {

	listIdString := GetListId(listId)

	// Build form request
	form := url.Values{
		"list":      {listIdString},
		"email":     {email},
		"name":      {first + " " + last},
		"FirstName": {first},
		"LastName":  {last},
		"api_key":   {os.Getenv("SENDY_API_KEY")},
		"Paid":      {paid},
	}

	// Do we have an ip address
	if len(ip) > 0 {
		form["ipaddress"] = []string{ip}
	}

	// Check to see if this user is unsubscripted.
	unsubscribed, err := IsUnSubscribed(listId, email)

	if err != nil {
		services.Info(errors.New("SendySubscribe - Unable to check if unsubscribed " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
	}

	// Log
	services.Info(errors.New("Subscribing " + email + " to Sendy List - " + listIdString))

	// Send request.
	resp, err := http.PostForm("https://sendy.cloudmanic.com/subscribe", form)

	if err != nil {
		services.Info(errors.New("SendySubscribe - Unable to subscribe " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
	}

	if resp.StatusCode != http.StatusOK {
		services.Info(errors.New("SendySubscribe (no 200) - Unable to subscribe " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
	}

	// If the user was already unsubscripted set it back to unsubbed. We often just want to update a var with the subscriber. The action above will subscribe them.
	if unsubscribed {
		Unsubscribe(listId, email)
	}

	defer resp.Body.Close()

}

//
// Unsubscribe to a sendy newsletter list
//
func Unsubscribe(listId string, email string) {

	listIdString := GetListId(listId)

	// Build form request
	form := url.Values{
		"list":    {listIdString},
		"email":   {email},
		"api_key": {os.Getenv("SENDY_API_KEY")},
	}

	// Log
	services.Info(errors.New("Unsubscribing " + email + " to Sendy List - " + listIdString))

	// Send request.
	resp, err := http.PostForm("https://sendy.cloudmanic.com/unsubscribe", form)

	if err != nil {
		services.Info(errors.New("SendySubscribe - Unable to unsubscribe " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
	}

	if resp.StatusCode != http.StatusOK {
		services.Info(errors.New("SendySubscribe (no 200) - Unable to unsubscribe " + email + " to Sendy Subscriber list. (" + err.Error() + ")"))
	}

	defer resp.Body.Close()

}

//
// Get the list id.
//
func GetListId(listId string) string {

	var listIdString = ""

	// Get the proper list id from our configs
	switch listId {

	case "expired":
		if len(os.Getenv("SENDY_EXPIRED_LIST")) > 0 {
			listIdString = os.Getenv("SENDY_EXPIRED_LIST")
		}

	case "subscribers":
		if len(os.Getenv("SENDY_SUBSCRIBE_LIST")) > 0 {
			listIdString = os.Getenv("SENDY_SUBSCRIBE_LIST")
		}

	}

	// Make sure we have a list id.
	if len(listIdString) == 0 {
		services.Critical(errors.New("No listIdString found in SendySubscribe : " + listId + " - " + listIdString))
		return ""
	}

	return listIdString
}

/* End File */
