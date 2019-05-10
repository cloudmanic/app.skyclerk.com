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

@Injectable({
	providedIn: 'root'
})

export class LedgerService {
	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// Get me
	//
	get(page: number, type: string, search: string): Observable<LedgerResponse> {
		let accountId = localStorage.getItem('account_id');
		let url = environment.app_server + '/api/v3/' + accountId + '/ledger?page=' + page + '&type=' + type + '&search=' + search;

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
		let url = environment.app_server + '/api/v3/' + accountId + '/ledger-summary?type=' + type;

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
	getLedgerPnlSummary(type: string, search: string): Observable<LedgerPnlSummary> {
		let accountId = localStorage.getItem('account_id');
		let url = environment.app_server + '/api/v3/' + accountId + '/ledger-pl-summary?type=' + type + '&search=' + search;

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
