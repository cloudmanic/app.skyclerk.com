//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AccountService } from 'src/app/services/account.service';
import { Router } from '@angular/router';
import { MeService } from 'src/app/services/me.service';

@Component({
	selector: 'app-wall',
	templateUrl: './wall.component.html'
})

export class WallComponent implements OnInit {
	//
	// Constructor
	//
	constructor(public accountService: AccountService, public router: Router, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Sumbit the close account.
	//
	closeAccount() {
		// Confirm
		let c = confirm("Are you sure you want to delete this account? ALL DATA WILL BE LOST FOREVER.");

		if (!c) {
			return
		}

		// Clear the account.
		this.accountService.delete().subscribe((_res) => {
			// Tell user TODO(spicer): Make this better in terms of UI
			alert("Your account was successfully deleted.");

			// TODO(spicer): Do we delete all the accounts? For now we assume people just have one account.
			// If they log back in they will just default to another account.

			// Log user out.
			this.meService.logout();
			this.router.navigate(['/login']);
		});
	}

}

/* End File */
