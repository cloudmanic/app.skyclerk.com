//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, EventEmitter, Output, Input } from '@angular/core';
import { CategoryService } from 'src/app/services/category.service';
import { Category } from 'src/app/models/category.model';

@Component({
	selector: 'app-categories-pop-out',
	templateUrl: './pop-out.component.html'
})

export class PopOutComponent implements OnInit {
	@Output() onCategory = new EventEmitter<Category>();
	@Input() type: string = "expense";
	@Input() show: boolean = false;

	name: string = "";
	errMsg: string = "";

	//
	// Constructor
	//
	constructor(public categoryService: CategoryService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Create a new category
	//
	save() {
		// Clear error
		this.errMsg = "";

		// Setup the category
		let cat = new Category();
		cat.Name = this.name;
		cat.Type = this.type;

		// Save category with BE
		this.categoryService.create(cat).subscribe(
			res => {
				this.name = "";
				this.onCategory.emit(res);
			},

			err => {
				this.errMsg = err.error.errors.name;
			}
		);
	}
}

/* End File */
