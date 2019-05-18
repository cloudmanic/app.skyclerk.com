//
// Date: 2019-05-17
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { UploadFile, UploadEvent, FileSystemFileEntry } from 'ngx-file-drop';
import { FileService } from 'src/app/services/file.service';
import { HttpErrorResponse } from '@angular/common/http';
import { File as FileModel } from '../../models/file.model';

@Component({
	selector: 'app-files-upload',
	templateUrl: './upload.component.html'
})

export class UploadComponent implements OnInit {
	files: FileUploadsWithStatus[] = [];

	//
	// Constructor.
	//
	constructor(public fileService: FileService, public ref: ChangeDetectorRef) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// When a file is dropped. Also when clicks on a file upload
	//
	fileDropped(event: UploadEvent) {
		// Add the files
		for (let i = 0; i < event.files.length; i++) {
			let f = event.files[i];

			// Is it a file?
			if (f.fileEntry.isFile) {
				let t = { file: f, status: "uploading", progress: 0, model: new FileModel() }
				this.files.push(t);
				this.uploadFile(t);
			}
		}
	}

	//
	// Upload file to server.
	//
	uploadFile(t: FileUploadsWithStatus) {
		const fileEntry = t.file.fileEntry as FileSystemFileEntry;

		// Clear errors
		//this.errMsg = "";

		// Get the file from the entry
		fileEntry.file((file: File) => {

			this.fileService.upload(file).subscribe(
				// Success
				(res) => {

					// Is this a progress upload.
					if (typeof res.message != "undefined") {
						t.progress = Number(res.message);

						if (t.progress == 100) {
							t.status = "processing";
						}
						return;
					}

					// Updated model
					if (typeof res == "object") {
						t.model = res;
						t.status = "done";
						console.log(t);
						return;
					}




					// t.status = "";
					// t.fileId = res.ID;
					// t.file.relativePath = res.Name;
					// this.onUpload.emit({ fileId: res.ID, docType: this.docPurpose });
				},

				// Error
				(err: HttpErrorResponse) => {

					console.log(err);

					// t.status = "Error";
					// if (err.error.fields.length) {
					// 	this.errMsg = err.error.fields[0].reason;
					// }
				}

			);

		});
	}

}

// FileUploadsWithStatus - keep files we are uploading and the status.
export interface FileUploadsWithStatus {
	file: UploadFile;
	model: FileModel;
	status: string;
	progress: number;
}

/* End File */
