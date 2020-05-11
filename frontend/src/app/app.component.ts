//
// Date: 2019-08-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { PingService } from './services/ping.service';

declare let _paq: any;
declare let gtag: Function;

@Component({
	selector: 'app-root',
	templateUrl: './app.component.html'
})

export class AppComponent {
	title = 'Skyclerk';

	//
	// Construct - Load services that need to run site wide
	//
	constructor(private router: Router, private titleService: Title, private pingService: PingService) {
		// Start server ping.
		this.pingService.startPing();

		// subscribe to router events and send page views to Analytics
		this.router.events.subscribe(event => {
			if (event instanceof NavigationEnd) {
				// We give it a timeout so we give time for the title to update.
				setTimeout(() => {
					// Set user id for piwik
					let email = localStorage.getItem('user_email');

					if (email.length) {
						//_paq.push(['setUserId', email]);

						// We do this instead of "setUserId" since it creates new logs for the user if they are
						// not logged in and we want to track public website actions too.
						_paq.push(['setCustomVariable', 1, "Email", email, "visit"]);
					}

					_paq.push(['setCustomUrl', event.urlAfterRedirects]);
					_paq.push(['setDocumentTitle', this.titleService.getTitle()]);
					_paq.push(['setGenerationTimeMs', 0]);
					_paq.push(['trackPageView']);
					_paq.push(['enableLinkTracking']); // Should be at end.

					// Google Analytics
					gtag('config', 'UA-102266466-2', { 'page_path': event.urlAfterRedirects });
				}, 50);
			}
		});
	}
}
