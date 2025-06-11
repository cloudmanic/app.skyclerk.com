//
// Date: 8/27/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package slack

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"app.skyclerk.com/backend/services"
)

// init sets default environment variables for slack functionality during tests
func init() {
	// Only set defaults during tests
	if flag.Lookup("test.v") != nil {
		setDefaultIfEmpty("APP_ENV", "test")
		setDefaultIfEmpty("SLACK_URL", "")
	}
}

// setDefaultIfEmpty sets an environment variable to a default value if it's not already set
func setDefaultIfEmpty(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}

//
// Notify to slack
//
func Notify(channel string, msg string) (string, error) {

	// Modfy message if this is testing.
	if os.Getenv("APP_ENV") == "local" {
		msg = "(local dev): " + msg
	}

	if len(os.Getenv("SLACK_HOOK")) > 0 {

		var jsonStr = []byte(`{"channel": "` + channel + `", "text": "` + msg + `"}`)

		// Creatre POST request
		req, err := http.NewRequest("POST", os.Getenv("SLACK_HOOK"), bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		// Send request.
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			services.Critical(errors.New(err.Error() + "Notify - Unable to send slack notice : " + msg + "."))
			return "", err
		}

		if resp.StatusCode != http.StatusOK {
			services.Critical(errors.New(err.Error() + "Notify (no 200) - Unable to send slack notice : " + msg + "."))
			return "", err
		}

		// Get the body.
		body, _ := ioutil.ReadAll(resp.Body)

		resp.Body.Close()

		// Return happy.
		return string(body), err

	}

	// Nothing happened.
	return "", errors.New("SLACK_HOOK is not set.")

}

/* End File */
