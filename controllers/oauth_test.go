//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"testing"

	"github.com/cloudmanic/skyclerk.com/library/test"
	"github.com/davecgh/go-spew/spew"
)

//
// TestDoOauthToken01 test to make sure auth works.
//
func TestDoOauthToken01(t *testing.T) {

	user := test.GetRandomUser(33)

	spew.Dump(user)

	app := test.GetRandomApplication()

	spew.Dump(app)

}

/* End File */
