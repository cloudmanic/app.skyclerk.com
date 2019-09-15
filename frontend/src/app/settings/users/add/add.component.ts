//
// Date: 2019-08-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { UserService } from 'src/app/services/user.service';
import { Router } from '@angular/router';

@Component({
	selector: 'app-add',
	templateUrl: './add.component.html'
})

export class AddComponent implements OnInit {
	first: string = "";
	last: string = "";
	email: string = "";
	message: string = "";
	errMsg: string = "";

	//
	// Constructor
	//
	constructor(public router: Router, public userService: UserService) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Submit invite request.
	//
	submit() {
		// Send invite request to server.
		this.userService.invite(this.first, this.last, this.email, this.message).subscribe(
			// success
			() => {
				this.router.navigate(['/settings/users'], { queryParams: { success: "true" } });
			},

			// Error
			(err) => {
				this.errMsg = err.error.error;
			}
		);
	}
}

/* End File */
