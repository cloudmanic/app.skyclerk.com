//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Contact implements Serializable {
	Id: number;
	AccountId: number;
	Name: string;
	FirstName: string;
	LastName: string;
	Email: string;

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.AccountId = json["account_id"];
		this.Name = json["name"];
		this.FirstName = json["first_name"];
		this.LastName = json["last_name"];
		this.Email = json["email"];
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Contact): Object {
		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			name: obj.Name,
			first_name: obj.FirstName,
			last_name: obj.LastName,
			email: obj.Email
		}
		return rt;
	}
}

/* End File */
