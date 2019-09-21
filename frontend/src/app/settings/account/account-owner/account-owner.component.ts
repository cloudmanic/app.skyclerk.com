//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { UserService } from 'src/app/services/user.service';
import { User } from 'src/app/models/user.model';
import { MeService } from 'src/app/services/me.service';
import { Me } from 'src/app/models/me.model';
import { AccountService } from 'src/app/services/account.service';
import { Account } from 'src/app/models/account.model';
import { Router } from '@angular/router';

@Component({
	selector: 'app-settings-account-account-owner',
	templateUrl: './account-owner.component.html'
})

export class AccountOwnerComponent {
	me: Me = new Me();
	editMode: boolean = false;
	users: User[] = [];
	account: Account = new Account();
	userName: string = "";

	//
	// Constructor
	//
	constructor(public meService: MeService, public userService: UserService, public accountService: AccountService, public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Get Me.
		this.meService.get().subscribe(res => {
			this.me = res;
		});

		// Get the active account.
		this.accountService.getAccount().subscribe(res => {
			this.account = res;
			this.setDisplayName();
		});

		// Get the users.
		this.loadUsers();
	}

	//
	// Load users
	//
	loadUsers() {
		this.userService.get().subscribe((res) => {
			this.users = res;
			this.setDisplayName();
		});
	}

	//
	// Set display name
	//
	setDisplayName() {
		for (let i = 0; i < this.users.length; i++) {
			if (this.account.OwnerId == this.users[i].Id) {
				this.userName = this.users[i].FirstName + " " + this.users[i].LastName;
				return;
			}
		}
	}

	//
	// Update the user title in the selector.
	//
	onChange() {
		console.log(this.account.OwnerId);
		this.setDisplayName();
	}

	//
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}

	//
	// On save.
	//
	onSave() {
		let c = confirm("Are you sure you want to change the owner? You will not be able to access account or billing settings after this.");

		if (!c) {
			return;
		}

		// Close widget.
		this.editModeToggle();

		// Send updated account to server.
		this.accountService.update(this.account).subscribe(() => {
			// Reset the current account.
			this.accountService.setActiveAccount();
			this.router.navigate(['/']);
		});
	}
}

/* End File */
