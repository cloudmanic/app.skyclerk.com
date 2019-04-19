//
// Date: 2019-04-19
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package files

import (
	"os"
)

//
// SizeWithError - Return the MD5 hash of a file with error.
//
func SizeWithError(path string) (int64, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

/* End File */
