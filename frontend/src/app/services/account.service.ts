//
// Date: 2019-09-20
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { Account } from '../models/account.model';
import { TrackService } from './track.service';

@Injectable({
	providedIn: 'root'
})

export class AccountService {
	// Used when account change happens.
	accountChange = new EventEmitter<number>();

	public activeAccount: Account = new Account();

	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) {
		this.setActiveAccount();
	}

	//
	// Get the active account.
	//
	getActiveAccount(): Account {
		return this.activeAccount;
	}

	//
	// Set active account.
	//
	setActiveAccount() {
		// Get the active account.
		this.getAccount().subscribe(res => {
			this.activeAccount = res;
		});
	}

	//
	// Get by ID
	//
	getAccount(): Observable<Account> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/account`;
		return this.http.get<Account>(url).pipe(map(res => new Account().deserialize(res)));
	}

	//
	// Update the account update
	//
	update(acct: Account): Observable<Account> {
		let accountId = localStorage.getItem('account_id');
		acct.Id = Number(accountId);

		return this.http.put<Account>(`${environment.app_server}/api/v3/${accountId}/account`, new Account().serialize(acct))
			.pipe(map(res => {
				let a = new Account().deserialize(res);

				// Track event.
				this.trackService.event('account-update', { app: "web", "accountId": accountId });

				return a;
			}));
	}

	//
	// Clear account - clears all account data.
	//
	clear(): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.post<Boolean>(`${environment.app_server}/api/v3/${accountId}/account/clear`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('account-clear', { app: "web", "accountId": accountId });

				return true;
			}));
	}
}

/* End File */
