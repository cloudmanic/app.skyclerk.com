import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

// Services
import { MeService } from './services/me.service';
import { AuthService } from './services/auth.service';
import { TokenInterceptor } from './services/token.interceptor';

// Components
import { AppComponent } from './app.component';
import { LoginComponent } from './auth/login/login.component';
import { AuthComponent as LayoutAuthComponent } from './layouts/auth/auth.component';
import { AppComponent as LayoutAppComponent } from './layouts/app/app.component';
import { LandingComponent as DashboardLandingComponent } from './dashboard/landing/landing.component';
import { SidebarComponent } from './layouts/sidebar/sidebar.component';
import { LandingComponent as LedgerLandingComponent } from './ledger/landing/landing.component';
import { LandingComponent as SnapclerkLandingComponent } from './snapclerk/landing/landing.component';
import { UsersComponent as SettingsUsersComponent } from './settings/users/users.component';
import { SubNavComponent as SettingsSubNavComponent } from './settings/sub-nav/sub-nav.component';
import { AccountComponent as SettingsAccountComponent } from './settings/account/account.component';


@NgModule({
	declarations: [
		AppComponent,
		LoginComponent,
		LayoutAuthComponent,
		LayoutAppComponent,
		DashboardLandingComponent,
		SidebarComponent,
		LedgerLandingComponent,
		SnapclerkLandingComponent,
		SettingsUsersComponent,
		SettingsSubNavComponent,
		SettingsAccountComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		AppRoutingModule,
		HttpClientModule
	],
	providers: [
		MeService,
		AuthService,
		{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true },
	],
	bootstrap: [AppComponent]
})
export class AppModule { }