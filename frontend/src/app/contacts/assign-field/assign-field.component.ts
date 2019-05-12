//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { ContactService } from 'src/app/services/contact.service';
import { Contact } from 'src/app/models/contact.model';

@Component({
	selector: 'app-contacts-assign-field',
	templateUrl: './assign-field.component.html'
})
export class AssignFieldComponent implements OnInit {
	@Output() addContactToggle = new EventEmitter<boolean>();

	selectedContact: Contact = null;
	showAddContact: boolean = false;
	contactSearchTerm: string = "";
	contactSearchResults: Contact[] = [];

	//
	// Constructor
	//
	constructor(public contactService: ContactService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Clear contact and start over
	//
	changeContact() {
		this.contactSearchTerm = "";
		this.selectedContact = null;
	}

	//
	// Select a contact for this field
	//
	setContact(contact: Contact) {
		this.selectedContact = contact;
		this.contactSearchResults = [];
	}

	//
	// Call this on key up on contact search.
	//
	searchForContact() {
		if (this.contactSearchTerm == "") {
			this.contactSearchResults = [];
			return;
		}

		this.contactService.get(50, this.contactSearchTerm).subscribe(res => {
			this.contactSearchResults = res;
		});
	}

	//
	// Toggle the add contact popover
	//
	addContactToggleClick() {
		this.addContactToggle.emit(true);
	}

}

/* End File */
