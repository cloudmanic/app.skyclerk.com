//
// Date: 2018-04-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-04-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"time"

	"github.com/araddon/dateparse"
)

//
// ParseDateNoError Parse a date return no error. Use this when we know
// there will not be errors.
//
func ParseDateNoError(d string) time.Time {
	t, _ := dateparse.ParseLocal(d)
	return t
}

/* End File */
