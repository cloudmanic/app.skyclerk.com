//
// Date: 9/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package slack

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

//
// Test - Notify
//
func TestNotify(t *testing.T) {
	// Get the path to the .env file relative to this test file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	envPath := filepath.Join(basepath, "..", "..", ".env")
	
	// Load .env file
	env.ReadEnv(envPath)

	if len(os.Getenv("SLACK_HOOK")) > 0 {

		// Flush pending mocks after test execution
		defer gock.Off()

		// Setup mock request.
		gock.New(os.Getenv("SLACK_HOOK")).
			Reply(200).
			BodyString("ok")

		// Run a test webhook to slack.
		res, err := Notify("#events", "Test from unit testing")

		if err != nil {
			panic(err)
		}

		// Verify the data was return as expected
		st.Expect(t, res, "ok")

		// Verify that we don't have pending mocks
		st.Expect(t, gock.IsDone(), true)

	}

}

/* End File */
