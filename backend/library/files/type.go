//
// Date: 2019-04-19
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package files

import (
	"io/ioutil"

	"github.com/h2non/filetype"
)

//
// FileContentTypeWithError - Return type filt hash of a file with error.
//
func FileContentTypeWithError(path string) (string, string, error) {
	// Read file
	buf, _ := ioutil.ReadFile(path)

	// Get kind
	kind, err := filetype.Match(buf)
	if kind == filetype.Unknown {
		return "", "", err
	}

	// Return happy
	return kind.MIME.Value, kind.Extension, nil
}

/* End File */
