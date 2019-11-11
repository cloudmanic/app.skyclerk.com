//
// Date: 2019-11-11
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

export class Numbers {
	//
	// Clean string to a float. Detect number format then convert to number.
	//
	public static toFloat(numberStr: string): number {
		// Test is this is a Brazilian style number
		let b = numberStr.match(/^\s*(?:[1-9]\d{0,2}(?:\.\d{3})*|0)(?:,\d{1,2})?$/g);

		if (b) {
			return parseFloat(b[0].replace(/\./g, "").replace(/,/g, "."));
		}

		// Must be an american style number.
		return parseFloat(numberStr.replace(/,/g, ""));
	}
}

/* End File */
