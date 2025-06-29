package avatar

import "testing"

func TestCleanString(t *testing.T) {
	// Input should equal output
	tests := map[string]string{}
	tests["AE"] = "AE"
	tests["ae"] = "AE"
	tests["a e"] = "AE"
	tests["andrew edwards"] = "AE"
	tests["andrew   edwards"] = "AN"
	tests["a"] = "A"
	tests["123"] = "12"
	tests["A 3"] = "A3"
	tests["B 3"] = "B3"

	for k, v := range tests {
		if cleanString(k) != v {
			t.Errorf("Received '%s', was expecting '%s' from '%s'", cleanString(k), v, k)
		}
	}
}
