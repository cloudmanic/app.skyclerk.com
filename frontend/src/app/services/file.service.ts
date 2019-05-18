//
// Date: 2019-05-17
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { File as FileModel } from '../models/file.model';

@Injectable({
	providedIn: 'root'
})

export class FileService {
	//
	// Constructor
	//
	constructor(private http: HttpClient) { }

	//
	// Upload a file. Returns the file ID
	//
	upload(file: File): Observable<any> {
		let accountId = localStorage.getItem('account_id');

		// Create upload form.
		const formData = new FormData()
		formData.append('file', file, file.name);

		// Upload file with progress
		return this.http.post<any>(`${environment.app_server}/api/v3/${accountId}/files`, formData, {
			reportProgress: true,
			observe: 'events'
		}).pipe(map((event) => {

			switch (event.type) {
				case HttpEventType.UploadProgress:
					const progress = Math.round(100 * event.loaded / event.total);
					return { status: 'progress', message: progress };

				case HttpEventType.Response:
					return new FileModel().deserialize(event.body);

				default:
					return `Unhandled event: ${event.type}`;
			}
		})

		);
	}
}

/* End File */
