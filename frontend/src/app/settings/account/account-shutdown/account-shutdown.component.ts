//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AccountService } from 'src/app/services/account.service';
import { MeService } from 'src/app/services/me.service';
import { Router } from '@angular/router';
import { Subject } from 'rxjs';
import { Account } from 'src/app/models/account.model';

@Component({
	selector: 'app-settings-account-account-shutdown',
	templateUrl: './account-shutdown.component.html'
})

export class AccountShutdownComponent implements OnInit {
	editMode: boolean = false;
	confirm: boolean = false;
	confirmRefund: boolean = false;
	account: Account = new Account();
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public accountService: AccountService, public meService: MeService, public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshAccount();
		});
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Refresh account.
	//
	refreshAccount() {
		// Get the active account.
		this.accountService.getAccount().subscribe(res => {
			this.account = res;
		});
	}

	//
	// Sumbit the clear account.
	//
	submit() {
		if (!this.confirm) {
			alert("Please check the confirmation box.");
			return;
		}

		if (!this.confirmRefund) {
			alert("Please check the refund policy box.");
			return;
		}

		let c = confirm("Are you sure you want to delete this account? ALL DATA WILL BE LOST FOREVER.");

		if (!c) {
			return
		}

		// Clear the account.
		this.accountService.delete().subscribe((res) => {
			// Tell user TODO(spicer): Make this better in terms of UI
			alert("Your account was successfully deleted.");

			// Do we have other accounts?
			if (res.length > 0) {
				this.doSelectAccount(res[0]);
				this.editModeToggle();
				this.router.navigate(['/']);
				return;
			}

			// Log user out.
			this.meService.logout();
			this.router.navigate(['/login']);
		});
	}

	//
	// Here we change the active account.
	//
	doSelectAccount(account: Account) {
		localStorage.setItem('account_id', account.Id.toString());
		this.account = account;

		// Reset the current account.
		this.accountService.setActiveAccount();

		// Tell the rest of the app the account switched.
		this.accountService.accountChange.emit(account.Id);
	}

	//
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}
}

/* End File */
