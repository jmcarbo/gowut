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

// w.AddHeadHTML(`<script src="https://cdn.jsdelivr.net/npm/@json-editor/json-editor/dist/jsoneditor.min.js"></script>`)

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

var editorConf = `<script>
    JSONEditor.defaults.languages.es = {
      error_notset: "Cal informar la propietat",
      error_notempty: "Cal un valor",
      error_enum: "Value must be one of the enumerated values",
      error_anyOf: "Value must validate against at least one of the provided schemas",
      error_oneOf: 'Value must validate against exactly one of the provided schemas. It currently validates against {{0}} of the schemas.',
      error_not: "Value must not validate against the provided schema",
      error_type_union: "Value must be one of the provided types",
      error_type: "Value must be of type {{0}}",
      error_disallow_union: "Value must not be one of the provided disallowed types",
      error_disallow: "Value must not be of type {{0}}",
      error_multipleOf: "Value must be a multiple of {{0}}",
      error_maximum_excl: "Value must be less than {{0}}",
      error_maximum_incl: "Value must be at most {{0}}",
      error_minimum_excl: "Value must be greater than {{0}}",
      error_minimum_incl: "Value must be at least {{0}}",
      error_maxLength: "Value must be at most {{0}} characters long",
      error_minLength: "Value must be at least {{0}} characters long",
      error_pattern: "Value must match the pattern {{0}}",
      error_additionalItems: "No additional items allowed in this array",
      error_maxItems: "Value must have at most {{0}} items",
      error_minItems: "Value must have at least {{0}} items",
      error_uniqueItems: "Array must have unique items",
      error_maxProperties: "Object must have at most {{0}} properties",
      error_minProperties: "Object must have at least {{0}} properties",
      error_required: "Object is missing the required property '{{0}}'",
      error_additional_properties: "No additional properties allowed, but property {{0}} is set",
      error_dependency: "Must have property {{0}}",
      error_date: 'Date must be in the format {{0}}',
      error_time: 'Time must be in the format {{0}}',
      error_datetime_local: 'Datetime must be in the format {{0}}',
      error_invalid_epoch: 'Date must be greater than 1 January 1970',
      button_delete_all: "Tot",
      button_delete_all_title: "Esborrar tot",
      button_delete_last: "Últim {{0}}",
      button_delete_last_title: "Delete Last {{0}}",
      button_add_row_title: "Afegir {{0}}",
      button_move_down_title: "More avall",
      button_move_up_title: "Moure amunt",
      button_delete_row_title: "Esborrar {{0}}",
      button_delete_row_title_short: "Esborrar",
      button_collapse: "Collapse",
      button_expand: "Expand",
      flatpickr_toggle_button: "Toggle",
      flatpickr_clear_button: "Netejar"
    };
    
    JSONEditor.defaults.default_language = 'es';
    </script>
`

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
        theme: 'bootstrap4',
        disable_collapse: true,
	disable_edit_json: true,
	disable_properties: true,
        show_errors: "change",
        iconlib: "bootstrap3"
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
        theme: 'bootstrap4',
        disable_collapse: true,
	disable_edit_json: true,
	disable_properties: true,
        show_errors: "change",
        iconlib: "bootstrap3"
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


