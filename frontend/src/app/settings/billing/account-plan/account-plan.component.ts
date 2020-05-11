//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Subject } from 'rxjs';
import { Billing } from 'src/app/models/billing.model';
import { AccountService } from 'src/app/services/account.service';

@Component({
	selector: 'app-settings-billing-account-plan',
	templateUrl: './account-plan.component.html'
})

export class AccountPlanComponent implements OnInit {
	checked: boolean = false;
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
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}

	//
	// Refresh account.
	//
	refreshBilling() {
		// Get the billing.
		this.accountService.getBilling().subscribe(res => {
			this.billing = res;

			// Set toggle
			if (this.billing.Subscription == "Yearly") {
				this.checked = true;
			} else {
				this.checked = false;
			}
		});
	}

	//
	// Toggle clicked
	//
	checkedClick() {
		if (this.checked) {
			this.checked = false;
		} else {
			this.checked = true;
		}
	}

	//
	// Submit the plan change to the backend.
	//
	changePlan() {
		// Figure out new plan from checkbox
		if (this.checked) {
			this.billing.Subscription = "Yearly";
		} else {
			this.billing.Subscription = "Monthy";
		}

		// Send request to backend
		this.accountService.updateSubscription(this.billing.Subscription).subscribe(
			// Success
			() => {
				this.editModeToggle();
				this.refreshBilling();
			},

			// Error
			() => {
				alert("There was an error changing your subscription. Please contact help@skyclerk.com.");
			}
		);

		// Return false;
		return false;
	}
}

/* End File */
