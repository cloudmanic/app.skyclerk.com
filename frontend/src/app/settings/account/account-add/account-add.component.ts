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
import { Me } from 'src/app/models/me.model';

@Component({
	selector: 'app-settings-account-account-add',
	templateUrl: './account-add.component.html'
})

export class AccountAddComponent implements OnInit {
	me: Me = new Me();
	name: string = "";
	errMsg: string = "";
	editMode: boolean = false;
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
		// Get data for page.
		this.getLoggedInUser();
		this.refreshAccount();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.getLoggedInUser();
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
	// Load up the logged in user.
	//
	getLoggedInUser() {
		this.meService.get().subscribe(res => {
			this.me = res;
		});
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
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}

	//
	// Sumbit the new account request.
	//
	submit() {
		let name = this.name.trim();

		// Poor man's validation
		if (name.length == 0) {
			this.errMsg = "Please enter an account name.";
			return;
		}

		// Make sure we are not already using this name.
		for (let i = 0; i < this.me.Accounts.length; i++) {
			if (this.me.Accounts[i].Name == name) {
				this.errMsg = "Account name is already in use.";
				return;
			}
		}

		// Send request to create a new account.
		this.accountService.create(name).subscribe(res => {
			this.doSelectAccount(res);
			this.editModeToggle();
			this.router.navigate(['/']);
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
}

/* End File */
