//
// Date: 2018-07-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-16
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import "sort"

//
// Round
//
func Round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}

//
// Find Min Float64 in a Slice
//
func MinFloat64Slice(v []float64) float64 {
	sort.Float64s(v)
	return v[0]
}

//
// Find Max Float64 in a Slice
//
func MaxFloat64Slice(v []float64) float64 {
	sort.Float64s(v)
	return v[len(v)-1]
}

//
// Find Min Int in a Slice
//
func MinIntSlice(v []int) int {
	sort.Ints(v)
	return v[0]
}

//
// Find Max Int in a Slice
//
func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}

/* End File */
