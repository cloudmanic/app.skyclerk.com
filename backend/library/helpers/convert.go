//
// Date: 3/5/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"strconv"
)

//
// Float to String
func FloatToString(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 2, 64)
}

//
// Convert a string to a int
//
func StringToInt(s string) int {
	idInt, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return idInt
}

//
// Convert a string to a uint
//
func StringToUint(s string) uint {
	idInt, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return uint(idInt)
}

//
// Convert a string to a float64
//
func StringToFloat64(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0
	}

	return num
}

/* End File */
