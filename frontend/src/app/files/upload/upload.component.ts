//
// Date: 2019-05-17
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, ChangeDetectorRef, Output, EventEmitter, Input, SimpleChanges } from '@angular/core';
import { UploadFile, UploadEvent, FileSystemFileEntry } from 'ngx-file-drop';
import { FileService } from 'src/app/services/file.service';
import { HttpErrorResponse } from '@angular/common/http';
import { File as FileModel } from '../../models/file.model';

@Component({
	selector: 'app-files-upload',
	templateUrl: './upload.component.html'
})

export class UploadComponent implements OnInit {
	@Output() onUpload = new EventEmitter<FileModel>();
	@Output() onDeleteFile = new EventEmitter<FileModel>();

	@Input() filesInput: FileModel[] = [];

	files: FileUploadsWithStatus[] = [];

	//
	// Constructor.
	//
	constructor(public fileService: FileService, public ref: ChangeDetectorRef) { }

	//
	// ngOnInit
	//
	ngOnInit() {

		console.log(this.filesInput);

	}

	//
	// Detect changes from properties.
	//
	ngOnChanges(changes: SimpleChanges) {
		// Detect type changes.
		if (typeof changes.filesInput != "undefined") {
			this.files = [];

			for (let i = 0; i < this.filesInput.length; i++) {
				this.files.push({ file: null, status: "done", progress: 0, error: "", model: this.filesInput[i] });
			}
		}
	}

	//
	// Delete a file
	//
	deleteClick(f: FileModel) {
		for (let i = 0; i < this.files.length; i++) {
			if (f.Id == this.files[i].model.Id) {
				this.files.splice(i, 1);
			}
		}

		// Tell parent about the delete.
		this.onDeleteFile.emit(f);
	}

	//
	// When a file is dropped. Also when clicks on a file upload
	//
	fileDropped(event: UploadEvent) {
		// Add the files
		for (let i = 0; i < event.files.length; i++) {
			let f = event.files[i];

			// Is it a file?
			if (f.fileEntry.isFile) {
				let t = { file: f, status: "uploading", progress: 0, error: "", model: new FileModel() }
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

						// Send uploaded file to the parent.
						this.onUpload.emit(t.model);
						return;
					}

				},

				// Error
				(err: HttpErrorResponse) => {
					t.status = "error";
					t.error = "Unknown error please contact help@options.cafe";

					// Should only be the file field
					if (typeof err.error.errors.file != "undefined") {
						t.error = err.error.errors.file;
					}
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
	error: string;
}

/* End File */
