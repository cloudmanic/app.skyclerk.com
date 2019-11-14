//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, EventEmitter, Output, Input, SimpleChanges } from '@angular/core';
import { CategoryService } from 'src/app/services/category.service';
import { Category } from 'src/app/models/category.model';

@Component({
	selector: 'app-categories-field',
	templateUrl: './field.component.html'
})
export class FieldComponent implements OnInit {
	@Input() type: string = "expense";
	@Input() theme: string = "add";
	@Input() errMsg: string = "";
	@Input() category: Category = new Category();
	@Output() onCategory = new EventEmitter<Category>();
	@Output() addCategoryToggle = new EventEmitter<boolean>();

	categories: Category[] = [];

	//
	// Constructor
	//
	constructor(public categoryService: CategoryService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Detect changes from properties.
	//
	ngOnChanges(changes: SimpleChanges) {
		// If we are in edit mode we do not do anything.
		if (this.theme == "edit") {
			this.getCategories(false);
			return;
		}

		// Detect type changes.
		if (typeof changes.type != "undefined") {
			this.category = new Category();
			this.getCategories(false);
		}

		// This is a change from adding a new category.
		if (typeof changes.category != "undefined") {
			this.deteactNewCategory();
		}
	}

	//
	// Detect a brand new category
	//
	deteactNewCategory() {
		// See if we have the current category in our list. If not refresh list.
		let found = false;

		for (let i = 0; i < this.categories.length; i++) {
			if (this.categories[i].Id == this.category.Id) {
				found = true;
			}
		}

		if (!found) {
			this.getCategories(true);
		}
	}

	//
	// On Category Change
	//
	onChange() {
		this.onCategory.emit(this.category);
	}

	//
	// Get categories
	//
	getCategories(refresh: boolean) {
		this.categoryService.get(this.type).subscribe(res => {
			this.categories = res;

			// Set the correct value (this is a little hacky)
			if (refresh) {
				for (let i = 0; i < this.categories.length; i++) {
					if (this.categories[i].Id == this.category.Id) {
						this.category = this.categories[i];
					}
				}
			}
		});
	}

	//
	// Toggle the add category popover
	//
	addCategoryToggleClick() {
		this.addCategoryToggle.emit(true);
	}
}

/* End File */
