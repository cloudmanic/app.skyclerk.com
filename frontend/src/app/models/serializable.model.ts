//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

export interface Serializable {
	serialize(obj: any): Object;
	deserialize(json: Object): this;
}

/* End File */
