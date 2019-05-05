//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Account implements Serializable {
	Id: number;
	Name: string;

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Name = json["name"];
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Account): Object {
		let rt = {
			id: obj.Id,
			name: obj.Name
		}
		return rt;
	}
}

/* End File */
