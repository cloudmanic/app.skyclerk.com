//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input, EventEmitter, Output, SimpleChanges } from '@angular/core';
import { Label } from 'src/app/models/label.model';
import { LabelService } from 'src/app/services/label.service';

@Component({
	selector: 'app-labels-field',
	templateUrl: './field.component.html'
})

export class FieldComponent implements OnInit {
	@Input() selectedLabels: Label[] = [];
	@Output() onLabelsChange = new EventEmitter<Label[]>();
	@Output() addLabelToggle = new EventEmitter<boolean>();

	labels: Label[] = [];

	//
	// Constructor
	//
	constructor(public labelService: LabelService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Load page data.
		this.getCategories();
	}

	//
	// Detect changes from properties.
	//
	ngOnChanges(changes: SimpleChanges) {
		// This is a change from adding a new label.
		if (typeof changes.selectedLabels != "undefined") {
			this.getCategories();
		}
	}

	//
	// Get labels
	//
	getCategories() {
		this.labelService.get().subscribe(res => {
			this.labels = res;
		});
	}

	//
	// Show add Label popout
	//
	showAddLabel() {
		this.addLabelToggle.emit(true);
	}

	//
	// Is the label checked
	//
	isChecked(lb: Label) {
		// Search selected label.
		for (let i = 0; i < this.selectedLabels.length; i++) {
			if (this.selectedLabels[i].Id == lb.Id) {
				return true;
			}
		}

		return false;
	}

	//
	// When we check a label
	//
	onChange(lb: Label) {
		// See if this is a uncheck
		for (let i = 0; i < this.selectedLabels.length; i++) {
			if (this.selectedLabels[i].Id == lb.Id) {
				this.selectedLabels.splice(i, 1);
				this.onLabelsChange.emit(this.selectedLabels);
				return true;
			}
		}

		// Add to selected
		this.selectedLabels.push(lb);
		this.onLabelsChange.emit(this.selectedLabels);
	}

}

/* End File */
