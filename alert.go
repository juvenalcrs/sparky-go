package sparky

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type alertType int

const (
	alertTypeInfo alertType = iota
	alertTypeSuccess
	alertTypeConfirm
	alertTypeInput
	alertTypeError
)

type alertContent struct {
	widget.BaseWidget
	alertType      alertType
	title          string
	message        string
	okBtnText      string
	onTappedOk     func() // for not-input alerts
	cancelBtnText  string
	onTappedCancel func()

	// for input alerts
	onTappedOkInput func(s string)
	isPasswordInput bool
}

func newAlertBase(alertType alertType, title, message string) *alertContent {
	a := &alertContent{}
	a.ExtendBaseWidget(a)
	a.alertType = alertType
	a.title = title
	a.message = message
	a.okBtnText = "Ok"
	a.cancelBtnText = "Cancel"
	return a
}

func (a *alertContent) hasTwoButtons() bool {
	return a.alertType == alertTypeConfirm || a.alertType == alertTypeInput
}

func (a *alertContent) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	bgIcon := &canvas.Image{}
	bg := canvas.NewRectangle(color.Black)
	title := canvas.NewText(a.title, theme.ForegroundColor())
	title.Alignment = fyne.TextAlignLeading
	title.TextStyle.Bold = true
	message := widget.NewLabelWithStyle(a.message, fyne.TextAlignCenter, fyne.TextStyle{})
	message.Wrapping = fyne.TextWrapWord
	okButton := widget.NewButton(a.okBtnText, a.onTappedOk)
	okButton.Hide()
	cancelButton := widget.NewButton(a.cancelBtnText, a.onTappedCancel)
	cancelButton.Hide()

	objects := []fyne.CanvasObject{
		bgIcon, bg, title, message, okButton, cancelButton,
	}

	var input *widget.Entry
	if a.alertType == alertTypeInput {
		if a.isPasswordInput {
			input = widget.NewPasswordEntry()
		} else {
			input = widget.NewEntry()
		}
		okButton.OnTapped = func() { a.onTappedOkInput(input.Text) }
		objects = append(objects, input)
	}

	r := &alertContentRenderer{
		bgIcon:       bgIcon,
		bg:           bg,
		title:        title,
		message:      message,
		input:        input,
		okButton:     okButton,
		cancelButton: cancelButton,
		widget:       a,
		objects:      objects,
	}
	r.Refresh()
	return r
}

type alertContentRenderer struct {
	bgIcon       *canvas.Image
	bg           *canvas.Rectangle
	title        *canvas.Text
	message      *widget.Label
	input        *widget.Entry
	okButton     *widget.Button
	cancelButton *widget.Button

	widget  *alertContent
	objects []fyne.CanvasObject
}

func (r *alertContentRenderer) Destroy() {}

func (r *alertContentRenderer) Layout(size fyne.Size) {
	insetPad := dialogInsetPad()
	pad := theme.Padding()

	// background
	r.bg.Move(fyne.NewPos(0, 0))
	r.bg.Resize(size)

	// bgIcon
	iconSize := float32(80)
	r.bgIcon.Resize(fyne.NewSize(iconSize, iconSize))
	r.bgIcon.Move(fyne.NewPos(size.Width-iconSize+pad, -pad))

	contentWidth := size.Width - 2*insetPad

	// title
	titleMinHeight := r.title.MinSize().Height
	r.title.Move(fyne.NewPos(insetPad, insetPad))
	r.title.Resize(fyne.NewSize(contentWidth, titleMinHeight))
	ypos := insetPad + titleMinHeight + pad

	// message
	messageMinHeight := r.message.MinSize().Height
	r.message.Move(fyne.NewPos(insetPad, ypos))
	r.message.Resize(fyne.NewSize(contentWidth, messageMinHeight))
	ypos += messageMinHeight + pad

	if r.input != nil {
		inputMinHeight := r.input.MinSize().Height
		r.input.Move(fyne.NewPos(insetPad, ypos))
		r.input.Resize(fyne.NewSize(contentWidth, inputMinHeight))
		ypos += inputMinHeight + 4*pad
	}

	if r.widget.hasTwoButtons() {
		// two buttons
		btnWidth := (contentWidth - pad) / 2
		btnHeight := MaxFloat32(r.okButton.MinSize().Height, r.cancelButton.MinSize().Height)
		r.cancelButton.Move(fyne.NewPos(insetPad, ypos))
		r.cancelButton.Resize(fyne.NewSize(btnWidth, btnHeight))
		r.okButton.Move(fyne.NewPos(insetPad+btnWidth+pad, ypos))
		r.okButton.Resize(fyne.NewSize(btnWidth, btnHeight))
	} else {
		// okButton only
		okBtnMinHeight := r.okButton.MinSize().Height
		r.okButton.Move(fyne.NewPos(insetPad, ypos))
		r.okButton.Resize(fyne.NewSize(contentWidth, okBtnMinHeight))
	}
}

func (r *alertContentRenderer) MinSize() fyne.Size {
	insetPad := dialogInsetPad()
	pad := theme.Padding()

	tmin := r.title.MinSize()
	mmin := r.message.MinSize()
	omin := r.okButton.MinSize()
	cmin := r.cancelButton.MinSize()

	min := fyne.NewSize(0, 0)
	min.Height = insetPad + tmin.Height + pad
	min.Height += mmin.Height + pad
	if r.input != nil {
		inputMin := r.input.MinSize()
		min.Width = inputMin.Width
		min.Height += inputMin.Height + 4*pad
	}
	if r.widget.hasTwoButtons() {
		min.Width = MaxFloat32(min.Width, tmin.Width, mmin.Width, omin.Width+pad+cmin.Width) + 2*insetPad
		min.Height += MaxFloat32(omin.Height, cmin.Height) + insetPad
	} else {
		min.Width = MaxFloat32(tmin.Width, mmin.Width, omin.Width) + 2*insetPad
		min.Height += omin.Height + insetPad
	}
	return min
}

func (r *alertContentRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *alertContentRenderer) Refresh() {
	r.bg.FillColor = dialogBackgroundColor()
	r.title.Text = r.widget.title
	r.title.TextSize = dialogTitleSize()
	r.okButton.Text = r.widget.okBtnText
	r.cancelButton.Text = r.widget.cancelBtnText
	r.message.SetText(r.widget.message)
	switch r.widget.alertType {
	case alertTypeInfo:
		r.bgIcon.Resource = theme.InfoIcon()
		r.title.Color = infoColor()
		r.okButton.Importance = widget.MediumImportance
		r.okButton.Show()
		r.cancelButton.Hide()
	case alertTypeSuccess:
		r.bgIcon.Resource = theme.ConfirmIcon()
		r.title.Color = successColor()
		r.okButton.Importance = widget.MediumImportance
		r.okButton.Show()
		r.cancelButton.Hide()
	case alertTypeError:
		r.bgIcon.Resource = theme.ErrorIcon()
		r.title.Color = theme.ErrorColor()
		r.okButton.Importance = widget.MediumImportance
		r.okButton.Show()
		r.cancelButton.Hide()
	case alertTypeConfirm:
		r.bgIcon.Resource = theme.QuestionIcon()
		r.title.Color = theme.ForegroundColor()
		r.okButton.Importance = widget.HighImportance
		r.okButton.Show()
		r.cancelButton.Importance = widget.MediumImportance
		r.cancelButton.Show()
	case alertTypeInput:
		r.bgIcon.Resource = theme.QuestionIcon()
		r.title.Color = theme.ForegroundColor()
		r.okButton.Importance = widget.HighImportance
		r.okButton.Show()
		r.cancelButton.Importance = widget.MediumImportance
		r.cancelButton.Show()
	}
	r.bg.Refresh()
	r.title.Refresh()
	r.bgIcon.Refresh()
}
