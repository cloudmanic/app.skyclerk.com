//
// Date: 2019-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';

@Injectable({
	providedIn: 'root'
})

export class ReportService {
	// Used when account change happens.
	accountChange = new EventEmitter<number>();

	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// Get PnlCurrentYear
	//
	getPnlCurrentYear(): Observable<PnlCurrentYear> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<PnlCurrentYear>(`${environment.app_server}/api/v3/${accountId}/reports/pnl-current-year`)
			.pipe(map(res => { return { Year: res["year"], Value: res["value"] } }));
	}
}

export interface PnlCurrentYear {
	Year: number,
	Value: number
}

/* End File */
