//
// Date: 2020-05-13
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CoreComponent } from './layout/core/core.component';
import { AccountsComponent } from './accounts/accounts.component';
import { ViewComponent } from './snapclerk/view/view.component';
import { SessionGuard } from './auth/guards/session.guard';

const routes: Routes = [

	// Centcom
	{
		path: '', component: CoreComponent, data: { tab: 'ff' }, children: [

			// users
			{ path: 'accounts', component: AccountsComponent, data: { tab: 'accounts', pageHeading: 'Acccounts' }, canActivate: [SessionGuard] },

			// snapclerk
			{ path: 'snapclerk', component: ViewComponent, data: { tab: 'snapclerk', pageHeading: 'Snap!Clerk' }, canActivate: [SessionGuard] },

		]
	},

	// Otherwise redirect to home
	{ path: '**', redirectTo: 'accounts' }

];

@NgModule({
	imports: [RouterModule.forRoot(routes)],
	exports: [RouterModule]
})
export class AppRoutingModule { }

/* End File */
