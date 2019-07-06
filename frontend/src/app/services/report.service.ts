//
// Date: 2019-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
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

	//
	// Get Pnl
	//
	getPnl(start: Date, end: Date, group: string, sort: string): Observable<Pnl[]> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<Pnl[]>(`${environment.app_server}/api/v3/${accountId}/reports/pnl?sort=${sort}&start=${moment(start).format('YYYY-MM-DD')}&end=${moment(end).format('YYYY-MM-DD')}&group=${group}`)
			.pipe(map(res => {
				let rt = [];

				for (let i = 0; i < res.length; i++) {
					rt.push({ Date: moment(res[i]["date"]).toDate(), Profit: res[i]["profit"], Expense: res[i]["expense"], Income: res[i]["income"] });
				}

				return rt;
			}));
	}

	//
	// Get profit loss by cateogry
	//
	getProfitLossByCategory(start: Date, end: Date, sort: string): Observable<PnlNameAmount[]> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<PnlNameAmount[]>(`${environment.app_server}/api/v3/${accountId}/reports/pnl-category?sort=${sort}&start=${moment(start).format('YYYY-MM-DD')}&end=${moment(end).format('YYYY-MM-DD')}`)
			.pipe(map(res => {
				let rt = [];

				for (let i = 0; i < res.length; i++) {
					rt.push({ Name: res[i]["name"], Amount: res[i]["amount"] });
				}

				return rt;
			}));
	}

	//
	// Get profit loss by label
	//
	getProfitLossByLabel(start: Date, end: Date, sort: string): Observable<PnlNameAmount[]> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<PnlNameAmount[]>(`${environment.app_server}/api/v3/${accountId}/reports/pnl-label?sort=${sort}&start=${moment(start).format('YYYY-MM-DD')}&end=${moment(end).format('YYYY-MM-DD')}`)
			.pipe(map(res => {
				let rt = [];

				for (let i = 0; i < res.length; i++) {
					rt.push({ Name: res[i]["name"], Amount: res[i]["amount"] });
				}

				return rt;
			}));
	}

	//
	// Get expense by contact
	//
	getExpenseByContact(start: Date, end: Date, sort: string): Observable<PnlNameAmount[]> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<PnlNameAmount[]>(`${environment.app_server}/api/v3/${accountId}/reports/expenses-by-contact?sort=${sort}&start=${moment(start).format('YYYY-MM-DD')}&end=${moment(end).format('YYYY-MM-DD')}`)
			.pipe(map(res => {
				let rt = [];

				for (let i = 0; i < res.length; i++) {
					rt.push({ Name: res[i]["name"], Amount: res[i]["amount"] });
				}

				return rt;
			}));
	}

	//
	// Get income by contact
	//
	getIncomeByContact(start: Date, end: Date, sort: string): Observable<PnlNameAmount[]> {
		let accountId = localStorage.getItem('account_id');
		return this.http.get<PnlNameAmount[]>(`${environment.app_server}/api/v3/${accountId}/reports/income-by-contact?sort=${sort}&start=${moment(start).format('YYYY-MM-DD')}&end=${moment(end).format('YYYY-MM-DD')}`)
			.pipe(map(res => {
				let rt = [];

				for (let i = 0; i < res.length; i++) {
					rt.push({ Name: res[i]["name"], Amount: res[i]["amount"] });
				}

				return rt;
			}));
	}
}

export interface PnlCurrentYear {
	Year: number,
	Value: number
}

export interface Pnl {
	Date: Date,
	Profit: number,
	Expense: number,
	Income: number
}

export interface PnlNameAmount {
	Name: string,
	Amount: number
}
/* End File */
