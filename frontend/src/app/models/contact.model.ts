//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Serializable } from './serializable.model';

export class Contact implements Serializable {
	Id: number = 0;
	AccountId: number = 0;
	Name: string = "";
	FirstName: string = "";
	LastName: string = "";
	Email: string = "";
	AvatarUrl: string = "";
	Address: string = "";
	City: string = "";
	State: string = "";
	Zip: string = "";
	Phone: string = "";
	Fax: string = "";
	Website: string = "";
	AccountNumber: string = "";
	Twitter: string = "";
	Facebook: string = "";
	Linkedin: string = "";
	Country: string = "";

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
		this.AvatarUrl = json["avatar_url"];
		this.Address = json["address"];
		this.City = json["city"];
		this.State = json["state"];
		this.Zip = json["zip"];
		this.Phone = json["phone"];
		this.Fax = json["fax"];
		this.Website = json["website"];
		this.AccountNumber = json["account_number"];
		this.Twitter = json["twitter"];
		this.Facebook = json["facebook"];
		this.Linkedin = json["linkedin"];
		this.Country = json["country"];
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
			email: obj.Email,
			address: obj.Address,
			city: obj.City,
			state: obj.State,
			zip: obj.Zip,
			phone: obj.Phone,
			fax: obj.Fax,
			website: obj.Website,
			account_number: obj.AccountNumber,
			twitter: obj.Twitter,
			facebook: obj.Facebook,
			linkedin: obj.Linkedin,
			country: obj.Country
		}
		return rt;
	}
}

/* End File */
