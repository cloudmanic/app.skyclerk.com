//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';
import { Account } from './account.model';

export class Me implements Serializable {
	Id: number;
	Email: string;
	FirstName: string;
	LastName: string;
	Accounts: Account[];

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Email = json["email"];
		this.FirstName = json["first_name"];
		this.LastName = json["last_name"];
		this.Accounts = [];

		// Add in the accounts.
		for (let i = 0; i < json["accounts"].length; i++) {
			this.Accounts.push(new Account().deserialize(json["accounts"][i]));
		}

		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Me): Object {
		let rt = {
			id: obj.Id,
			email: obj.Email,
			first_name: obj.FirstName,
			last_name: obj.LastName,
			accounts: []
		}

		// Add in the accounts.
		for (let i = 0; i < obj.Accounts.length; i++) {
			rt.accounts.push(new Account().serialize(obj.Accounts[i]));
		}

		return rt;
	}
}

/* End File */
