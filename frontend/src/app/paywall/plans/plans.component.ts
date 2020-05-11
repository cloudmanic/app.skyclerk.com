//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';
import { MeService } from 'src/app/services/me.service';

@Component({
	selector: 'app-plans',
	templateUrl: './plans.component.html'
})
export class PlansComponent implements OnInit {
	back: string = "";
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor.
	//
	constructor(public route: ActivatedRoute, public accountService: AccountService, public router: Router, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set plan on load.
		this.back = this.route.snapshot.queryParamMap.get("back");

		// When plan changes
		this.route.queryParamMap.takeUntil(this.destory).subscribe(queryParams => {
			this.back = queryParams.get("back");
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
