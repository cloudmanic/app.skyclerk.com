//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccountService } from 'src/app/services/account.service';
import { environment } from 'src/environments/environment';
import { Subject } from 'rxjs';

declare var Stripe: any;

@Component({
	selector: 'app-payment',
	templateUrl: './payment.component.html'
})

export class PaymentComponent implements OnInit {
	back: string = "";
	plan: string = "Monthly";
	today: number = Date.now();
	errorMsg: string = "";
	stripe: any = null;
	cardNumber: any = null;
	cvc: any = null;
	expiry: any = null;
	elements: any = null;
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor.
	//
	constructor(public route: ActivatedRoute, public accountService: AccountService, public router: Router) { }

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

		// Set plan on load.
		this.plan = this.route.snapshot.queryParamMap.get("plan");

		// When plan changes
		this.route.queryParamMap.subscribe(queryParams => {
			this.plan = queryParams.get("plan");
		});

		// Setup stripe fields
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
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Submit payment.
	//
	submitPayment() {
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

		return false;
	}

	//
	// sendTokenToServer send token to server to for storage.
	//
	sendTokenToServer(token: string) {
		// Send to server.
		this.accountService.stripeToken(token, this.plan).subscribe(
			// Success
			_res => {
				this.router.navigate(['/paywall/success']);
			},

			// Error
			err => {
				console.log(err);
				alert("Very sorry, there was an error with your credit card. Please contact help@skyclerk.com");
			}
		);
	}
}

/* End File */
