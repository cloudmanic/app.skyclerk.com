//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
	providedIn: 'root'
})
export class SessionGuard implements CanActivate {
	//
	// Construct - Inject router so we can redirect later.
	//
	constructor(private router: Router) { }

	//
	// Is the user allowed to use this page?
	//
	canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
		// Make sure we have an access token.
		if (localStorage.getItem('access_token') && localStorage.getItem('user_id')) {
			return true;
		}

		// Not logged in so redirect to login page with the return url
		this.router.navigate(['/login'], { queryParams: { redirect_to: state.url } });
		return false;
	}
}

/* End File */
