//
// Date: 2019-04-23
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';
import { File as FileModel } from './file.model';

export class SnapClerk implements Serializable {
	Id: number;
	AccountId: number;
	Status: string;
	File: FileModel;
	LedgerId: number;
	Amount: number;
	Contact: string;
	Category: string;
	Labels: string;
	Note: string;
	Lat: string;
	Lon: string;
	CreatedAt: Date;
	ProcessedAt: Date;

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.AccountId = json["account_id"];
		this.Status = json["status"];
		this.File = new FileModel().deserialize(json["file"]);
		this.LedgerId = json["ledger_id"];
		this.Amount = json["amount"];
		this.Contact = json["contact"];
		this.Category = json["category"];
		this.Labels = json["labels"];
		this.Note = json["note"];
		this.Lat = json["lat"];
		this.Lon = json["lon"];
		this.CreatedAt = moment(json["created_at"]).toDate();
		this.ProcessedAt = moment(json["processed_at"]).toDate();
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: SnapClerk): Object {
		let rt = {
			id: obj.Id,
			account_id: obj.AccountId,
			status: obj.Status,
			file: new FileModel().serialize(obj.File),
			ledger_id: obj.LedgerId,
			amount: obj.Amount,
			contact: obj.Contact,
			category: obj.Category,
			labels: obj.Labels,
			note: obj.Note,
			lat: obj.Lat,
			lon: obj.Lon
		}

		return rt;
	}
}

/* End File */
