//
// Date: 2020-05-13
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
	selector: 'app-users',
	templateUrl: './users.component.html'
})

export class UsersComponent implements OnInit {
	accounts: Account[] = [];

	//
	// Constructor.
	//
	constructor(public http: HttpClient) {
		this.getAccounts();
	}

	//
	// ngOnInit.
	//
	ngOnInit() { }

	//
	// getAccounts from the backend
	//
	getAccounts() {
		this.requestAccounts().subscribe(res => {
			this.accounts = res;
		});
	}

	//
	// Request accounts
	//
	requestAccounts(): Observable<Account[]> {
		return this.http.get<Account[]>(environment.app_server + '/api/admin/accounts')
			.pipe(map(res => res.map(res => new Account().deserialize(res))));
	}
}

// Accounts returned from admin API
export class Account {
	AccountId: number = 0;
	CreatedAt: Date = new Date();
	LastActivity: Date = new Date();
	Name: string = "";
	FirstName: string = "";
	LastName: string = "";
	Email: string = "";
	Status: string = "";
	LedgerCount: number = 0;

	//
	// Json to Object.
	//
	deserialize(json: Account): this {
		this.AccountId = json["account_id"];
		this.CreatedAt = moment(json["created_at"]).toDate();
		this.LastActivity = moment(json["last_activity"]).toDate();
		this.Name = json["name"];
		this.FirstName = json["first_name"];
		this.LastName = json["last_name"];
		this.Email = json["email"];
		this.Status = json["status"];
		this.LedgerCount = json["ledger_count"];
		return this;
	}
}

/* End File */
