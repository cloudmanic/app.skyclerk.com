//
// Date: 2019-05-20
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService } from 'src/app/services/ledger.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Ledger } from 'src/app/models/ledger.model';

@Component({
	selector: 'app-ledger-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	ledger: Ledger = new Ledger();

	//
	// Constructor
	//
	constructor(public ledgerService: LedgerService, public route: ActivatedRoute, public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Is this an edit action?
		let ledgerId = this.route.snapshot.params['id'];

		// Get the ledger based on the id we passed in.
		this.ledgerService.getById(ledgerId).subscribe(res => {
			this.ledger = res;
		});
	}

	//
	// Delete ledger
	//
	deleteLedger() {
		let c = confirm("Are you sure you want to delete this ledger entry?");

		if (!c) {
			return;
		}

		// Send delete request.
		this.ledgerService.delete(this.ledger).subscribe(() => {
			this.router.navigate(['/ledger']);
		});
	}

}

/* End File */
