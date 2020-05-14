//
// Date: 2020-05-13
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';

// providers
import { TokenInterceptor } from './services/token.interceptor';

// declarations
import { AppComponent } from './app.component';
import { AccountsComponent } from './accounts/accounts.component';
import { CoreComponent } from './layout/core/core.component';
import { ViewComponent } from './snapclerk/view/view.component';

@NgModule({
	declarations: [
		AppComponent,
		AccountsComponent,
		CoreComponent,
		ViewComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		AppRoutingModule,
		HttpClientModule
	],
	providers: [
		{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true },
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
