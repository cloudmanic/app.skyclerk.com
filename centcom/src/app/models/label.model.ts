//
// Date: 2019-04-27
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Label implements Serializable {
	Id: number;
	AccountId: number;
	Name: string;
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
		this.Count = json["count"];
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Label): Object {
		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			name: obj.Name,
			count: obj.Count
		}
		return rt;
	}
}

/* End File */
