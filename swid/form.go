package swid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Form defines form widget.
type Form struct {
	widget.BaseWidget
	OnChanged           func()
	OnValidationChanged func(valid bool)
	cols                int
	fields              []FormField
	isValid             bool
	submitButton        *widget.Button
}

// NewForm creates a new form widget.
//
func NewForm(cols int, fields ...FormField) *Form {
	f := &Form{cols: cols, fields: fields, isValid: true}
	f.ExtendBaseWidget(f)
	return f
}

// IsValid returns true if the form is valid.
func (f *Form) IsValid() bool {
	f.validate()
	return f.isValid
}

// Reset resets or clears all the inputs.
//
func (f *Form) Reset() {
	for _, field := range f.fields {
		field.Reset()
	}
}

// Save triggers onSaved callback of all FormFields.
//
func (f *Form) Save() {
	for _, field := range f.fields {
		field.Save()
	}
}

// SubmitButton returns the form submit button.
func (f *Form) SubmitButton(text string, onTapped func()) *widget.Button {
	if f.submitButton != nil {
		f.submitButton.Text = text
		f.submitButton.OnTapped = onTapped
		f.submitButton.Refresh()
		return f.submitButton
	}
	f.submitButton = widget.NewButton(text, onTapped)
	f.submitButton.Importance = widget.HighImportance
	return f.submitButton
}

// ResetButton returns the form reset button.
func (f *Form) ResetButton(text string) *widget.Button {
	return widget.NewButton(text, f.Reset)
}

// updates submit button state if there is one.
func (f *Form) updateSubmitButtonState() {
	if f.submitButton == nil {
		return
	}
	if f.isValid {
		f.submitButton.Enable()
	} else {
		f.submitButton.Disable()
	}
}

// Validate validates the form. If it is invalid, it will return
// the first error found.
func (f *Form) validate() {
	isValid := true
	for _, field := range f.fields {
		// use only validationError because the validation is done
		// automatically by the form fields itself.
		if err := field.ValidationError(); err != nil && isValid {
			isValid = false
			// do not return here, to ensure we validate all fields
		}
	}
	f.isValid = isValid
}

// fieldDidChange must be called from a form field.
func (f *Form) fieldDidChange() {
	if f.OnChanged != nil {
		f.OnChanged()
	}
	prev := f.isValid
	f.validate()
	if prev == f.isValid {
		return
	}
	f.updateSubmitButtonState()
	if f.OnValidationChanged != nil {
		f.OnValidationChanged(f.isValid)
	}
}

// ===============================================================
// FormRenderer
// ===============================================================

// CreateRenderer implements fyne.WidgetRenderer.
//
func (f *Form) CreateRenderer() fyne.WidgetRenderer {
	f.ExtendBaseWidget(f)
	objects := make([]fyne.CanvasObject, len(f.fields))
	f.isValid = true
	for i, field := range f.fields {
		field.setParentForm(f)
		if err := field.Validate(); err != nil && f.isValid {
			f.isValid = false
		}
		objects[i] = field
	}
	f.updateSubmitButtonState()
	if f.OnValidationChanged != nil {
		f.OnValidationChanged(f.isValid)
	}
	return &formRenderer{
		layout:  layout.NewGridLayoutWithColumns(f.cols),
		objects: objects,
	}
}

type formRenderer struct {
	layout  fyne.Layout
	objects []fyne.CanvasObject
}

func (r *formRenderer) Destroy() {}

func (r *formRenderer) Layout(size fyne.Size) {
	r.layout.Layout(r.objects, size)
}

func (r *formRenderer) MinSize() fyne.Size {
	return r.layout.MinSize(r.objects)
}

func (r *formRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *formRenderer) Refresh() {}
