//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Category implements Serializable {
	Id: number;
	AccountId: number;
	Name: string;
	Type: string;
	Count: number;
	EditMode: boolean = false;
	ErrorMsg: string = "";

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.AccountId = json["account_id"];
		this.Name = json["name"];
		this.Type = json["type"];
		this.Count = json["count"];

		if (json["type"] == "1") {
			this.Type = "expense";
		}

		if (json["type"] == "2") {
			this.Type = "income";
		}

		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Category): Object {

		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			name: obj.Name,
			type: obj.Type,
			count: obj.Count
		}

		if (rt.type == "expense") {
			rt.type = "1";
		}

		if (rt.type == "income") {
			rt.type = "2";
		}

		return rt;
	}
}

/* End File */
