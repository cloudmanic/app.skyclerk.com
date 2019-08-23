//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ContactService } from 'src/app/services/contact.service';
import { Contact } from 'src/app/models/contact.model';

@Component({
	selector: 'app-settings-contacts',
	templateUrl: './contacts.component.html'
})

export class ContactsComponent implements OnInit {
	errors: any = [];
	limit: number = 50;
	search: string = "";
	contacts: Contact[] = [];

	//
	// Constructor
	//
	constructor(public contactService: ContactService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		this.refreshContacts();
	}

	//
	// Refresh contacts
	//
	refreshContacts() {
		// Get categories by type.
		this.contactService.get(this.limit, this.search).subscribe(res => {
			this.contacts = res;
		});
	}

	//
	// Add contact click.
	//
	addContactClick() {
		let c = new Contact();
		c.EditMode = true;
		c.Country = "United States";
		this.contacts.unshift(c);
	}

	//
	// Save Contact.
	//
	save(row: Contact) {
		if (row.Id > 0) {
			this.updateContact(row);
		} else {
			this.createContact(row);
		}
	}

	//
	// Delete contact
	//
	deleteContact(row: Contact) {
		// Delete contact on server.
		this.contactService.delete(row).subscribe(
			// Success
			() => {
				this.refreshContacts();
			},

			// Error
			err => {
				alert(err.error.error);
			}
		);
	}

	//
	// Create Contact
	//
	createContact(row: Contact) {
		// Send change to server.
		this.contactService.create(row).subscribe(
			// Success
			() => {
				this.refreshContacts();
			},

			// Error
			err => {
				// Show errors
				this.errors = err.error.errors;
			}
		);
	}

	//
	// Update Contact
	//
	updateContact(row: Contact) {
		// Send change to server.
		this.contactService.update(row).subscribe(
			// Success
			(res) => {
				console.log(res);
				this.refreshContacts();
			},

			// Error
			err => {
				// Show errors
				this.errors = err.error.errors;
			}
		);
	}

	//
	// Toggle Edit Mode
	//
	toggleEditMode(row: Contact, open: string) {
		// Only one edit mode at a time.
		for (let i = 0; i < this.contacts.length; i++) {
			this.contacts[i].EditMode = false;
		}

		// Toggle open / close
		if (open == "open") {
			row.EditMode = true;

			// Default to United States
			if (!row.Country.length) {
				row.Country = "United States";
			}
		} else {
			row.EditMode = false;

			// Clear value if the user updated it.
			this.refreshContacts();
		}
	}

	//
	// Check to see if we have errors for a field.
	//
	hasError(field: string) {
		if (this.errors[field]) {
			return this.errors[field];
		}

		return "";
	}
}

/* End File */
