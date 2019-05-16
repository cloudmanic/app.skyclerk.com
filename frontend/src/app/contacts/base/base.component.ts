//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { Contact } from 'src/app/models/contact.model';
import { ContactService } from 'src/app/services/contact.service';

@Component({
	selector: 'app-base',
	templateUrl: './base.component.html'
})

export class BaseComponent implements OnInit {
	@Output() onContact = new EventEmitter<Contact>();

	fieldErrors: any;
	contact: Contact = new Contact();

	fields: Field[] = [
		{ Name: "FirstName", Title: "First Name", Field: "first_name", Type: "text", Placeholder: "", Section: 1 },
		{ Name: "LastName", Title: "Last Name", Field: "last_name", Type: "text", Placeholder: "", Section: 1 },
		{ Name: "Name", Title: "Company", Field: "name", Type: "text", Placeholder: "", Section: 1 },

		{ Name: "Email", Title: "Email", Field: "email", Type: "text", Placeholder: "", Section: 2 },
		{ Name: "Address", Title: "Address", Field: "address", Type: "text", Placeholder: "", Section: 2 },
		{ Name: "City", Title: "City", Field: "city", Type: "text", Placeholder: "", Section: 2 },

		// We hard code this so these are ignored in the popup contact
		{ Name: "State", Title: "State", Field: "state", Type: "text", Placeholder: "", Section: 3 },
		{ Name: "Zip", Title: "Zip", Field: "zip", Type: "text", Placeholder: "", Section: 3 },

		{ Name: "Country", Title: "Country", Field: "country", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Phone", Title: "Phone", Field: "phone", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Fax", Title: "Fax", Field: "fax", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Website", Title: "Website", Field: "website", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "AccountNumber", Title: "Account #", Field: "account_number", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Facebook", Title: "Facebook", Field: "facebook", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Twitter", Title: "Twitter", Field: "twitter", Type: "text", Placeholder: "", Section: 4 },
		{ Name: "Linkedin", Title: "Linkedin", Field: "linkedin", Type: "text", Placeholder: "", Section: 4 }
	];

	//
	// Constructor
	//
	constructor(public contactService: ContactService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Save a new contact
	//
	save() {
		this.fieldErrors = null;

		// Send request to BE to save.
		this.contactService.create(this.contact).subscribe(
			res => {
				this.contact = new Contact();
				this.onContact.emit(res);
			},

			err => {
				this.fieldErrors = err.error.errors;
			}
		);
	}

	//
	// Does an field have an error?
	//
	hasError(field: string): string {

		for (let row in this.fieldErrors) {
			if (row == field) {
				return this.fieldErrors[row];
			}
		}

		return "";
	}

}

// Interface for a field.
export interface Field {
	Name: string,
	Title: string,
	Field: string,
	Type: string,
	Placeholder: string,
	Section: number
}

/* End File */
