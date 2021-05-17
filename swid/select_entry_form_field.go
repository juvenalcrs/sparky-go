package swid

import (
	"fyne.io/fyne/v2"
)

// SelectEntryFormField defines a special select entry field for Forms.
type SelectEntryFormField struct {
	BaseFormField

	TextStyle   fyne.TextStyle
	Placeholder string
	Wrapping    fyne.TextWrap
	Validator   fyne.StringValidator

	OnChanged func(string) `json:"-"`
	OnSaved   func(s string)

	selectEntryField *SelectEntryField
	initialText      string
}

// NewSelectEntryFormField creates a new select entry form field.
func NewSelectEntryFormField(label, initialValue string, options []string) *SelectEntryFormField {
	s := &SelectEntryFormField{}
	s.ExtendBaseFormField(s)
	s.Label = label
	s.Wrapping = fyne.TextTruncate
	s.initialText = initialValue
	s.setupSelectEntryField(options)
	return s
}

// ===============================================================
// Constructor shortcuts
// ===============================================================

// WithValidator adds a validator to the text form field. This function does not
// call Refresh() and its use is only to define validators when creating the widget
// in the same line.
func (s *SelectEntryFormField) WithValidator(v fyne.StringValidator) *SelectEntryFormField {
	s.Validator = v
	return s
}

// WithOnSaved adds OnSaved callback to the text form field. This function does not
// call Refresh() and its use is only to add OnSaved callback when creating the widget
// in the same line.
func (s *SelectEntryFormField) WithOnSaved(onSaved func(s string)) *SelectEntryFormField {
	s.OnSaved = onSaved
	return s
}

// ===============================================================
// Methods
// ===============================================================

// Text returns the current text value.
func (s *SelectEntryFormField) Text() string {
	return s.selectEntryField.Text
}

// SetText manually sets the text of the TextFormField to the given text value.
func (s *SelectEntryFormField) SetText(text string) {
	s.selectEntryField.Text = text
	s.Refresh() // refresh the whole widget
}

// SetOptions sets the options the user might select from.
func (s *SelectEntryFormField) SetOptions(options []string) {
	s.selectEntryField.SetOptions(options)
}

// Reset resets the text value to the initial value.
func (s *SelectEntryFormField) Reset() {
	s.dirty = false
	s.SetText(s.initialText)
	s.didChange()
}

// Save triggers the OnSaved callback.
func (s *SelectEntryFormField) Save() {
	if s.OnSaved != nil {
		s.OnSaved(s.selectEntryField.Text)
	}
}

// Validate validates the field.
func (s *SelectEntryFormField) Validate() error {
	if s.Validator != nil {
		// means that this was called before CreateRenderer and
		// then Validator field is not copy to the textField yet,
		// so Refresh to generate it
		if s.selectEntryField.Validator == nil {
			s.ExtendBaseFormField(s)
			s.Refresh()
		}
		return s.validationError
	}
	return nil
}

func (s *SelectEntryFormField) setupSelectEntryField(options []string) {
	s.selectEntryField = NewSelectEntryField(options)
	s.selectEntryField.Text = s.initialText
	s.selectEntryField.OnChanged = func(text string) {
		if s.OnChanged != nil {
			s.OnChanged(text)
		}
		s.didChange()
		if text == "" || !s.selectEntryField.focused {
			s.Refresh()
		}
	}
	s.selectEntryField.onFocusChanged = func(bool) {
		s.Refresh()
	}
	s.selectEntryField.SetOnValidationChanged(func(e error) {
		s.validationError = e
		s.Refresh()
	})
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer implements fyne.Widget.
func (s *SelectEntryFormField) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseFormField(s)

	s.selectEntryField.Validator = s.Validator
	s.selectEntryField.Validate() // validates as soon as it is created

	isFieldEmpty := func() bool {
		return s.selectEntryField.Text == ""
	}

	isFieldFocused := func() bool {
		return s.selectEntryField.focused
	}

	updateInternalField := func() {
		s.selectEntryField.TextStyle = s.TextStyle
		s.selectEntryField.Wrapping = s.Wrapping
		// TODO change SetPlaceholder by r.widget.selectEntryField.PlaceHolder when it is fixed in fyne
		if s.selectEntryField.focused && s.selectEntryField.Text == "" {
			s.selectEntryField.SetPlaceHolder(s.Placeholder)
		} else {
			s.selectEntryField.SetPlaceHolder("")
		}
		s.selectEntryField.Wrapping = s.Wrapping
		s.selectEntryField.Validator = s.Validator
		if s.Disabled() {
			s.selectEntryField.Disable()
		} else {
			s.selectEntryField.Enable()
		}
		s.selectEntryField.Refresh()
	}

	return s.CreateBaseRenderer(
		s.Label, s.Hint, s.selectEntryField,
		isFieldEmpty, isFieldFocused,
		updateInternalField,
	)
}
