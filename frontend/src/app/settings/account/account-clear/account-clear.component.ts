//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';
import { Account } from 'src/app/models/account.model';

@Component({
	selector: 'app-settings-account-account-clear',
	templateUrl: './account-clear.component.html'
})

export class AccountClearComponent {
	confirm: boolean = false;
	editMode: boolean = false;
	account: Account = new Account();
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public accountService: AccountService) { }

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

		let c = confirm("Are you sure you want to clear this account? ALL DATA WILL BE LOST FOREVER.");

		if (!c) {
			return
		}

		// Clear the account.
		this.editModeToggle();
		this.accountService.clear().subscribe(() => {
			// Tell the rest of the app to update the account.
			this.accountService.accountChange.emit(this.account.Id);
		});
	}

	//
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}
}

/* End File */
