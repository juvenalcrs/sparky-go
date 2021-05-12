package sparky

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ValueKey defines key type for injected values.
type ValueKey int

// DialogStyle define context dialog style.
type DialogStyle struct {
	MinWidth     float32
	LoaderTitles LoaderTitles
}

// Context defines a sparky context.
type Context interface {
	// Child creates a child context for a child window. The new context
	// will inherit all the parent context values.
	Child(win fyne.Window) Context
	// Window returns the underlying window.
	Window() fyne.Window
	// PutValue puts a value into the context, that later can be retrieved
	// with GetValue.
	PutValue(key ValueKey, v interface{})
	// GetValue gets a value previously added with PutValue.
	GetValue(key ValueKey) interface{}
	// ShowLoader shows a loader dialog.
	ShowLoader(message string) *Loader
	// ShowModal shows a modal with the specified content.
	ShowModal(content fyne.CanvasObject) *widget.PopUp
	// ShowInfo shows an information alert.
	ShowInfo(title, message string)
	// ShowSuccess shows a success alert.
	ShowSuccess(title, message string)
	// ShowError shows an error alert.
	ShowError(title, message string)
	// ShowConfirm shows a confirm alert. This will return a boolean
	// channel that will have the confirmation response.
	ShowConfirm(title, message, confirmButtonText string) <-chan bool
	// ShowTextInput shows a text input dialog. It will return a nil string
	// if the user cancels the dialog.
	ShowTextInput(title, message, submitText string) <-chan *string
	// ShowPasswordInput shows a password input dialog. It will return a nil string
	// if the user cancels the dialog.
	ShowPasswordInput(title, message, submitText string) <-chan *string
}

// NewContext creates a new sparky context.
func NewContext(win fyne.Window) Context {
	return &contextImpl{
		win:    win,
		values: map[ValueKey]interface{}{},
		dialogStyle: &DialogStyle{
			MinWidth: 300,
			LoaderTitles: LoaderTitles{
				Loading: "Processing!",
				Done:    "Done!",
				Error:   "Error!",
			},
		},
	}
}

// NewContextWithStyle creates a new sparky context with style options.
func NewContextWithStyle(win fyne.Window, s *DialogStyle) Context {
	return &contextImpl{
		win:         win,
		values:      map[ValueKey]interface{}{},
		dialogStyle: &(*s),
	}
}

// ===============================================================
// Implementation
// ===============================================================

type contextImpl struct {
	win         fyne.Window
	values      map[ValueKey]interface{}
	dialogStyle *DialogStyle
}

func (c *contextImpl) Child(win fyne.Window) Context {
	return &contextImpl{
		win:         win,
		values:      c.values,
		dialogStyle: c.dialogStyle,
	}
}

func (c *contextImpl) Window() fyne.Window {
	return c.win
}

func (c *contextImpl) PutValue(key ValueKey, v interface{}) {
	c.values[key] = v
}

func (c *contextImpl) GetValue(key ValueKey) interface{} {
	return c.values[key]
}

func (c *contextImpl) ShowLoader(message string) *Loader {
	return newLoader(c.win, message, c.dialogStyle.MinWidth, &c.dialogStyle.LoaderTitles)
}

func (c *contextImpl) ShowModal(content fyne.CanvasObject) *widget.PopUp {
	m := widget.NewModalPopUp(content, c.win.Canvas())
	m.Show()
	return m
}

func (c *contextImpl) ShowConfirm(title, message, confirmButtonText string) <-chan bool {
	resp := make(chan bool, 1)
	alert := newAlertBase(alertTypeConfirm, title, message)
	alert.okBtnText = confirmButtonText
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOk = func() {
		d.Hide()
		resp <- true
		close(resp)
	}
	alert.onTappedCancel = func() {
		d.Hide()
		close(resp)
	}
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
	return resp
}

func (c *contextImpl) ShowInfo(title, message string) {
	alert := newAlertBase(alertTypeInfo, title, message)
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOk = func() { d.Hide() }
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
}

func (c *contextImpl) ShowError(title, message string) {
	alert := newAlertBase(alertTypeError, title, message)
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOk = func() { d.Hide() }
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
}

func (c *contextImpl) ShowSuccess(title, message string) {
	alert := newAlertBase(alertTypeSuccess, title, message)
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOk = func() { d.Hide() }
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
}

func (c *contextImpl) ShowTextInput(title, message, submitText string) <-chan *string {
	resp := make(chan *string, 1)
	alert := newAlertBase(alertTypeInput, title, message)
	alert.okBtnText = submitText
	alert.isPasswordInput = false
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOkInput = func(s string) {
		d.Hide()
		resp <- &s
		close(resp)
	}
	alert.onTappedCancel = func() {
		d.Hide()
		close(resp)
	}
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
	return resp
}

func (c *contextImpl) ShowPasswordInput(title, message, submitText string) <-chan *string {
	resp := make(chan *string, 1)
	alert := newAlertBase(alertTypeInput, title, message)
	alert.okBtnText = submitText
	alert.isPasswordInput = true
	d := widget.NewModalPopUp(alert, c.win.Canvas())
	alert.onTappedOkInput = func(s string) {
		d.Hide()
		resp <- &s
		close(resp)
	}
	alert.onTappedCancel = func() {
		d.Hide()
		close(resp)
	}
	d.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	alert.Resize(fyne.NewSize(c.dialogStyle.MinWidth, 0))
	d.Resize(fyne.NewSize(c.dialogStyle.MinWidth, alert.MinSize().Height))
	return resp
}
