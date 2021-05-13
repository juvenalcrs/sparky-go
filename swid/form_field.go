package swid

import "fyne.io/fyne/v2"

// FormField defines a widget that can be used inside a Form.
type FormField interface {
	fyne.Widget
	Reset()
	Save()
	Validate() error

	setParentForm(f *Form)
	didChange()
	reset()
}

// BaseFormField defines a base form field.
type BaseFormField struct {
	form *Form
}

func (b *BaseFormField) setParentForm(f *Form) {
	b.form = f
}

func (b *BaseFormField) didChange() {
	if b.form == nil {
		return
	}
	b.form.fieldDidChange()
}

func (b *BaseFormField) reset() {
	if b.form == nil {
		return
	}
	b.form.fieldDidChange()
}
