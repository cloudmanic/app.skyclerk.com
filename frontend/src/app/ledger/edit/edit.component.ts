//
// Date: 2019-05-20
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService } from 'src/app/services/ledger.service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
	selector: 'app-ledger-edit',
	templateUrl: './edit.component.html'
})
export class EditComponent implements OnInit {
	ledgerId: number = 0;

	//
	// Constructor
	//
	constructor(public ledgerService: LedgerService, public route: ActivatedRoute, public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Is this an edit action?
		this.ledgerId = this.route.snapshot.params['id'];
	}

	//
	// Save
	//
	save() {
		this.router.navigate([`/ledger/${this.ledgerId}`]);
	}
}

/* End File */
