//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { MeService } from 'src/app/services/me.service';
import { Me } from 'src/app/models/me.model';
import { Router } from '@angular/router';
import { Account } from 'src/app/models/account.model';

@Component({
	selector: 'app-layouts-app',
	templateUrl: './app.component.html'
})
export class AppComponent implements OnInit {
	me: Me = new Me();
	account: Account = new Account();
	accountToggle: boolean = false;

	//
	// Constructor.
	//
	constructor(public meService: MeService, public router: Router) { }

	//
	// NgOnInit
	//
	ngOnInit() {
		// Load the logged in user.
		this.getLoggedInUser();
	}

	//
	// Here we change the active account.
	//
	doSelectAccount(account: Account) {
		localStorage.setItem('account_id', account.Id.toString());
		this.account = account;
		this.accountToggle = false;
	}

	//
	// Do account select toggle.
	//
	doAccountToggle() {
		if (this.accountToggle) {
			this.accountToggle = false;
		} else {
			this.accountToggle = true;
		}
	}

	//
	// Load up the logged in user.
	//
	getLoggedInUser() {
		this.meService.get().subscribe(res => {
			this.me = res;

			// Get the account id
			let accountId = Number(localStorage.getItem("account_id"));

			// Set our default account.
			for (let i = 0; i < this.me.Accounts.length; i++) {
				if (this.me.Accounts[i].Id == accountId) {
					this.account = this.me.Accounts[i];
					break;
				}
			}
		});
	}

	//
	// Log user out.
	//
	doLogOut() {
		this.meService.logout();
		this.router.navigate(['/login']);
	}
}

/* End File */
