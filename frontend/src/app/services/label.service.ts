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
import { TrackService } from './track.service';

@Injectable({
	providedIn: 'root'
})

export class LabelService {
	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) { }

	//
	// Get labels
	//
	get(): Observable<Label[]> {
		let accountId = localStorage.getItem('account_id');
		let url = environment.app_server + '/api/v3/' + accountId + '/labels';
		return this.http.get<Label[]>(url)
			.pipe(map(res => res.map(res => new Label().deserialize(res))));
	}

	//
	// Create a new label
	//
	create(lb: Label): Observable<Label> {
		let accountId = localStorage.getItem('account_id');
		lb.AccountId = Number(accountId);

		return this.http.post<number>(`${environment.app_server}/api/v3/${accountId}/labels`, new Label().serialize(lb))
			.pipe(map(res => {
				let lb = new Label().deserialize(res);

				// Track event.
				this.trackService.event('label-create', { app: "web", "accountId": accountId });

				return lb;
			}));
	}

	//
	// Update a label
	//
	update(label: Label): Observable<Label> {
		let accountId = localStorage.getItem('account_id');

		let put = {
			name: label.Name,
			account_id: Number(accountId)
		}

		return this.http.put<Label>(`${environment.app_server}/api/v3/${accountId}/labels/${label.Id}`, put)
			.pipe(map(res => {
				let lb = new Label().deserialize(res);

				// Track event.
				this.trackService.event('label-update', { app: "web", "accountId": accountId });

				return lb;
			}));
	}

	//
	// Delete a label
	//
	delete(label: Label): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');
		label.AccountId = Number(accountId);

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/labels/${label.Id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('label-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}
}


/* End File */
