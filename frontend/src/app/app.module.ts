import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule } from '@angular/common/http';

// Services
import { AuthService } from './services/auth.service';

// Components
import { AppComponent } from './app.component';
import { LoginComponent } from './auth/login/login.component';
import { AuthComponent as LayoutAuthComponent } from './layouts/auth/auth.component';
import { AppComponent as LayoutAppComponent } from './layouts/app/app.component';
import { LandingComponent as DashboardLandingComponent } from './dashboard/landing/landing.component';

@NgModule({
	declarations: [
		AppComponent,
		LoginComponent,
		LayoutAuthComponent,
		LayoutAppComponent,
		DashboardLandingComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		AppRoutingModule,
		HttpClientModule
	],
	providers: [
		AuthService
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
