import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FileDropModule } from 'ngx-file-drop';

// Services
import { MeService } from './services/me.service';
import { AuthService } from './services/auth.service';
import { TokenInterceptor } from './services/token.interceptor';
import { ReportService } from './services/report.service';
import { ContactService } from './services/contact.service';

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
import { AddComponent as LedgerAddComponent } from './ledger/add/add.component';
import { AddEditPopOutComponent as ContactsAddEditPopOutComponent } from './contacts/add-edit-pop-out/add-edit-pop-out.component';
import { AssignFieldComponent as ContactsAssignFieldComponent } from './contacts/assign-field/assign-field.component';
import { BaseComponent as ContactsBase } from './contacts/base/base.component';
import { PopOutComponent as CategoriesPopOutComponent } from './categories/pop-out/pop-out.component';
import { FieldComponent as CategoriesFieldComponent } from './categories/field/field.component';
import { FieldComponent as LabelsFieldComponent } from './labels/field/field.component';
import { PopOutComponent as LabelsPopOutComponent } from './labels/pop-out/pop-out.component';
import { UploadComponent as FilesUploadComponent } from './files/upload/upload.component';
import { ViewComponent as LedgerViewComponent } from './ledger/view/view.component';
import { EditComponent as LedgerEditComponent } from './ledger/edit/edit.component';

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
		SettingsAccountComponent,
		LedgerAddComponent,
		ContactsAssignFieldComponent,
		ContactsAddEditPopOutComponent,
		ContactsBase,
		CategoriesPopOutComponent,
		CategoriesFieldComponent,
		LabelsFieldComponent,
		LabelsPopOutComponent,
		FilesUploadComponent,
		LedgerViewComponent,
		LedgerEditComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		AppRoutingModule,
		HttpClientModule,
		BrowserAnimationsModule,
		FileDropModule
	],
	providers: [
		MeService,
		AuthService,
		ReportService,
		ContactService,
		{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true },
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
