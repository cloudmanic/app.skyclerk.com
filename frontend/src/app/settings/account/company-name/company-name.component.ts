//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { AccountService } from 'src/app/services/account.service';
import { Account } from 'src/app/models/account.model';
import { Me } from 'src/app/models/me.model';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-settings-account-company-name',
	templateUrl: './company-name.component.html'
})

export class CompanyNameComponent {
	me: Me = new Me();
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
		this.refreshAccount();

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
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
		this.refreshAccount();
	}

	//
	// Save Name
	//
	save() {
		// Make sure the account name is not empty.
		if (this.account.Name.length == 0) {
			alert("Account name can not be empty.");
			return;
		}

		// Send updated account to server.
		this.accountService.update(this.account).subscribe(
			// Success
			() => {
				// Reset the current account.
				this.editModeToggle();
				this.accountService.setActiveAccount();

				// Tell the rest of the app to update the account.
				this.accountService.accountChange.emit(this.account.Id);
			},

			// error
			(err) => {
				alert(err.error.errors.name);
			}
		);
	}
}

/* End File */
