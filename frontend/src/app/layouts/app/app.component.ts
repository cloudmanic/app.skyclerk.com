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
import { ReportService, PnlCurrentYear } from 'src/app/services/report.service';
import { AccountService } from 'src/app/services/account.service';

@Component({
	selector: 'app-layouts-app',
	templateUrl: './app.component.html'
})
export class AppComponent implements OnInit {
	me: Me = new Me();
	pnl: PnlCurrentYear = { Year: 0, Value: 0 };
	account: Account = new Account();
	accountToggle: boolean = false;

	//
	// Constructor.
	//
	constructor(public meService: MeService, public accountService: AccountService, public router: Router, public reportService: ReportService) { }

	//
	// NgOnInit
	//
	ngOnInit() {
		// Load data for page
		this.loadPageData();

		// Listen for account changes.
		this.meService.accountChange.subscribe(() => {
			this.loadPageData();
		});
	}

	//
	// loadPageData
	//
	loadPageData() {
		// Load the logged in user.
		this.getLoggedInUser();

		//  Load PnL
		this.getPnL();
	}

	//
	// Get the Year PNL
	//
	getPnL() {
		this.reportService.getPnlCurrentYear().subscribe(res => {
			this.pnl = res;
		});
	}

	//
	// Here we change the active account.
	//
	doSelectAccount(account: Account) {
		localStorage.setItem('account_id', account.Id.toString());
		this.account = account;
		this.accountToggle = false;

		// Reset the current account.
		this.accountService.setActiveAccount();

		// Tell the rest of the app the account switched.
		this.meService.accountChange.emit(account.Id);
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
