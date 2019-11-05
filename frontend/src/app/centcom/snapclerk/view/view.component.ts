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
import { Me } from 'src/app/models/me.model';
import { Account } from 'src/app/models/account.model';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	user: Me = new Me();
	account: Account = new Account();
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
		this.loadUser();
	}

	//
	// On account change.
	//
	onAccountChange() {
		this.contactsResults = [];
		this.loadContacts();
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
	// Load user
	//
	loadUser() {
		this.requestUser().subscribe(res => {
			this.user = res;
			this.account = this.user.Accounts[0];
		});
	}

	//
	// Request User
	//
	requestUser(): Observable<Me> {
		return this.http.get<Me>(environment.app_server + '/api/admin/users/109')
			.pipe(map(res => { return new Me().deserialize(res); }));
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
		let url = `${environment.app_server}/api/admin/contacts?limit=20&account_id=${this.account.Id}`;

		// Are we searching?
		if (this.contactInput.length > 0) {
			url = url + "&search=" + this.contactInput
		}

		return this.http.get<Contact[]>(url)
			.pipe(map(res => res.map(res => new Contact().deserialize(res))));
	}
}

/* End File */
