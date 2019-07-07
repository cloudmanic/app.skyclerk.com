//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { Activity } from '../models/activity.model';

@Injectable({
	providedIn: 'root'
})

export class ActivityService {
	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// Get Activity - with a limit.
	//
	getWithLimit(page: number, limit: number): Observable<Activity[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/activities?page=${page}&limit=${limit}`;
		return this.http.get<Activity[]>(url)
			.pipe(map(res => res.map(res => new Activity().deserialize(res))));
	}

	//
	// Get Activity
	//
	get(page: number, limit: number, order: string, sort: string): Observable<ActivityResponse> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/activities?page=${page}&order=${order}&sort=${sort}&limit=${limit}`;

		return this.http.get<ActivityResponse[]>(url, { observe: 'response' }).pipe(map((res) => {
			// Setup data
			let data: Activity[] = [];
			let lastPage = false;

			// Serialize the response.
			for (let i = 0; i < res.body.length; i++) {
				data.push(new Activity().deserialize(res.body[i]));
			}

			// Build last page
			if (res.headers.get('X-Last-Page') == "true") {
				lastPage = true;
			}

			// Return happy.
			return new ActivityResponse(lastPage, Number(res.headers.get('X-Offset')), Number(res.headers.get('X-Limit')), Number(res.headers.get('X-No-Limit-Count')), data);
		}));
	}

	//
	// Get Activity by Ledger ID
	//
	getByLedgerId(ledgerId: number): Observable<Activity[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/activities?ledger_id=${ledgerId}`;
		return this.http.get<Activity[]>(url)
			.pipe(map(res => res.map(res => new Activity().deserialize(res))));
	}
}

//
// Activity Response
//
export class ActivityResponse {
	constructor(
		public LastPage: boolean,
		public Offset: number,
		public Limit: number,
		public NoLimitCount: number,
		public Data: Activity[]
	) { }
}


/* End File */
