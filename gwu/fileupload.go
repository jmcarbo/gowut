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
  "io"
  "io/ioutil"
  "path"
  "os"
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

  filename string
  originalFilename string
}


// NewFileUpload creates a new FileUpload.
func NewFileUpload() FileUpload {
	c := newFileUploadImpl(strEncURIThisV)
	c.Style().AddClass("gwu-FileUpload")
	return &c
}

// newFileUploadImpl creates a new fileUploadImpl.
func newFileUploadImpl(valueProviderJs []byte) fileUploadImpl {
	c := fileUploadImpl{newCompImpl(valueProviderJs), newHasEnabledImpl(), "", ""}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (c *fileUploadImpl) preprocessEvent(event Event, r *http.Request) {
  file, handler, err := r.FormFile("cval")
  if err != nil {
    fmt.Println(err)
    return
  }
  defer file.Close()
  //fmt.Printf("%+v\n", handler)
  myPath:=path.Join("tempfiles")
  os.MkdirAll(myPath, 0755)
  f, err :=ioutil.TempFile(myPath, "*_" + handler.Filename)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer f.Close()
  io.Copy(f, file)
  c.filename=f.Name()
  c.originalFilename=handler.Filename
}

var (
	strFileUpload  = []byte(`<input type="file"`) // `<input type="`
	//strInputCl  = []byte(`"/>`)           // `"/>`
)

func (c *fileUploadImpl) Render(w Writer) {
	w.Write(strFileUpload)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	//c.renderEHandlers(w)
	w.Write(strInputCl)
        w.Write([]byte(fmt.Sprintf(`<input type="button" value="Upload" id="upload-button-%d" />`, c.id)))
        w.Write([]byte(fmt.Sprintf(`<div id="upload-status-%d"></div>`, c.id)))
        w.Write([]byte(fmt.Sprintf(`<div id="progress-%d"></div>`, c.id)))
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
            formData.append('et', 11);
            formData.append('cval', file);
	    formData.append("cid", "%d")
	    if (document.activeElement.id != null)
	      formData.append("fcid", document.activeElement.id)
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
           var div = document.getElementById('upload-status-%d');
           div.innerHTML = 'Upload started.';
         }
         // Handle the end of the transmission
         function onloadHandler(evt) {
           var div = document.getElementById('upload-status-%d');
           div.innerHTML += '<' + 'br>File uploaded. Waiting for response.';
         }
         // Handle the progress
         function onprogressHandler(evt) {
           var div = document.getElementById('progress-%d');
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
             var status = document.getElementById('upload-status-%d');
             status.innerHTML += '<' + 'br>Success!';
             var result = document.getElementById('result');
             result.innerHTML = '<p>The server saw it as:</p><pre>' + evt.target.responseText + '</pre>';
           }
         }
         </script>
  `, c.id, c.id, c.id, c.id, c.id, c.id, c.id, c.id, c.id)))
}


