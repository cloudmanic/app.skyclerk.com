//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';
import { Contact } from './contact.model';
import { Category } from './category.model';

export class Ledger implements Serializable {
	Id: number;
	AccountId: number;
	Date: Date;
	Amount: number;
	Note: string;
	Contact: Contact;
	Category: Category;

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.AccountId = json["account_id"];
		this.Date = moment(json["date"]).toDate();
		this.Amount = json["amount"];
		this.Note = json["note"];
		this.Contact = new Contact().deserialize(json["contact"]);
		this.Category = new Category().deserialize(json["category"]);
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Ledger): Object {
		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			date: obj.Date,
			amount: obj.Amount,
			note: obj.Note,
			contact: new Contact().serialize(obj.Contact),
			category: new Category().serialize(obj.Category)
		}
		return rt;
	}
}

/* End File */
