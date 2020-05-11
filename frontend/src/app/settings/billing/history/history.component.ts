//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { Subject } from 'rxjs';
import { AccountService } from 'src/app/services/account.service';
import { BillingInvoice } from 'src/app/models/billing-invoice.model';

@Component({
	selector: 'app-settings-billing-history',
	templateUrl: './history.component.html'
})

export class HistoryComponent {
	editMode: boolean = false;
	history: BillingInvoice[] = [];
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
		this.refreshBillingHistory();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshBillingHistory();
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
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}

	//
	// Refresh account.
	//
	refreshBillingHistory() {
		// Get the billing.
		this.accountService.getBillingHistory().subscribe(res => {
			this.history = res;
		});
	}
}

/* End File */
