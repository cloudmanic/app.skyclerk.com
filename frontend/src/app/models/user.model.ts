//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';

export class User implements Serializable {
	Id: number;
	Email: string;
	FirstName: string;
	LastName: string;
	LastActivity: Date;

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Email = json["email"];
		this.FirstName = json["first_name"];
		this.LastName = json["last_name"];
		this.LastActivity = moment(json["last_activity"]).toDate();
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: User): Object {
		let rt = {
			id: obj.Id,
			email: obj.Email,
			first_name: obj.FirstName,
			last_name: obj.LastName,
			last_activity: obj.LastActivity
		}
		return rt;
	}
}

/* End File */
