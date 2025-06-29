//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { Ledger } from '../models/ledger.model';
import { Category } from '../models/category.model';
import { Label } from '../models/label.model';
import { TrackService } from './track.service';

@Injectable({
	providedIn: 'root'
})

export class LedgerService {
	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) { }

	//
	// Create a new ledger
	//
	create(ledger: Ledger): Observable<Ledger> {
		let accountId = localStorage.getItem('account_id');
		ledger.AccountId = Number(accountId);

		return this.http.post<Ledger>(`${environment.app_server}/api/v3/${accountId}/ledger`, new Ledger().serialize(ledger))
			.pipe(map(res => {
				let lg = new Ledger().deserialize(res);
				let type = "expense";

				if (lg.Amount > 0) {
					type = "income";
				}

				// Track event.
				this.trackService.event('ledger-create', { ledgerEntryType: type, app: "web", "accountId": accountId });

				return lg;
			}));
	}

	//
	// Update a ledger
	//
	update(ledger: Ledger): Observable<Ledger> {
		let accountId = localStorage.getItem('account_id');
		ledger.AccountId = Number(accountId);

		return this.http.put<Ledger>(`${environment.app_server}/api/v3/${accountId}/ledger/${ledger.Id}`, new Ledger().serialize(ledger))
			.pipe(map(res => {
				let lg = new Ledger().deserialize(res);

				// Track that a ledger entry was updated.
				this.trackService.event('ledger-update', { app: "web", "accountId": accountId });

				return lg;
			}));
	}

	//
	// Delete a ledger
	//
	delete(ledger: Ledger): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');
		ledger.AccountId = Number(accountId);

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/ledger/${ledger.Id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('ledger-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}

	//
	// Get by ID
	//
	getById(id: number): Observable<Ledger> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/ledger/${id}`;
		return this.http.get<Ledger>(url).pipe(map(res => new Ledger().deserialize(res)));
	}

	//
	// Get me
	//
	get(page: number, type: string, search: string, category: Category, labels: Label[], year: number): Observable<LedgerResponse> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/ledger?page=${page}&type=${type}&search=${search}`;

		// Do we have a category?
		if (category) {
			url = url + `&category_id=${category.Id}`
		}

		// Do we have a labels?
		if (labels.length > 0) {
			let ll = [];

			for (let i = 0; i < labels.length; i++) {
				ll.push(labels[i].Id);
			}

			url = url + `&label_ids=${ll.join(",")}`
		}

		// Do we have a year?
		if (year) {
			url = url + `&year=${year}`
		}

		return this.http.get<Ledger[]>(url, { observe: 'response' }).pipe(map((res) => {
			// Setup data
			let data: Ledger[] = [];
			let lastPage = false;

			// Serialize the response.
			for (let i = 0; i < res.body.length; i++) {
				data.push(new Ledger().deserialize(res.body[i]));
			}

			// Build last page
			if (res.headers.get('X-Last-Page') == "true") {
				lastPage = true;
			}

			// Return happy.
			return new LedgerResponse(lastPage, Number(res.headers.get('X-Offset')), Number(res.headers.get('X-Limit')), Number(res.headers.get('X-No-Limit-Count')), data);
		}));
	}

	//
	// Get getLedgerSummary
	//
	getLedgerSummary(type: string): Observable<LedgerSummaryResponse> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/ledger-summary?type=${type}`;

		return this.http.get<LedgerSummaryResponse[]>(url, { observe: 'response' }).pipe(map((res) => {
			// Response object
			let ls = new LedgerSummaryResponse([], [], []);

			// Add in years
			for (let i = 0; i < res.body["years"].length; i++) {
				let r = res.body["years"][i];
				ls.Years.push(new LedgerYearSummaryResult(r["year"], r["count"]));
			}

			// Add in labels
			for (let i = 0; i < res.body["labels"].length; i++) {
				let r = res.body["labels"][i];
				ls.Labels.push(new LedgerSummaryResult(r["id"], r["name"], r["count"]));
			}

			// Add in Categories
			for (let i = 0; i < res.body["categories"].length; i++) {
				let r = res.body["categories"][i];
				ls.Categories.push(new LedgerSummaryResult(r["id"], r["name"], r["count"]));
			}

			return ls;
		}));
	}

	//
	// Get getLedgerPnlSummary
	//
	getLedgerPnlSummary(type: string, search: string, category: Category, labels: Label[], year: number): Observable<LedgerPnlSummary> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/ledger-pl-summary?type=${type}&search=${search}`;

		// Do we have a category?
		if (category) {
			url = url + `&category_id=${category.Id}`
		}

		// Do we have a labels?
		if (labels.length > 0) {
			let ll = [];

			for (let i = 0; i < labels.length; i++) {
				ll.push(labels[i].Id);
			}

			url = url + `&label_ids=${ll.join(",")} `
		}

		// Do we have a year?
		if (year) {
			url = url + `&year=${year}`
		}

		return this.http.get<LedgerPnlSummary>(url)
			.pipe(map(res => { return new LedgerPnlSummary(res["income"], res["expense"], res["profit"]) }));
	}
}

//
// Ledger Response
//
export class LedgerResponse {
	constructor(
		public LastPage: boolean,
		public Offset: number,
		public Limit: number,
		public NoLimitCount: number,
		public Data: Ledger[]
	) { }
}

//
// Ledger Summary Response
//
export class LedgerSummaryResponse {
	constructor(
		public Years: LedgerYearSummaryResult[],
		public Labels: LedgerSummaryResult[],
		public Categories: LedgerSummaryResult[]
	) { }
}

//
// LedgerSummaryResult
//
export class LedgerSummaryResult {
	constructor(
		public Id: number,
		public Name: string,
		public Count: number
	) { }
}

//
// LedgerYearSummaryResult
//
export class LedgerYearSummaryResult {
	constructor(
		public Year: number,
		public Count: number
	) { }
}

//
// LedgerPnlSummary
//
export class LedgerPnlSummary {
	constructor(
		public Income: number,
		public Expense: number,
		public Profit: number
	) { }
}

/* End File */
