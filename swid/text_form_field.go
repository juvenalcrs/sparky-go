package swid

import (
	"fyne.io/fyne/v2"
)

// TextFormField defines a special text field for Forms.
type TextFormField struct {
	BaseFormField

	TextStyle   fyne.TextStyle
	Placeholder string
	Wrapping    fyne.TextWrap
	Validator   fyne.StringValidator
	// ActionItem is a small item which is displayed at the outer right of the entry (like a password revealer)
	ActionItem fyne.CanvasObject
	MaxLength  int

	OnChanged func(s string)
	OnSaved   func(s string)

	textField       *TextField
	initialText     string
	isPasswordField bool
}

// NewTextFormField creates a new special text field for Forms.
func NewTextFormField(label, initialText string) *TextFormField {
	t := &TextFormField{}
	t.ExtendBaseFormField(t)
	t.Label = label
	t.Wrapping = fyne.TextTruncate
	t.initialText = initialText
	t.setupTextField()
	return t
}

// NewRestrictTextFormField creates a new text form field that accepts an input
// type.
func NewRestrictTextFormField(label, initialText string, input RestrictInput) *TextFormField {
	t := NewTextFormField(label, initialText)
	t.textField.restriction = input
	return t
}

// NewPasswordTextFormField creates a new password text field.
func NewPasswordTextFormField(label, initialText string) *TextFormField {
	t := NewTextFormField(label, initialText)
	t.isPasswordField = true
	t.textField.Password = true
	return t
}

// NewMaskedTextFormField creates a new text form field with a mask.
// Mask definitions:
//
//  9: Represents a numeric character (0-9)
//  a: Represents an alpha character (A-Z,a-z)
//  *: Represents an alphanumeric character (A-Z,a-z,0-9)
func NewMaskedTextFormField(label, initialText, mask, placeHolder string) *TextFormField {
	t := NewTextFormField(label, initialText)
	t.Placeholder = placeHolder
	t.textField.mask = []rune(mask)
	return t
}

// ===============================================================
// Constructor shortcuts
// ===============================================================

// WithValidator adds a validator to the text form field. This function does not
// call Refresh() and its use is only to define validators when creating the widget
// in the same line.
func (t *TextFormField) WithValidator(v fyne.StringValidator) *TextFormField {
	t.Validator = v
	return t
}

// WithOnSaved adds OnSaved callback to the text form field. This function does not
// call Refresh() and its use is only to add OnSaved callback when creating the widget
// in the same line.
func (t *TextFormField) WithOnSaved(onSaved func(s string)) *TextFormField {
	t.OnSaved = onSaved
	return t
}

// ===============================================================
// Methods
// ===============================================================

// Text returns the current text value.
func (t *TextFormField) Text() string {
	return t.textField.Text
}

// SetText manually sets the text of the TextFormField to the given text value.
func (t *TextFormField) SetText(text string) {
	t.textField.Text = text
	t.Refresh() // refresh the whole widget
}

// Reset resets the text value to the initial value.
func (t *TextFormField) Reset() {
	t.dirty = false
	t.SetText(t.initialText)
	t.didChange()
}

// Save triggers the OnSaved callback.
func (t *TextFormField) Save() {
	if t.OnSaved != nil {
		t.OnSaved(t.textField.Text)
	}
}

// Validate validates the field.
func (t *TextFormField) Validate() error {
	if t.Validator != nil {
		// means that this was called before CreateRenderer and
		// then Validator field is not copy to the textField yet,
		// so Refresh to generate it
		if t.textField.Validator == nil {
			t.ExtendBaseFormField(t)
			t.Refresh()
		}
		return t.validationError
	}
	return nil
}

func (t *TextFormField) setupTextField() {
	t.textField = NewTextField()
	t.textField.Text = t.initialText
	t.textField.OnChanged = func(s string) {
		if t.OnChanged != nil {
			t.OnChanged(s)
		}
		t.didChange()
		if s == "" {
			t.Refresh()
		}
	}
	t.textField.onFocusChanged = func(bool) {
		t.Refresh()
	}
	t.textField.SetOnValidationChanged(func(e error) {
		t.validationError = e
		t.Refresh()
	})
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer implements fyne.Widget.
func (t *TextFormField) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseFormField(t)

	if !t.isPasswordField {
		t.textField.ActionItem = t.ActionItem
	}
	t.textField.Validator = t.Validator
	t.textField.Validate() // validates as soon as it is created

	isFieldEmpty := func() bool {
		return t.textField.Text == ""
	}

	isFieldFocused := func() bool {
		return t.textField.focused
	}

	updateInternalField := func() {
		t.textField.TextStyle = t.TextStyle
		t.textField.Wrapping = t.Wrapping
		if !t.isPasswordField {
			// REVIEW this won't work until it is fixed in Fyne.
			t.textField.ActionItem = t.ActionItem
		}
		// TODO change SetPlaceholder by r.widget.textField.PlaceHolder when it is fixed in fyne
		if t.textField.focused && t.textField.Text == "" {
			t.textField.SetPlaceHolder(t.Placeholder)
		} else {
			t.textField.SetPlaceHolder("")
		}
		t.textField.Wrapping = t.Wrapping
		t.textField.MaxLength = t.MaxLength
		t.textField.Validator = t.Validator
		if t.Disabled() {
			t.textField.Disable()
		} else {
			t.textField.Enable()
		}
		t.textField.Refresh()
	}

	return t.CreateBaseRenderer(
		t.Label, t.Hint, t.textField,
		isFieldEmpty, isFieldFocused,
		updateInternalField,
	)
}
