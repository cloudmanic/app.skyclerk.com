//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService, LedgerResponse, LedgerSummaryResponse, LedgerPnlSummary } from 'src/app/services/ledger.service';
import { MeService } from 'src/app/services/me.service';
import { Category } from 'src/app/models/category.model';
import { Label } from 'src/app/models/label.model';
import { ActivatedRoute } from '@angular/router';
import { trigger, transition, animate, style } from '@angular/animations';
import { Ledger } from 'src/app/models/ledger.model';
import { Title } from '@angular/platform-browser';
import { environment } from 'src/environments/environment';
import { AccountService } from 'src/app/services/account.service';

const pageTitle: string = environment.title_prefix + "Ledger";

@Component({
	selector: 'app-landing',
	templateUrl: './landing.component.html',
	animations: [
		trigger('fadeIn', [
			transition(':enter', [style({ opacity: 0 }), animate(700)]),
			transition(':leave', animate(700, style({ opacity: 0 })))
		])
	]
})

export class LandingComponent implements OnInit {
	pageTitle: string = "Skyclerk | Dashboard Reports";
	page: number = 1;
	selected: Ledger[] = [];
	plSummary: LedgerPnlSummary = new LedgerPnlSummary(0, 0, 0);
	pageRangeSelect: number = 1;
	ledgerSummary: LedgerSummaryResponse = new LedgerSummaryResponse([], [], []);
	ledgers: LedgerResponse = new LedgerResponse(false, 0, 25, 0, []);

	// Active filters
	type: string = "";
	search: string = "";
	activeYear: number = null;
	activeLabels: Label[] = [];
	activeCategory: Category = null;

	// Add / Edit ledger stuff.
	showAddEditType: string = "income";
	showAddEditLedger: boolean = false;


	//
	// Construct
	//
	constructor(public ledgerService: LedgerService, public meService: MeService, public activeRoute: ActivatedRoute, private titleService: Title, public accountService: AccountService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Load starting data
		this.refreshLedger();

		// Listen for account changes.
		this.accountService.accountChange.subscribe(() => {
			this.refreshLedger();
		});

		// Listen for GET parm changes
		this.activeRoute.queryParams.subscribe(params => {
			if (typeof params['add'] == "undefined") {
				this.showAddEditLedger = false;
				return;
			}

			this.showAddEditType = params['add'];
			this.showAddEditLedger = true;
		});
	}

	//
	// Refresh ledger.
	//
	refreshLedger() {
		this.loadLedgerData();
		this.getLedgerSummary();
		this.getLedgerPnlSummary();
	}

	//
	// Refresh ledger on add
	//
	refreshFromAdd(_l: Ledger) {
		this.refreshLedger();
	}

	//
	// Load page data
	//
	loadPageData() {
		// Clear selected
		this.selected = [];

		// Reset filters, and such.
		this.page = 1;
		this.type = "";
		this.search = "";
		this.pageRangeSelect = 1;

		// Load ledger data
		this.refreshLedger();
	}

	//
	// Load ledger data.
	//
	loadLedgerData() {
		// Load ledger entries
		this.ledgerService.get(this.page, this.type, this.search, this.activeCategory, this.activeLabels, this.activeYear).subscribe(res => {
			this.ledgers = res;
			this.pageRangeSelect = this.page;
		});
	}

	//
	// getLedgerPnlSummary data.
	//
	getLedgerPnlSummary() {
		this.ledgerService.getLedgerPnlSummary(this.type, this.search, this.activeCategory, this.activeLabels, this.activeYear).subscribe(res => {
			this.plSummary = res;
		});
	}

	//
	// Get ledger summary
	//
	getLedgerSummary() {
		this.ledgerService.getLedgerSummary(this.type).subscribe(res => {
			this.ledgerSummary = res;
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
		this.refreshLedger();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.refreshLedger();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.refreshLedger();
	}

	//
	// Mark a ledger entry as selected
	//
	setSelected(entry: Ledger) {
		// See if we are deselecting
		for (let i = 0; i < this.selected.length; i++) {
			if (this.selected[i].Id == entry.Id) {
				this.selected.splice(i, 1);
				return;
			}
		}

		// Must be adding.
		this.selected.push(entry);
	}

	//
	// See if entry is checked.
	//
	isSelectedChecked(entry: Ledger) {
		// See if it is selected
		for (let i = 0; i < this.selected.length; i++) {
			if (this.selected[i].Id == entry.Id) {
				return true;
			}
		}

		return false;
	}

	//
	// Check all entries
	//
	checkAllEntries() {
		if (this.selected.length) {
			this.selected = [];
		} else {
			for (let i = 0; i < this.ledgers.Data.length; i++) {
				this.selected.push(this.ledgers.Data[i]);
			}
		}
	}

	//
	// Delete selected ledger
	//
	deleteSelectedEntries() {
		let c = confirm(`Are you sure you want to delete ${this.selected.length} ledger entries?`);

		if (!c) {
			return;
		}

		// Loop through and delete entries.
		for (let i = 0; i < this.selected.length; i++) {
			this.ledgerService.delete(this.selected[i]).subscribe();
		}

		// Small delay before reloading
		setTimeout(() => { this.loadPageData(); }, 1000);
	}

	// ------------- Filter Functions --------------- //

	//
	// Do we show filters
	//
	doWeShowFilters() {
		if (this.activeYear) {
			return true;
		}

		if (this.activeCategory) {
			return true;
		}

		if (this.activeLabels.length > 0) {
			return true;
		}

		return false;
	}

	//
	// Add Label fitter
	//
	addLabelFilter(lb: Label) {
		this.page = 1;
		this.activeLabels.push(lb);
		this.refreshLedger();
	}

	//
	// Remove Label fitter
	//
	removeLabelFilter(lb: Label) {
		this.page = 1;
		this.activeLabels = this.activeLabels.filter(obj => obj !== lb);
		this.refreshLedger();
	}

	//
	// Set year filter
	//
	setYearFilter(year: number) {
		this.page = 1;
		this.activeYear = year;
		this.refreshLedger();
	}

	//
	// Set category filter
	//
	setCategoryFilter(cat: Category) {
		this.page = 1;
		this.activeCategory = cat;
		this.refreshLedger();
	}

	//
	// Clear all sidebar filters
	//
	clearAllSidebarFilters() {
		this.page = 1;
		this.activeYear = null;
		this.activeLabels = [];
		this.activeCategory = null;
		this.refreshLedger();
	}

	//
	// Do search.
	//
	doSearch() {
		this.page = 1;
		this.refreshLedger();
	}

	//
	// Change the type we are filtring by
	//
	doTypeClick(type: string) {
		this.type = type;
		this.page = 1;
		this.refreshLedger();
	}

	//
	// Return the count of a year
	//
	filterGetYearCount(year: number) {
		// Match up the label
		for (let i = 0; i < this.ledgerSummary.Years.length; i++) {
			if (this.ledgerSummary.Years[i].Year == year) {
				return this.ledgerSummary.Years[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}

	//
	// Return the count of a label
	//
	filterGetLabelCount(lb: Label) {
		// Match up the label
		for (let i = 0; i < this.ledgerSummary.Labels.length; i++) {
			if (this.ledgerSummary.Labels[i].Id == lb.Id) {
				return this.ledgerSummary.Labels[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}

	//
	// Filter get category count
	//
	filterGetCategoryCount(cat: Category) {
		// Match up the category
		for (let i = 0; i < this.ledgerSummary.Categories.length; i++) {
			if (this.ledgerSummary.Categories[i].Id == cat.Id) {
				return this.ledgerSummary.Categories[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}
}

/* End File */
