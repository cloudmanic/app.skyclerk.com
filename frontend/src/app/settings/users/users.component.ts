//
// Date: 2019-08-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { UserService } from 'src/app/services/user.service';
import { User } from 'src/app/models/user.model';
import { Me } from 'src/app/models/me.model';
import { MeService } from 'src/app/services/me.service';
import { ActivatedRoute } from '@angular/router';

const pageTitle: string = environment.title_prefix + "Settings Users";

@Component({
	selector: 'app-users',
	templateUrl: './users.component.html'
})

export class UsersComponent implements OnInit {
	me: Me = new Me();
	users: User[] = [];
	successMsg: string = "";

	//
	// Construct.
	//
	constructor(public titleService: Title, public meService: MeService, public userService: UserService, public route: ActivatedRoute) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		this.route.queryParams.subscribe(params => {
			if (params["success"]) {
				this.successMsg = "User successfully invited!";
				setTimeout(() => { this.successMsg = ""; }, 3000);
			}
		});

		// Get Me.
		this.meService.get().subscribe(res => {
			this.me = res;
		});

		// Get the users.
		this.userService.get().subscribe((res) => {
			this.users = res;
		});

		// TODO(spicer): Get the account and the owner. We can't delete an account owner.
	}
}

/* End File */
