package swid

import (
	"fyne.io/fyne/v2"
)

// TextFormField defines a special text field for Forms.
type TextFormField struct {
	BaseFormField

	TextStyle   fyne.TextStyle
	Placeholder string
	Password    bool
	Wrapping    fyne.TextWrap
	Validator   fyne.StringValidator
	// ActionItem is a small item which is displayed at the outer right of the entry (like a password revealer)
	ActionItem fyne.CanvasObject
	MaxLength  int

	OnChanged func(s string)
	OnSaved   func(s string)

	textField   *TextField
	initialText string
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
		t.textField.ActionItem = t.ActionItem
		t.textField.Wrapping = t.Wrapping
		t.textField.Password = t.Password
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
