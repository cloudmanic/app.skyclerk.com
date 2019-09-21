//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { LabelService } from 'src/app/services/label.service';
import { Label } from 'src/app/models/label.model';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-settings-categories-labels-labels',
	templateUrl: './labels.component.html'
})

export class LabelsComponent {
	labels: Label[] = [];
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public labelService: LabelService, public accountService: AccountService) { }

	//
	// ngOnInit - since we have Input()'s we have to use ngOnInit instead of constructor
	//
	ngOnInit() {
		this.refreshLabels();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshLabels();
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
	// Refresh labels
	//
	refreshLabels() {
		// Get labels
		this.labelService.get().subscribe(res => {
			this.labels = res;
		});
	}

	//
	// Add label click.
	//
	addLabelClick() {
		let lb = new Label();
		lb.EditMode = true;
		this.labels.push(lb);
	}

	//
	// Save new label name.
	//
	save(row: Label) {
		if (row.Id > 0) {
			this.updateLabel(row);
		} else {
			this.createLabel(row);
		}
	}

	//
	// Delete label
	//
	deleteLabel(row: Label) {
		// Delete label on server.
		this.labelService.delete(row).subscribe(
			// Success
			() => {
				this.refreshLabels();
			},

			// Error
			err => {
				alert(err.error.error);
			}
		);
	}

	//
	// Create label
	//
	createLabel(row: Label) {
		// Send change to server.
		this.labelService.create(row).subscribe(
			// Success
			() => {
				this.refreshLabels();
			},

			// Error
			err => {
				row.ErrorMsg = err.error.errors.name;
			}
		);
	}

	//
	// Update label
	//
	updateLabel(row: Label) {
		// Send change to server.
		this.labelService.update(row).subscribe(
			// Success
			() => {
				this.refreshLabels();
			},

			// Error
			err => {
				row.ErrorMsg = err.error.errors.name;
			}
		);
	}

	//
	// Toggle Edit Mode
	//
	toggleEditMode(row: Label, open: string) {
		// Only one edit mode at a time.
		for (let i = 0; i < this.labels.length; i++) {
			this.labels[i].EditMode = false;
		}

		// Toggle open / close
		if (open == "open") {
			row.EditMode = true;
		} else {
			row.EditMode = false;

			// Clear value if the user updated it.
			this.refreshLabels();
		}
	}
}

/* End File */
