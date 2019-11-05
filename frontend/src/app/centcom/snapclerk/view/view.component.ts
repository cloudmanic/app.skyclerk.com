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
import { SnapClerk } from 'src/app/models/snapclerk.model';
import { Category } from 'src/app/models/category.model';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	user: Me = new Me();
	remaining: number = 0;
	account: Account = new Account();
	snapclerk: SnapClerk = new SnapClerk();
	contact: Contact = new Contact();
	contactInput: string = "";
	contactsResults: Contact[] = [];
	category: Category = new Category();
	categoriesInput: string = "";
	categoriesResults: Category[] = [];

	//
	// Constructor
	//
	constructor(public http: HttpClient) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		this.loadSnapClerks();
	}

	//
	// Save snapclerk
	//
	save() {
		this.snapclerk.Contact = this.contactInput;
		this.snapclerk.Category = this.categoriesInput;

		console.log(this.snapclerk);
	}

	//
	// On account change.
	//
	onAccountChange() {
		this.contactsResults = [];
		this.loadContacts();
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
	// Load SnapClerk
	//
	loadSnapClerks() {
		this.requestSnapClerks().subscribe(res => {
			if (res.length == 0) {
				this.remaining = 0;
				this.snapclerk = new SnapClerk();
				return;
			}

			// Set active snapclerk
			this.snapclerk = res[0];
			this.remaining = (res.length - 1);

			// Set the account.
			this.loadUser();
		});
	}

	//
	// Set active account.
	//
	setActiveAccount() {
		for (let i = 0; i < this.user.Accounts.length; i++) {
			if (this.snapclerk.AccountId == this.user.Accounts[i].Id) {
				this.account = this.user.Accounts[i];
				this.setActiveCategory();
				return;
			}
		}
	}

	//
	// Set active category
	//
	setActiveCategory() {
		// Get all categories.
		this.getAllCategories().subscribe(res => {
			for (let i = 0; i < res.length; i++) {
				if (res[i].Name == this.snapclerk.Category) {
					this.category = res[i];
					this.categoriesInput = res[i].Name;
					return;
				}
			}
		});

		// Must be a new category.
		this.category = new Category();
		this.category.Name = this.snapclerk.Category;
		this.category.Type = "expense";
	}

	//
	// Get all categories
	//
	getAllCategories() {
		let url = `${environment.app_server}/api/admin/categories?account_id=${this.account.Id}&type=expense`;

		return this.http.get<Category[]>(url)
			.pipe(map(res => res.map(res => new Category().deserialize(res))));
	}

	//
	// Request Snapclerks
	//
	requestSnapClerks(): Observable<SnapClerk[]> {
		return this.http.get<SnapClerk[]>(environment.app_server + '/api/admin/snapclerk')
			.pipe(map(res => res.map(res => new SnapClerk().deserialize(res))));
	}

	//
	// Load user
	//
	loadUser() {
		this.requestUser().subscribe(res => {
			this.user = res;
			this.account = this.user.Accounts[0];
			this.setActiveAccount();
		});
	}

	//
	// Request User
	//
	requestUser(): Observable<Me> {
		return this.http.get<Me>(environment.app_server + '/api/admin/users/' + this.snapclerk.AddedById)
			.pipe(map(res => { return new Me().deserialize(res); }));
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

	//
	// onCategoryChange do search
	//
	onCategoryChange() {
		this.loadCategories();
	}

	//
	// onCategorySelect
	//
	onCategorySelect(result: Category) {
		this.category = result;
		this.categoriesInput = result.Name;
		this.categoriesResults = [];
	}

	//
	// Load categories.
	//
	loadCategories() {
		// Do nothing is search is empty.
		if (this.categoriesInput.length == 0) {
			this.categoriesResults = [];
			return;
		}

		// Search Categories list.
		this.requestCategories().subscribe(res => {
			this.categoriesResults = res;
		});
	}

	//
	// Request categories
	//
	requestCategories(): Observable<Category[]> {
		let url = `${environment.app_server}/api/admin/categories?account_id=${this.account.Id}&type=expense`;

		// Are we searching?
		if (this.categoriesInput.length > 0) {
			url = url + "&search=" + this.categoriesInput
		}

		return this.http.get<Category[]>(url)
			.pipe(map(res => res.map(res => new Category().deserialize(res))));
	}
}

/* End File */
