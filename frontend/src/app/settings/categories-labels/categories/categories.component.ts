//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, Input, OnInit } from '@angular/core';
import { CategoryService } from 'src/app/services/category.service';
import { Category } from 'src/app/models/category.model';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-settings-categories-labels-categories',
	templateUrl: './categories.component.html'
})

export class CategoriesComponent implements OnInit {
	@Input() type: string = "income";

	categories: Category[] = [];
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public categoryService: CategoryService, public accountService: AccountService) { }

	//
	// ngOnInit - since we have Input()'s we have to use ngOnInit instead of constructor
	//
	ngOnInit() {
		this.refreshCategories();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshCategories();
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
	// Refresh categories
	//
	refreshCategories() {
		// Get categories by type.
		this.categoryService.get(this.type).subscribe(res => {
			//console.log(res);
			this.categories = res;
		});
	}

	//
	// Add category click.
	//
	addCategoryClick() {
		let cat = new Category();
		cat.EditMode = true;
		cat.Type = this.type;
		this.categories.push(cat);
	}

	//
	// Save new category name.
	//
	save(row: Category) {
		if (row.Id > 0) {
			this.updateCategory(row);
		} else {
			this.createCategory(row);
		}
	}

	//
	// Delete category
	//
	deleteCategory(row: Category) {
		// Delete category on server.
		this.categoryService.delete(row).subscribe(
			// Success
			() => {
				this.refreshCategories();
			},

			// Error
			err => {
				alert(err.error.error);
			}
		);
	}

	//
	// Create category
	//
	createCategory(row: Category) {
		// Send change to server.
		this.categoryService.create(row).subscribe(
			// Success
			() => {
				this.refreshCategories();
			},

			// Error
			err => {
				row.ErrorMsg = err.error.errors.name;
			}
		);
	}

	//
	// Update category
	//
	updateCategory(row: Category) {
		// Send change to server.
		this.categoryService.update(row).subscribe(
			// Success
			() => {
				this.refreshCategories();
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
	toggleEditMode(row: Category, open: string) {
		// Only one edit mode at a time.
		for (let i = 0; i < this.categories.length; i++) {
			this.categories[i].EditMode = false;
		}

		// Toggle open / close
		if (open == "open") {
			row.EditMode = true;
		} else {
			row.EditMode = false;

			// Clear value if the user updated it.
			this.refreshCategories();
		}
	}

}

/* End File */
