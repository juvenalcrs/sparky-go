package swid

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SelectFormField defines a special select field for Forms.
type SelectFormField struct {
	BaseFormField

	Options     []string
	Placeholder string
	Validator   fyne.StringValidator

	OnChanged func(string) `json:"-"`
	OnSaved   func(s string)

	selectField  *SelectField
	initialValue string
	isRendered   bool // TODO remove when Fyne has a way to check if the widget has been renderered or not
}

// NewSelectFormField creates a new select form field.
func NewSelectFormField(label, initialValue string, options []string) *SelectFormField {
	s := &SelectFormField{}
	s.ExtendBaseFormField(s)
	s.Label = label
	s.Options = options
	s.initialValue = initialValue
	s.setupSelectField()
	return s
}

// ===============================================================
// Constructor shortcuts
// ===============================================================

// WithValidator adds a validator to the text form field. This function does not
// call Refresh() and its use is only to define validators when creating the widget
// in the same line.
// TODO wait for fyne
// func (s *SelectFormField) WithValidator(v fyne.StringValidator) *SelectFormField {
// 	s.Validator = v
// 	return s
// }

// WithOnSaved adds OnSaved callback to the text form field. This function does not
// call Refresh() and its use is only to add OnSaved callback when creating the widget
// in the same line.
func (s *SelectFormField) WithOnSaved(onSaved func(s string)) *SelectFormField {
	s.OnSaved = onSaved
	return s
}

// ===============================================================
// Methods
// ===============================================================

// Selected returns the selected value.
func (s *SelectFormField) Selected() string {
	return s.selectField.Selected
}

// SetSelected sets the current option.
func (s *SelectFormField) SetSelected(text string) {
	s.selectField.Selected = text
	s.Refresh() // refresh the whole widget
}

// Reset resets the text value to the initial value.
func (s *SelectFormField) Reset() {
	s.dirty = false
	s.SetSelected(s.initialValue)
	s.didChange()
}

// Save triggers the OnSaved callback.
func (s *SelectFormField) Save() {
	if s.OnSaved != nil {
		s.OnSaved(s.selectField.Selected)
	}
}

// ValidationError returns the underlying validation error.
func (s *SelectFormField) ValidationError() error {
	if s.Validator != nil {
		// TODO remove when Fyne has a way to check if the widget has been renderered or not
		// means that this was called before CreateRenderer so create it by refreshing.
		if !s.isRendered {
			s.ExtendBaseFormField(s)
			s.Refresh()
		}
		return s.validationError
	}
	return nil
}

// Validate validates the field.
func (s *SelectFormField) Validate() error {
	if s.Validator != nil {
		s.ExtendBaseFormField(s)
		err := s.Validator(s.selectField.Selected)
		if s.validationError != err {
			s.validationError = err
			s.Refresh()
		}
		return s.validationError
	}
	return nil
}

func (s *SelectFormField) setupSelectField() {
	s.selectField = NewSelectField(s.Options, nil)
	s.selectField.Selected = s.initialValue
	s.selectField.OnChanged = func(text string) {
		if s.Validator != nil {
			if err := s.Validator(text); s.validationError != err {
				s.validationError = err
				s.Refresh()
			}
		}
		if s.OnChanged != nil {
			s.OnChanged(text)
		}
		s.didChange()
		if text == "" || !s.selectField.focused {
			s.Refresh()
		}
	}
	s.selectField.onFocusChanged = func(bool) {
		s.Refresh()
	}
	s.selectField.onHoverChanged = func(bool) {
		s.Refresh()
	}
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer implements fyne.Widget.
func (s *SelectFormField) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseFormField(s)

	if s.Validator != nil {
		s.validationError = s.Validator(s.selectField.Selected)
	}

	isFieldEmpty := func() bool {
		return false // Fyne force always a widget.Select placeholder so make this field unempty.
	}

	isFieldFocused := func() bool {
		return s.selectField.focused
	}

	updateInternalField := func() {
		s.selectField.Options = s.Options
		s.selectField.PlaceHolder = s.Placeholder
		if s.Disabled() {
			s.selectField.Disable()
		} else {
			s.selectField.Enable()
		}
		s.selectField.Refresh()
	}

	r := s.CreateBaseRenderer(
		s.Label, s.Hint, s.selectField,
		isFieldEmpty, isFieldFocused,
		updateInternalField,
	)

	r.(*formFieldRenderer).labelBgColor = func() color.Color {
		if s.selectField.hovered && !s.selectField.Disabled() {
			return theme.HoverColor()
		}
		return theme.InputBackgroundColor()
	}

	s.isRendered = true // TODO remove when Fyne has a way to check if the widget has been renderered or not

	return r
}
