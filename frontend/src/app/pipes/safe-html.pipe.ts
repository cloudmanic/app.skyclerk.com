//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

@Pipe({
	name: 'safeHtml'
})
export class SafeHtmlPipe implements PipeTransform {
	constructor(private sanitizer: DomSanitizer) { }

	transform(value: any, _args?: any): any {
		return this.sanitizer.bypassSecurityTrustHtml(value);
	}
}
