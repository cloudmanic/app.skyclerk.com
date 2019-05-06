//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService, LedgerResponse } from 'src/app/services/ledger.service';

@Component({
	selector: 'app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	page: number = 1;
	type: string = "";
	search: string = "";
	pageRangeSelect: number = 1;
	ledgers: LedgerResponse = new LedgerResponse(false, 0, 50, 0, []);

	//
	// Construct
	//
	constructor(public ledgerService: LedgerService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		this.loadLedgerData();
	}

	//
	// Load ledger data.
	//
	loadLedgerData() {
		this.ledgerService.get(this.page, this.type, this.search).subscribe(res => {
			this.ledgers = res;
			this.pageRangeSelect = this.page;
		});
	}

	//
	// Return the page list for ledger
	//
	getPageRange() {
		let pages = [];

		if (this.ledgers.Data.length == 0) {
			return [1];
		}

		let pageCount = Math.ceil(this.ledgers.NoLimitCount / this.ledgers.Limit);

		for (let i = 1; i <= pageCount; i++) {
			pages.push(i);
		}

		return pages;
	}

	//
	// Paging select change
	//
	doPageSelectChange() {
		this.page = this.pageRangeSelect;
		this.loadLedgerData();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.loadLedgerData();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.loadLedgerData();
	}
}

/* End File */
