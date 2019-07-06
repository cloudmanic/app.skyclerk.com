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
	// Get Activity
	//
	get(page: number, limit: number): Observable<Activity[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/activities?page=${page}&limit=${limit}`;
		return this.http.get<Activity[]>(url)
			.pipe(map(res => res.map(res => new Activity().deserialize(res))));
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


/* End File */
