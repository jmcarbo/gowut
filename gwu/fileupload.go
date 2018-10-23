// Copyright (C) 2013 Andras Belicza. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Defines the TextBox component.

package gwu

import (
	"net/http"
  "fmt"
//	"strconv"
)

// FileUpload interface defines a component for file upload purpose.
//
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the FileUpload is synchronized with the server
// on ETypeChange event which is when the TextBox loses focus
// or when the ENTER key is pressed.
//
// Default style class: "gwu-TextBox"
type FileUpload interface {
	// FileUpload is a component.
	Comp

	// FileUpload can be enabled/disabled.
	HasEnabled

}


// FileUpload implementation.
type fileUploadImpl struct {
	compImpl       // Component implementation
	hasEnabledImpl // Has enabled implementation

}


// NewFileUpload creates a new FileUpload.
func NewFileUpload() FileUpload {
	c := newFileUploadImpl(strEncURIThisV)
	c.Style().AddClass("gwu-FileUpload")
	return &c
}

// newFileUploadImpl creates a new fileUploadImpl.
func newFileUploadImpl(valueProviderJs []byte) fileUploadImpl {
	c := fileUploadImpl{newCompImpl(valueProviderJs), newHasEnabledImpl()}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (c *fileUploadImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for text box is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0
        /*
	value := r.FormValue(paramCompValue)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[paramCompValue] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
        */
}

var (
	strFileUpload  = []byte(`<input type="file"`) // `<input type="`
	//strInputCl  = []byte(`"/>`)           // `"/>`
)

func (c *fileUploadImpl) Render(w Writer) {
	w.Write(strFileUpload)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)
	w.Write(strInputCl)
        w.Write([]byte(fmt.Sprintf(`<input type="button" value="Upload" id="upload-button-%d" />`, c.id)))
        w.Write([]byte(fmt.Sprintf(`
          <script>
          var uploadBtn%d = document.getElementById('upload-button-%d');
          uploadBtn%d.onclick = function (evt) {
            var formData = new FormData();
            // Since this is the file only, we send it to a specific location
            var action = _pathUpload;
            // FormData only has the file
            var fileInput = document.getElementById('%d');
            var file = fileInput.files[0];
            formData.append('_pCompValue', file);
	    formData.append("_pCompId", "%d")
	    if (document.activeElement.id != null)
	      formData.append("_pFocCompId", document.activeElement.id)
            // Code common to both variants
            sendXHRequest(formData, action);
         }
         // Once the FormData instance is ready and we know
         // where to send the data, the code is the same
         // for both variants of this technique
         function sendXHRequest(formData, uri) {
           // Get an XMLHttpRequest instance
           var xhr = new XMLHttpRequest();
           // Set up events
           xhr.upload.addEventListener('loadstart', onloadstartHandler, false);
           xhr.upload.addEventListener('progress', onprogressHandler, false);
           xhr.upload.addEventListener('load', onloadHandler, false);
           xhr.addEventListener('readystatechange', onreadystatechangeHandler, false);
           // Set up request
           xhr.open('POST', uri, true);
           // Fire!
           xhr.send(formData);
         }
         // Handle the start of the transmission
         function onloadstartHandler(evt) {
           var div = document.getElementById('upload-status');
           div.innerHTML = 'Upload started.';
         }
         // Handle the end of the transmission
         function onloadHandler(evt) {
           var div = document.getElementById('upload-status');
           div.innerHTML += '<' + 'br>File uploaded. Waiting for response.';
         }
         // Handle the progress
         function onprogressHandler(evt) {
           var div = document.getElementById('progress');
           var percent = evt.loaded/evt.total*100;
           div.innerHTML = 'Progress: ' + percent + '%';
         }
         // Handle the response from the server
         function onreadystatechangeHandler(evt) {
           var status, text, readyState;
           try {
             readyState = evt.target.readyState;
             text = evt.target.responseText;
             status = evt.target.status;
           }
           catch(e) {
             return;
           }
           if (readyState == 4 && status == '200' && evt.target.responseText) {
             var status = document.getElementById('upload-status');
             status.innerHTML += '<' + 'br>Success!';
             var result = document.getElementById('result');
             result.innerHTML = '<p>The server saw it as:</p><pre>' + evt.target.responseText + '</pre>';
           }
         }
         </script>
  `, c.id, c.id, c.id, c.id, c.id)))
}


