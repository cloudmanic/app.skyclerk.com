//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit, Input, EventEmitter, Output } from '@angular/core';
import { LabelService } from 'src/app/services/label.service';
import { Label } from 'src/app/models/label.model';

@Component({
	selector: 'app-labels-pop-out',
	templateUrl: './pop-out.component.html'
})

export class PopOutComponent implements OnInit {
	@Input() show: boolean = false;
	@Output() onLabel = new EventEmitter<Label>();

	name: string = "";
	errMsg: string = "";

	//
	// Constructor
	//
	constructor(public labelService: LabelService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Create a new label
	//
	save() {
		// Clear error
		this.errMsg = "";

		// Setup the category
		let lb = new Label();
		lb.Name = this.name;

		// Save category with BE
		this.labelService.create(lb).subscribe(
			res => {
				this.name = "";
				this.onLabel.emit(res);
			},

			err => {
				this.errMsg = err.error.errors.name;
			}
		);
	}
}

/* End File */
