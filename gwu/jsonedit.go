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
type JSONEdit interface {
	// FileUpload is a component.
	Comp

	HasText
	// FileUpload can be enabled/disabled.
	HasEnabled

	SetSchema(s string)
}


// FileUpload implementation.
type jsonEditImpl struct {
  compImpl       // Component implementation
  hasTextImpl // Has enabled implementation
  hasEnabledImpl // Has enabled implementation

  schema string
}


// NewJSONEdit creates a new FileUpload.
func NewJSONEdit() JSONEdit {
	c := newJSONEditImpl(strEncURIThisV)
	c.Style().AddClass("gwu-JSONEdit")
	return &c
}

// newFileUploadImpl creates a new fileUploadImpl.
func newJSONEditImpl(valueProviderJs []byte) jsonEditImpl {
	c := jsonEditImpl{newCompImpl(valueProviderJs), newHasTextImpl(""), newHasEnabledImpl(), ""}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (c *jsonEditImpl) SetSchema(s string) {
	c.schema = s
}

func (c *jsonEditImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for text box is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0
	value := r.FormValue(paramCompValue)
        //fmt.Printf("Getting value %+v\n", value)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[paramCompValue] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
}


func (c *jsonEditImpl) Render(w Writer) {
/*
	w.Write(strFileUpload)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	//c.renderEHandlers(w)
	w.Write(strInputCl)
*/
        w.Write([]byte(fmt.Sprintf(`<div id="%d" />`, c.id)))
	if c.text != "" {
        w.Write([]byte(fmt.Sprintf(`
<script>
      // Initialize the editor with a JSON schema
      var editor%d = new JSONEditor(document.getElementById('%d'),{
        schema: %s,
	startval: %s,
        theme: 'bootstrap4'
      });
      
      // Hook up the submit button to log to the console
      //document.getElementById('submit').addEventListener('click',function() {
        // Get the value from the editor
       // console.log(editor.getValue());
      //});
      editor%d.on('change',function() {
  	// Do something
        value = editor%d.getValue()
        console.log("The data has changed!" + JSON.stringify(value) ); 
        se2(null,11,%d,JSON.stringify(value))
       });
    </script>
  `, c.id, c.id, c.schema, c.text, c.id, c.id, c.id)))
	} else {
        w.Write([]byte(fmt.Sprintf(`
<script>
      // Initialize the editor with a JSON schema
      var editor%d = new JSONEditor(document.getElementById('%d'),{
        schema: %s,
        theme: 'bootstrap4'
      });
      
      editor%d.on('change',function() {
  	// Do something
        value = editor%d.getValue()
        console.log("The data has changed!" + JSON.stringify(value) ); 
        se2(null,11,%d,JSON.stringify(value))
       });
	
    </script>
  `, c.id, c.id, c.schema, c.id, c.id, c.id)))
	}
}


