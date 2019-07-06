//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';

export class Activity implements Serializable {
	Id: number = 0;
	Name: string = "";
	Action: string = "";
	SubAction: string = "";
	Message: string = "";
	LedgerId: number = 0;
	ContactId: number = 0;
	LabelId: number = 0;
	CategoryId: number = 0;
	SnapSnapClerkId: number = 0;
	CreatedAt: Date = new Date();

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Name = json["name"];
		this.Action = json["action"];
		this.SubAction = json["sub_action"];
		this.Message = json["message"];
		this.LedgerId = json["ledger_id"];
		this.ContactId = json["contact_id"];
		this.LabelId = json["label_id"];
		this.CategoryId = json["category_id"];
		this.SnapSnapClerkId = json["snapclerk_id"];
		this.CreatedAt = moment(json["created_at"]).toDate();
		return this;
	}

	//
	// Model to JS Object. - We do not send activity back to the server.
	//
	serialize(obj: Activity): Object {
		let rt = {}
		return rt;
	}
}

/* End File */
