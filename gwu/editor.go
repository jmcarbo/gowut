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

// Defines the Editor component.

package gwu

import (
	"net/http"
  "fmt"
//	"strconv"
)

// Editor interface defines a component for text input purpose.
//
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the TextBox is synchronized with the server
// on ETypeChange event which is when the TextBox loses focus
// or when the ENTER key is pressed.
// If you want a TextBox to synchronize values during editing
// (while you type in characters), add the ETypeKeyUp event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETypeKeyUp)
//
// Default style class: "gwu-TextBox"
type Editor interface {
	// Editor is a component.
	Comp

	// Editor has text.
	HasText

	// Editor can be enabled/disabled.
	HasEnabled

	// ReadOnly returns if the editor is read-only.
	ReadOnly() bool

	// SetReadOnly sets if the text box is read-only.
	SetReadOnly(readOnly bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this a simple, one-line input text box,
	// rows>1 will make this a text area.
	SetRows(rows int)

	// Cols returns the number of displayed columns.
	Cols() int

	// SetCols sets the number of displayed columns.
	SetCols(cols int)

	// MaxLength returns the maximum number of characters
	// allowed in the text box.
	// -1 is returned if there is no maximum length set.
	MaxLength() int

	// SetMaxLength sets the maximum number of characters
	// allowed in the text box.
	// Pass -1 to not limit the maximum length.
	SetMaxLength(maxLength int)
}

// TextBox implementation.
type editorImpl struct {
	compImpl       // Component implementation
	hasTextImpl    // Has text implementation
	hasEnabledImpl // Has enabled implementation

	rows, cols int  // Number of displayed rows and columns.
}

var (
//	strEncURIThisV = []byte("encodeURIComponent(this.value)") // "encodeURIComponent(this.value)"
)

// NewEditor creates a new Editor
func NewEditor(text string) *editorImpl {
	c := newEditorImpl(strEncURIThisV, text, false)
	c.Style().AddClass("gwu-Editor")
	return &c
}

// newTextBoxImpl creates a new textBoxImpl.
func newEditorImpl(valueProviderJs []byte, text string, isPassw bool) editorImpl {
	c := editorImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl(), 1, 20}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (c *editorImpl) ReadOnly() bool {
	ro := c.Attr("readonly")
	return len(ro) > 0
}

func (c *editorImpl) SetReadOnly(readOnly bool) {
	if readOnly {
		c.SetAttr("readonly", "readonly")
	} else {
		c.SetAttr("readonly", "")
	}
}

func (c *editorImpl) Rows() int {
	return c.rows
}

func (c *editorImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *editorImpl) Cols() int {
	return c.cols
}

func (c *editorImpl) SetCols(cols int) {
	c.cols = cols
}

func (c *editorImpl) preprocessEvent(event Event, r *http.Request) {
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


var (
	strDivOp   = []byte("<div")   // "<div"
	//strRows         = []byte(` rows="`)     // ` rows="`
	//strCols         = []byte(`" cols="`)    // `" cols="`
	strDivOpCl = []byte("\">\n")       // "\">\n"
	strDivCl   = []byte("</div>") // "</textarea>"
)

// renderTextArea renders the component as an textarea HTML tag.
func (c *editorImpl) Render(w Writer) {
  w.Write(strDivOp)
  c.renderAttrsAndStyle(w)
  c.renderEnabled(w)
  c.renderEHandlers(w)
  w.Write(strDivOpCl)

  w.Writes(c.text)
  w.Write(strDivCl)

  w.Write([]byte(fmt.Sprintf(`<script>
    /**
    * This code is based on <https://github.com/pourquoi/ckeditor5-simple-upload>
    * and will be implemented by <https://github.com/mecha-cms/extend.c-k-editor> in the future!
    */

  // The upload adapter
  var Adapter = function(loader, urlOrObject, t) {

    var $ = this;

    $.loader = loader;
    $.urlOrObject = urlOrObject;
    $.t = t;

    $.upload = function() {
        return new Promise(function(resolve, reject) {
            $._initRequest();
            $._initListeners(resolve, reject);
            $._sendRequest();
        });
    };

    $.abort = function() {
        $.xhr && $.xhr.abort();
    };

    $._initRequest = function() {
        $.xhr = new XMLHttpRequest();
        var url = $.urlOrObject,
            headers;
        if (typeof url === "object") {
            url = url.url;
            headers = url.headers;
        }
        $.xhr.withCredentials = true;
        $.xhr.open('POST', url, true);
        if (headers) {
            for (var key in headers) {
                if (typeof headers[key] === "function") {
                    $.xhr.setRequestHeader(key, headers[key]());
                } else {
                    $.xhr.setRequestHeader(key, headers[key]);
                }
            }
        }
        $.xhr.responseType = 'json';
    };

    $._initListeners = function(resolve, reject) {
        var xhr = $.xhr,
            loader = $.loader,
            t = $.t,
            genericError = t('Cannot upload file:') + ' ' + loader.file.name;
        xhr.addEventListener('error', function() {
            reject(genericError);
        });
        xhr.addEventListener('abort', reject);
        xhr.addEventListener('load', function() {
            var response = xhr.response;
            if (!response || !response.uploaded) {
                return reject(response && response.error && response.error.message ? response.error.message : genericError);
            }
            resolve({
                'default': response.url
            });
        });
        if (xhr.upload) {
            xhr.upload.addEventListener('progress', function(evt) {
                if (evt.lengthComputable) {
                    loader.uploadTotal = evt.total;
                    loader.uploaded = evt.loaded;
                }
            });
        }
    }

    $._sendRequest = function() {
        var data = new FormData();
        data.append('upload', $.loader.file);
        $.xhr.send(data);
    };

  };

    let editor%d;
    InlineEditor
    .create( document.getElementById( '%d' ) )
    .then( newEditor => {
      editor%d = newEditor;
      editor%d.model.document.on( 'change:data', () => { 
        value = editor%d.getData()
        //console.log("The data has changed!" + value ); 
        se(null,11,%d,encodeURIComponent(value))
      } );
      editor%d.plugins.get('FileRepository').createUploadAdapter = function(loader) {
        return new Adapter(loader, '/upload.php?token=b4d455', editor%d.t);
      };
    })
    .catch( error => {
      console.error( error );
    });
    </script>`, c.id, c.id, c.id, c.id, c.id, c.id, c.id, c.id)))
}
