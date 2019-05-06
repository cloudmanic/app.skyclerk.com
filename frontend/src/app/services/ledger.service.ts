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

/* End File */
