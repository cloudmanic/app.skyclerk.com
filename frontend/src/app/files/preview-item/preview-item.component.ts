//
// Date: 2019-05-17
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';
import { FileUploadsWithStatus } from '../upload/upload.component';
import { File as FileModel } from '../../models/file.model';

@Component({
	selector: 'app-files-preview-item',
	templateUrl: './preview-item.component.html'
})

export class PreviewItemComponent implements OnInit {
	@Input() file: FileUploadsWithStatus = { file: null, status: "uploading", progress: 0, model: new FileModel() };

	//
	// Constructor.
	//
	constructor() {
		console.log(this.file);
	}

	//
	// ngOnInit
	//
	ngOnInit() { }
}

/* End File */
