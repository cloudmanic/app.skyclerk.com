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
import { LandingComponent as LedgerLandingComponent } from './ledger/landing/landing.component';
import { LandingComponent as SnapclerkLandingComponent } from './snapclerk/landing/landing.component';
import { UsersComponent as SettingsUsersComponent } from './settings/users/users.component';
import { AccountComponent as SettingsAccountComponent } from './settings/account/account.component';
import { ViewComponent as LedgerViewComponent } from './ledger/view/view.component';
import { EditComponent as LedgerEditComponent } from './ledger/edit/edit.component';
import { ActivityComponent } from './activity/activity.component';
import { GraphsComponent } from './dashboard/graphs/graphs.component';
import { ReportsComponent } from './dashboard/reports/reports.component';

const routes: Routes = [

	// Core App - with main css div
	{
		path: '', component: LayoutAppComponent, children: [

			// dashboard
			{ path: '', component: DashboardLandingComponent, canActivate: [SessionGuard] },
			{ path: 'dashboard/graphs', component: GraphsComponent, canActivate: [SessionGuard] },
			{ path: 'dashboard/reports', component: ReportsComponent, canActivate: [SessionGuard] },

			// activity
			{ path: 'activity', component: ActivityComponent, canActivate: [SessionGuard] },

			// ledger
			{ path: 'ledger', component: LedgerLandingComponent, canActivate: [SessionGuard] },
			{ path: 'ledger/:id', component: LedgerViewComponent, canActivate: [SessionGuard] },
			{ path: 'ledger/:id/edit', component: LedgerEditComponent, canActivate: [SessionGuard] },

			// snapclerk
			{ path: 'snapclerk', component: SnapclerkLandingComponent, canActivate: [SessionGuard] },

			// settings
			{ path: 'settings/users', component: SettingsUsersComponent, canActivate: [SessionGuard] },
			{ path: 'settings/account', component: SettingsAccountComponent, canActivate: [SessionGuard] },
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
