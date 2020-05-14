//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Account implements Serializable {
	Id: number = 0;
	Name: string = "";
	OwnerId: number = 0;
	Locale: string = "en-US";
	Currency: string = "USD";


	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Name = json["name"];
		this.OwnerId = json["owner_id"];
		this.Locale = json["locale"];
		this.Currency = json["currency"];
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Account): Object {
		let rt = {
			id: obj.Id,
			name: obj.Name,
			owner_id: obj.OwnerId,
			locale: obj.Locale,
			currency: obj.Currency,
		}
		return rt;
	}
}

/* End File */
