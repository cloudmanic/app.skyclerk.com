//
// Date: 2019-04-27
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { Label } from '../models/label.model';

@Injectable({
	providedIn: 'root'
})

export class LabelService {
	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// Get labels
	//
	get(): Observable<Label[]> {
		let accountId = localStorage.getItem('account_id');
		let url = environment.app_server + '/api/v3/' + accountId + '/labels';
		return this.http.get<Label[]>(url)
			.pipe(map(res => res.map(res => new Label().deserialize(res))));
	}
}


/* End File */
