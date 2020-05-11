//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';

export class Billing implements Serializable {
	Id: number = 0;
	Subscription: string = "";
	Status: string = "en-US";
	TrialExpire: Date = new Date();


	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Subscription = json["subscription"];
		this.Status = json["status"];
		this.TrialExpire = moment(json["trial_expire"]).toDate();
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Billing): Object {
		let rt = {
			id: obj.Id,
			subscription: obj.Subscription,
			status: obj.Status,
			trial_expire: obj.TrialExpire,
		}
		return rt;
	}
}

/* End File */
