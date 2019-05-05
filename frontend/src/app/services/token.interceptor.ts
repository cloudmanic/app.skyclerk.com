//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { HttpRequest, HttpHandler, HttpEvent, HttpInterceptor } from '@angular/common/http';
import { Observable } from 'rxjs';

export class TokenInterceptor implements HttpInterceptor {
	//
	// Construct
	//
	constructor() { }

	//
	// Run intercept
	//
	intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
		// Set the access token
		let access_token = localStorage.getItem('access_token');

		// Add in the headers for the API.
		request = request.clone({ setHeaders: { Authorization: `Bearer ${access_token}` } });

		// Let the request continue.
		return next.handle(request);
	}
}

/* End File */
