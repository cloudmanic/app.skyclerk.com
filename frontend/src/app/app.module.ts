import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FileDropModule } from 'ngx-file-drop';
import { AgmCoreModule } from '@agm/core';

// Pipes
import { CallbackPipe } from './pipes/callback.pipe';
import { SafeHtmlPipe } from './pipes/safe-html.pipe';
import { CurrencyFormatPipe } from './pipes/currency-format.pipe';

// Services
import { MeService } from './services/me.service';
import { UserService } from './services/user.service';
import { AuthService } from './services/auth.service';
import { TokenInterceptor } from './services/token.interceptor';
import { ReportService } from './services/report.service';
import { ContactService } from './services/contact.service';
import { ActivityService } from './services/activity.service';
import { AccountService } from './services/account.service';

// Components
import { AppComponent } from './app.component';
import { LoginComponent } from './auth/login/login.component';
import { AuthComponent as LayoutAuthComponent } from './layouts/auth/auth.component';
import { AppComponent as LayoutAppComponent } from './layouts/app/app.component';
import { SummaryComponent as DashboardSummaryComponent } from './dashboard/summary/summary.component';
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
import { ActivityComponent } from './activity/activity.component';
import { GraphsComponent } from './dashboard/graphs/graphs.component';
import { SubNavComponent } from './dashboard/sub-nav/sub-nav.component';
import { ReportsComponent } from './dashboard/reports/reports.component';
import { BillingComponent } from './settings/billing/billing.component';
import { ContactsComponent } from './settings/contacts/contacts.component';
import { CategoriesLabelsComponent } from './settings/categories-labels/categories-labels.component';
import { CompanyNameComponent } from './settings/account/company-name/company-name.component';
import { CurrencyComponent } from './settings/account/currency/currency.component';
import { AccountOwnerComponent } from './settings/account/account-owner/account-owner.component';
import { AccountAddComponent } from './settings/account/account-add/account-add.component';
import { AccountShutdownComponent } from './settings/account/account-shutdown/account-shutdown.component';
import { AccountClearComponent } from './settings/account/account-clear/account-clear.component';
import { AccountPlanComponent } from './settings/billing/account-plan/account-plan.component';
import { PaymentMethodComponent } from './settings/billing/payment-method/payment-method.component';
import { NextPaymentComponent } from './settings/billing/next-payment/next-payment.component';
import { HistoryComponent } from './settings/billing/history/history.component';
import { CategoriesComponent } from './settings/categories-labels/categories/categories.component';
import { LabelsComponent } from './settings/categories-labels/labels/labels.component';
import { AddComponent as SettingsUsersAdd } from './settings/users/add/add.component';
import { ForgotPasswordComponent } from './auth/forgot-password/forgot-password.component';

@NgModule({
	declarations: [
		CallbackPipe,
		SafeHtmlPipe,
		CurrencyFormatPipe,
		AppComponent,
		LoginComponent,
		LayoutAuthComponent,
		LayoutAppComponent,
		DashboardSummaryComponent,
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
		LedgerEditComponent,
		ActivityComponent,
		GraphsComponent,
		SubNavComponent,
		ReportsComponent,
		BillingComponent,
		ContactsComponent,
		CategoriesLabelsComponent,
		CompanyNameComponent,
		CurrencyComponent,
		AccountOwnerComponent,
		AccountAddComponent,
		AccountShutdownComponent,
		AccountClearComponent,
		AccountPlanComponent,
		PaymentMethodComponent,
		NextPaymentComponent,
		HistoryComponent,
		CategoriesComponent,
		LabelsComponent,
		SettingsUsersAdd,
		ForgotPasswordComponent
	],
	imports: [
		FormsModule,
		BrowserModule,
		AppRoutingModule,
		HttpClientModule,
		BrowserAnimationsModule,
		FileDropModule,
		AgmCoreModule.forRoot({ apiKey: 'AIzaSyCc8fAAyASKh3FzA0IXCjIKFl5oFF5i1zU' })
	],
	providers: [
		MeService,
		UserService,
		AuthService,
		ReportService,
		AccountService,
		ContactService,
		ActivityService,
		{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true },
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
