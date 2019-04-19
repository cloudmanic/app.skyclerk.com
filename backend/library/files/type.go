//
// Date: 2019-04-19
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package files

import (
	"net/http"
	"os"
)

//
// FileContentTypeWithError - Return type filt hash of a file with error.
//
func FileContentTypeWithError(path string) (string, error) {
	// Open File
	f, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err2 := f.Read(buffer)

	if err2 != nil {
		return "", err2
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	// Return happy
	return contentType, nil
}

/* End File */
