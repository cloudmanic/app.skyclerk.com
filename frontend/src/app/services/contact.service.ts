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
import { Contact } from '../models/contact.model';
import { TrackService } from './track.service';

@Injectable({
	providedIn: 'root'
})

export class ContactService {
	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) { }

	//
	// Get contacts
	//
	get(limit: number, search: string): Observable<Contact[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/contacts?limit=${limit}`;

		// Are we searching?
		if (search.length > 0) {
			url = url + "&search=" + search
		}

		return this.http.get<Contact[]>(url)
			.pipe(map(res => res.map(res => new Contact().deserialize(res))));
	}

	//
	// Create a new contact
	//
	create(contact: Contact): Observable<Contact> {
		let accountId = localStorage.getItem('account_id');
		contact.AccountId = Number(accountId);

		return this.http.post<number>(`${environment.app_server}/api/v3/${accountId}/contacts`, new Contact().serialize(contact))
			.pipe(map(res => {
				let con = new Contact().deserialize(res);

				// Track event.
				this.trackService.event('contact-create', { app: "web", "accountId": accountId });

				return con;
			}));
	}

	//
	// Update a contact
	//
	update(contact: Contact): Observable<Contact> {
		let accountId = localStorage.getItem('account_id');
		contact.AccountId = Number(accountId);

		return this.http.put<Contact>(`${environment.app_server}/api/v3/${accountId}/contacts/${contact.Id}`, new Contact().serialize(contact))
			.pipe(map(res => {
				let con = new Contact().deserialize(res);

				// Track event.
				this.trackService.event('contact-update', { app: "web", "accountId": accountId });

				return con;
			}));
	}

	//
	// Delete a contact
	//
	delete(contact: Contact): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');
		contact.AccountId = Number(accountId);

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/contacts/${contact.Id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('contact-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}
}


/* End File */
