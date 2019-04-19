//
// Date: 2018-03-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package files

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

//
// Return the MD5 of a file.
// Returns an empty string if there is an error
//
func Md5(path string) string {
	hash, err := Md5WithError(path)

	if err != nil {
		return ""
	}

	return hash
}

//
// Return the MD5 hash of a file with error.
//
func Md5WithError(path string) (string, error) {

	// Check the MD5 the current file.
	f, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

/* End File */
