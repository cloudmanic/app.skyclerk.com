//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

// Guards
import { SessionGuard } from './auth/guards/session.guard';

// Components
import { LoginComponent } from './auth/login/login.component';
import { AuthComponent as LayoutAuthComponent } from './layouts/auth/auth.component';
import { AppComponent as LayoutAppComponent } from './layouts/app/app.component';
import { LandingComponent as DashboardLandingComponent } from './dashboard/landing/landing.component';

const routes: Routes = [

	// Core App - with main css div
	{
		path: '', component: LayoutAppComponent, children: [

			// dashboard
			{ path: '', component: DashboardLandingComponent, canActivate: [SessionGuard] },

		]
	},

	// login
	{
		path: '', component: LayoutAuthComponent, children: [
			{ path: 'login', component: LoginComponent }
		]
	},


	// Otherwise redirect to home
	{ path: '**', redirectTo: '' }

];

@NgModule({
	imports: [RouterModule.forRoot(routes)],
	exports: [RouterModule]
})
export class AppRoutingModule { }

/* End File */
