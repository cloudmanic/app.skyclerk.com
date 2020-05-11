//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { Billing } from 'src/app/models/billing.model';
import { Subject } from 'rxjs';
import { AccountService } from 'src/app/services/account.service';

@Component({
	selector: 'app-settings-billing-next-payment',
	templateUrl: './next-payment.component.html'
})

export class NextPaymentComponent {
	editMode: boolean = false;
	billing: Billing = new Billing();
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public accountService: AccountService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
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
			console.log(this.billing.CurrentPeriodEnd);
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
