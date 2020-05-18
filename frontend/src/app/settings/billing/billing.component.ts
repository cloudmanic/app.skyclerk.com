//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';
import { Billing } from 'src/app/models/billing.model';

const pageTitle: string = environment.title_prefix + "Settings Billing";

@Component({
	selector: 'app-settings-billing',
	templateUrl: './billing.component.html'
})

export class BillingComponent implements OnInit {
	billing: Billing = new Billing();
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Construct.
	//
	constructor(public titleService: Title, public accountService: AccountService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Refresh account.
		this.refreshBilling();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshBilling();
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
	refreshBilling() {
		// Get the billing.
		this.accountService.getBilling().subscribe(res => {
			this.billing = res;

			console.log(this.billing);
		});
	}

	//
	// Returns the number of says until expire
	//
	daysToExpire(): number {
		let today = moment();
		let expire = moment(this.billing.TrialExpire);
		let days = expire.diff(today, 'days');
		return days;
	}
}

/* End File */
