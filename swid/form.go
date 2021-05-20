package swid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// Form defines form widget.
type Form struct {
	widget.BaseWidget
	OnChanged           func()
	OnValidationChanged func(valid bool)

	container     *fyne.Container
	cols          int
	fields        []FormField
	isValid       bool
	submitButtons []*widget.Button
}

// NewForm creates a new form widget.
//
func NewForm(cols int, fields ...FormField) *Form {
	f := &Form{cols: cols, fields: fields, isValid: true}
	f.ExtendBaseWidget(f)
	return f
}

// NewCustomForm creates a new custom form from a container.
// Internally, it will extract the form fields from the container
// and attach them to this form.
func NewCustomForm(cont *fyne.Container) *Form {
	f := &Form{container: cont, fields: make([]FormField, 0, 20), isValid: true}
	f.fields = fieldsFromContent(f.fields, cont)
	f.ExtendBaseWidget(f)
	f.fields = f.fields[:len(f.fields):len(f.fields)]
	return f
}

// ===============================================================
// Methods
// ===============================================================

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

// CreateSubmitButton creates a new form submit button.
func (f *Form) CreateSubmitButton(text string, onTapped func()) *widget.Button {
	btn := widget.NewButton(text, onTapped)
	btn.Importance = widget.HighImportance
	f.submitButtons = append(f.submitButtons, btn)
	return btn
}

// CreateResetButton creates a new form reset button.
func (f *Form) CreateResetButton(text string) *widget.Button {
	return widget.NewButton(text, f.Reset)
}

// updates submit button state if there is one.
func (f *Form) updateSubmitButtonState() {
	isValid := f.isValid
	for _, btn := range f.submitButtons {
		if isValid {
			btn.Enable()
		} else {
			btn.Disable()
		}
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
	if f.container == nil {
		return &formRenderer{
			layout:  layout.NewGridLayoutWithColumns(f.cols),
			objects: objects,
		}
	}
	return &containerFormRenderer{widget: f}
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

type containerFormRenderer struct {
	widget *Form
}

func (r *containerFormRenderer) Destroy() {}

func (r *containerFormRenderer) Layout(size fyne.Size) {
	r.widget.container.Layout.Layout(r.widget.container.Objects, size)
}

func (r *containerFormRenderer) MinSize() fyne.Size {
	return r.widget.container.MinSize()
}

func (r *containerFormRenderer) Objects() []fyne.CanvasObject {
	return r.widget.container.Objects
}

func (r *containerFormRenderer) Refresh() {
	r.widget.container.Refresh()
}

// ===============================================================
// Private helpers
// ===============================================================

func fieldsFromContent(fields []FormField, content fyne.CanvasObject) []FormField {
	if content == nil {
		return fields
	}
	switch o := content.(type) {
	case fyne.Widget:
		for _, co := range test.WidgetRenderer(o).Objects() {
			fields = fieldsFromContent(fields, co)
		}
		if ff, ok := o.(FormField); ok {
			fields = append(fields, ff)
		}
	case *fyne.Container:
		for _, co := range o.Objects {
			fields = fieldsFromContent(fields, co)
		}
	}
	return fields
}
