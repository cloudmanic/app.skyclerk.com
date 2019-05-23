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
import { Label } from './label.model';
import { File as FileModel } from './file.model';

export class Ledger implements Serializable {
	Id: number = 0;
	AccountId: number = 0;
	Date: Date = new Date();
	Amount: number = 0.00;
	Note: string = "";
	Lat: number = 0.00;
	Lon: number = 0.00;
	Contact: Contact = new Contact();
	Category: Category = new Category();
	Labels: Label[] = [];
	Files: FileModel[] = [];

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.AccountId = json["account_id"];
		this.Date = moment(json["date"]).toDate();
		this.Amount = json["amount"];
		this.Note = json["note"];
		this.Lat = json["lat"];
		this.Lon = json["lon"];
		this.Contact = new Contact().deserialize(json["contact"]);
		this.Category = new Category().deserialize(json["category"]);
		this.Labels = [];
		this.Files = [];

		// Deal with labels.
		if (json["labels"]) {
			for (let i = 0; i < json["labels"].length; i++) {
				this.Labels.push(new Label().deserialize(json["labels"][i]));
			}
		}

		// Deal with files.
		if (json["files"]) {
			for (let i = 0; i < json["files"].length; i++) {
				this.Files.push(new FileModel().deserialize(json["files"][i]));
			}
		}

		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Ledger): Object {
		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			date: moment(obj.Date).toDate(),
			amount: Number(obj.Amount),
			note: obj.Note,
			lat: obj.Lat,
			lon: obj.Lon,
			contact: new Contact().serialize(obj.Contact),
			category: new Category().serialize(obj.Category),
			labels: [],
			files: []
		}

		// Deal with labels.
		for (let i = 0; i < obj.Labels.length; i++) {
			rt.labels.push(new Label().serialize(obj.Labels[i]));
		}

		// Deal with files.
		for (let i = 0; i < obj.Files.length; i++) {
			rt.files.push(new FileModel().serialize(obj.Files[i]));
		}

		return rt;
	}
}

/* End File */
