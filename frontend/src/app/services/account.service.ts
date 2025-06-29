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
import { Me } from '../models/me.model';
import { User } from '../models/user.model';
import { Billing } from '../models/billing.model';
import { BillingInvoice } from '../models/billing-invoice.model';

@Injectable({
	providedIn: 'root'
})

export class AccountService {
	// Used when account change happens.
	accountChange = new EventEmitter<number>();

	activeAccount: Account = new Account();

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

	//
	// Delete account - delete account.
	//
	delete(): Observable<Account[]> {
		let accountId = localStorage.getItem('account_id');

		return this.http.post<Account[]>(`${environment.app_server}/api/v3/${accountId}/account/delete`, {})
			.pipe(map((res) => {
				// Track event.
				this.trackService.event('account-delete', { app: "web", "accountId": accountId });

				let a = [];

				for (let i = 0; i < res.length; i++) {
					a.push(new Account().deserialize(res[i]));
				}

				return a;
			}));
	}

	//
	// New account - create new account.
	//
	create(name: string): Observable<Account> {
		let accountId = localStorage.getItem('account_id');

		return this.http.post<Account>(`${environment.app_server}/api/v3/${accountId}/account/new`, { name: name })
			.pipe(map(res => {
				let a = new Account().deserialize(res);

				// Track event.
				this.trackService.event('account-new', { app: "web", "accountId": a.Id });

				return a;
			}));
	}

	//
	// stripeToken - send token to the backend
	//
	stripeToken(token: string, plan: string): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.post<Boolean>(`${environment.app_server}/api/v3/${accountId}/account/stripe-token`, { token: token, plan: plan })
			.pipe(map(_res => {
				// Track event.
				this.trackService.event('account-stripe-token', { app: "web", "accountId": accountId });

				return true;
			}));
	}

	//
	// Get billing
	//
	getBilling(): Observable<Billing> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/account/billing`;
		return this.http.get<Billing>(url).pipe(map(res => new Billing().deserialize(res)));
	}

	//
	// Get billing history
	//
	getBillingHistory(): Observable<BillingInvoice[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/account/billing-history`;
		return this.http.get<BillingInvoice[]>(url)
			.pipe(map(res => res.map(res => new BillingInvoice().deserialize(res))));
	}

	//
	// updateSubscription - update subscription monthly / yearly
	//
	updateSubscription(plan: string): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.put<Boolean>(`${environment.app_server}/api/v3/${accountId}/account/subscription`, { plan: plan })
			.pipe(map(_res => {
				// Track event.
				this.trackService.event('account-subscription-change', { app: "web", "accountId": accountId });

				return true;
			}));
	}
}

/* End File */
