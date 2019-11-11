//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Component, OnInit, Input, EventEmitter, Output } from '@angular/core';
import { Contact } from 'src/app/models/contact.model';
import { Ledger } from 'src/app/models/ledger.model';
import { Category } from 'src/app/models/category.model';
import { Label } from 'src/app/models/label.model';
import { File as FileModel } from 'src/app/models/file.model';
import { Router } from '@angular/router';
import { LedgerService } from 'src/app/services/ledger.service';
import { MeService } from 'src/app/services/me.service';
import { Subject } from 'rxjs';
import { AccountService } from 'src/app/services/account.service';
import { Numbers } from 'src/app/library/numbers';

@Component({
	selector: 'app-ledger-add',
	templateUrl: './add.component.html'
})
export class AddComponent implements OnInit {
	@Input() type: string = "income";
	@Output() refreshLedger = new EventEmitter<Ledger>();

	errors: any = [];
	amount: string = "";
	ledger: Ledger = new Ledger();
	showAddLabel: boolean = false;
	showAddContact: boolean = false;
	showAddCategory: boolean = false;
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public ledgerService: LedgerService, public router: Router, public meService: MeService, public accountService: AccountService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.router.navigate([`/ledger`]);
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
	// Save the new ledger entry.
	//
	save() {
		// Clear errors
		this.errors = [];

		// Clean up ledger Amount
		this.ledger.Amount = Numbers.toFloat(this.amount);

		// Is this an expense
		if (this.type == "expense") {
			this.ledger.Amount = (this.ledger.Amount * -1);
		}

		// Send this ledger to BE.
		this.ledgerService.create(this.ledger).subscribe(
			// Sucesss
			res => {
				this.refreshLedger.emit(res);

				// Hack to refresh P&L
				this.accountService.accountChange.emit(this.ledger.AccountId);

				this.router.navigate(['/ledger']);
			},

			// Error
			err => {
				// Reset negative value
				if (this.type == "expense") {
					this.ledger.Amount = (this.ledger.Amount * -1);
				}

				// Show errors
				this.errors = err.error.errors;
			}
		);
	}

	//
	// We call this on assigning a contact.
	//
	assignContact(contact: Contact) {
		this.showAddContact = false;
		this.ledger.Contact = contact;
	}

	//
	// We call this whenever someone checks or unchecks a label
	//
	onLabelsChange(lbs: Label[]) {
		this.ledger.Labels = lbs;
	}

	//
	// We call this when a file is added to the ledger.
	//
	onAddFile(f: FileModel) {
		this.ledger.Files.push(f);
	}

	//
	// Call this when we delete a file.
	//
	onDeleteFile(f: FileModel) {
		for (let i = 0; i < this.ledger.Files.length; i++) {
			if (f.Id == this.ledger.Files[i].Id) {
				this.ledger.Files.splice(i, 1);
			}
		}
	}

	//
	// We call this after assigning a label.
	//
	assignLabel(lb: Label) {
		this.showAddLabel = false;
		this.ledger.Labels.push(lb);

		// Hack to make sure things update.
		this.onLabelsChange(this.ledger.Labels.slice());
	}

	//
	// We call this after assigning a category.
	//
	assignCategory(category: Category) {
		this.showAddCategory = false;
		this.ledger.Category = category;
	}

	//
	// Add label click
	//
	addLabelToggle() {
		if (this.showAddLabel) {
			this.showAddLabel = false;
		} else {
			this.showAddLabel = true;
		}
	}

	//
	// Add category click
	//
	addCategoryToggle() {
		if (this.showAddCategory) {
			this.showAddCategory = false;
		} else {
			this.showAddCategory = true;
		}
	}

	//
	// Add contact click
	//
	addContactToggle() {
		if (this.showAddContact) {
			this.showAddContact = false;
		} else {
			this.showAddContact = true;
		}
	}

	//
	// Check to see if we have errors for a field.
	//
	hasError(field: string) {
		if (this.errors[field]) {
			return this.errors[field];
		}

		return "";
	}
}

/* End File */
