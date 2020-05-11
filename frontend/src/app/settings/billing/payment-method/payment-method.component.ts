//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';
import { Billing } from 'src/app/models/billing.model';
import { environment } from 'src/environments/environment';

declare var Stripe: any;

@Component({
	selector: 'app-settings-billing-payment-method',
	templateUrl: './payment-method.component.html'
})

export class PaymentMethodComponent {
	errorMsg: string = "";
	editMode: boolean = false;
	billing: Billing = new Billing();
	stripe: any = null;
	cardNumber: any = null;
	cvc: any = null;
	expiry: any = null;
	elements: any = null;
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
		});
	}

	//
	// Toggle to edit mode.
	//
	editModeToggle() {
		// Toggle edit mode.
		this.editMode = !this.editMode;

		// Setup stripe on show.
		if (this.editMode) {
			// Setup stripe fields
			setTimeout(() => {
				this.stripe = Stripe(environment.stripe_pub_key);
				this.elements = this.stripe.elements();

				// Card number field.
				this.cardNumber = this.elements.create('cardNumber', {
					style: { base: { fontSize: '18px', '::placeholder': { color: '#efefef' } } },
					classes: { base: "field-stripe" }
				});
				this.cardNumber.mount('#card-number');

				// Card expiry field.
				this.expiry = this.elements.create('cardExpiry', {
					style: { base: { fontSize: '18px', '::placeholder': { color: '#efefef' } } },
					classes: { base: "field-stripe" }
				});
				this.expiry.mount('#card-expiry');

				// Card cvc field.
				this.cvc = this.elements.create('cardCvc', {
					style: { base: { fontSize: '18px', '::placeholder': { color: '#efefef' } } },
					classes: { base: "field-stripe" },
					placeholder: "123"
				});
				this.cvc.mount('#card-cvc');
			}, 100);
		}
	}

	//
	// Submit new credit card number to BE
	//
	updateCard() {
		// Clear error
		this.errorMsg = "";

		// Get stripe token from stripe
		this.stripe.createToken(this.cardNumber).then((result: any) => {
			// Is this an error?
			if (result.error) {
				this.errorMsg = result.error.message;
				return;
			}

			// Is this an success?
			if (result.token) {
				this.sendTokenToServer(result.token.id);
				return;
			}
		});
	}

	//
	// sendTokenToServer send token to server to for storage.
	//
	sendTokenToServer(token: string) {
		// Send to server.
		this.accountService.stripeToken(token, this.billing.Subscription).subscribe(
			// Success
			_res => {
				this.editModeToggle();
				let accountId = localStorage.getItem('account_id');
				this.accountService.accountChange.emit(Number(accountId));
			},

			// Error
			err => {
				if (err.error) {
					this.errorMsg = err.error.error;
				} else {
					this.errorMsg = "Very sorry, there was an error with your credit card. Please contact help@skyclerk.com";
				}
			}
		);
	}
}

/* End File */
