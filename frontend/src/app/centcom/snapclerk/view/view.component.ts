//
// Date: 2019-11-04
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { Contact } from 'src/app/models/contact.model';
import { map } from 'rxjs/operators';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	accountID: number = 101;
	contact: Contact = new Contact();
	contactInput: string = "";
	contactsResults: Contact[] = [];

	//
	// Constructor
	//
	constructor(public http: HttpClient) { }

	//
	// ngOnInit
	//
	ngOnInit() {

	}

	//
	// onContactChange do search
	//
	onContactChange() {
		this.loadContacts();
	}

	//
	// onContactSelect
	//
	onContactSelect(result: Contact) {
		this.contact = result;
		this.contactInput = this.displayContact(result);
		this.contactsResults = [];
	}

	//
	// Display contact name.
	//
	displayContact(row: Contact) {
		if (row.Name.length > 0) {
			return row.Name;
		}

		return row.FirstName + " " + row.LastName;
	}

	//
	// Load contacts.
	//
	loadContacts() {
		// Do nothing is search is empty.
		if (this.contactInput.length == 0) {
			this.contactsResults = [];
			return;
		}

		// Search contact list.
		this.requestContacts().subscribe(res => {
			this.contactsResults = res;
		});
	}

	//
	// Request contacts
	//
	requestContacts(): Observable<Contact[]> {
		let url = `${environment.app_server}/api/admin/contacts?limit=20&account_id=${this.accountID}`;

		// Are we searching?
		if (this.contactInput.length > 0) {
			url = url + "&search=" + this.contactInput
		}

		return this.http.get<Contact[]>(url)
			.pipe(map(res => res.map(res => new Contact().deserialize(res))));
	}
}

/* End File */
