package sparky

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LoaderTitles defines the title text on the specified loader state.
type LoaderTitles struct {
	Loading string
	Done    string
	Error   string
}

// Loader defines sparky loader.
type Loader struct {
	popup   *widget.PopUp
	content *loaderContent
}

// newLoader creates a new sparky loader.
func newLoader(win fyne.Window, message string, minWidth float32, titles *LoaderTitles) *Loader {
	l := &Loader{}
	l.content = newLoaderContent(message)
	l.content.loadingTitle = titles.Loading
	l.content.doneTitle = titles.Done
	l.content.errorTitle = titles.Error
	l.popup = widget.NewModalPopUp(l.content, win.Canvas())
	l.popup.Show()
	// this fixes the initial big min height at start because of the label
	// text wrapping, so given it the disired width, solve this problem
	l.content.Resize(fyne.NewSize(minWidth, 0))
	l.popup.Resize(fyne.NewSize(minWidth, l.content.MinSize().Height))
	return l
}

// UpdateMessage updates loader message.
func (l *Loader) UpdateMessage(s string) {
	l.content.message = s
	l.content.Refresh()
}

// Done stops the loader with a done message.
func (l *Loader) Done(s string) <-chan interface{} {
	done := make(chan interface{})
	l.content.SetDone(s, func() {
		l.Hide()
		close(done)
	})
	return done
}

// Error stops the loader with an error message.
func (l *Loader) Error(s string) <-chan interface{} {
	done := make(chan interface{})
	l.content.SetError(s, func() {
		l.Hide()
		close(done)
	})
	return done
}

// Hide hides the loader.
func (l *Loader) Hide() {
	if l.popup != nil {
		l.popup.Hide()
	}
}

// ===============================================================
// Content
// ===============================================================

type loaderState int

const (
	loaderStateLoading loaderState = iota
	loaderStateDone
	loaderStateError
)

type loaderContent struct {
	widget.BaseWidget
	message      string
	loadingTitle string
	doneTitle    string
	errorTitle   string
	onTappedOk   func()

	state loaderState
}

func newLoaderContent(message string) *loaderContent {
	l := &loaderContent{}
	l.ExtendBaseWidget(l)
	l.message = message
	l.state = loaderStateLoading
	return l
}

func (l *loaderContent) SetDone(message string, onTappedOk func()) {
	l.message = message
	l.state = loaderStateDone
	l.onTappedOk = onTappedOk
	l.Refresh()
}

func (l *loaderContent) SetError(message string, onTappedOk func()) {
	l.message = message
	l.state = loaderStateError
	l.onTappedOk = onTappedOk
	l.Refresh()
}

func (l *loaderContent) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	bgIcon := &canvas.Image{}
	bg := canvas.NewRectangle(color.Black)
	title := canvas.NewText("", theme.ForegroundColor())
	title.Alignment = fyne.TextAlignLeading
	title.TextStyle.Bold = true
	message := widget.NewLabelWithStyle(l.message, fyne.TextAlignCenter, fyne.TextStyle{})
	message.Wrapping = fyne.TextWrapWord
	okButton := widget.NewButton("Ok", l.onTappedOk)
	okButton.Hide()
	progressIndicator := widget.NewProgressBarInfinite()
	progressIndicator.Hide()

	r := &loaderContentRenderer{
		widget:            l,
		bgIcon:            bgIcon,
		bg:                bg,
		title:             title,
		message:           message,
		okButton:          okButton,
		progressIndicator: progressIndicator,
		objects: []fyne.CanvasObject{
			bgIcon, bg, title, message, okButton, progressIndicator,
		},
	}
	r.Refresh() // REVIEW is this needed?
	return r
}

type loaderContentRenderer struct {
	bgIcon            *canvas.Image
	bg                *canvas.Rectangle
	title             *canvas.Text
	message           *widget.Label
	okButton          *widget.Button
	progressIndicator fyne.CanvasObject

	objects []fyne.CanvasObject
	widget  *loaderContent
}

func (r *loaderContentRenderer) Destroy() {}

func (r *loaderContentRenderer) Layout(size fyne.Size) {
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

	// okButton
	okButtonMinHeight := r.okButton.MinSize().Height
	r.okButton.Move(fyne.NewPos(insetPad, ypos))
	r.okButton.Resize(fyne.NewSize(contentWidth, okButtonMinHeight))

	// progressIndicator (it would be at the same place of okButton)
	pIndicatorMinHeight := r.progressIndicator.MinSize().Height
	r.progressIndicator.Move(fyne.NewPos(insetPad, ypos))
	r.progressIndicator.Resize(fyne.NewSize(contentWidth, pIndicatorMinHeight))
}

func (r *loaderContentRenderer) MinSize() fyne.Size {
	insetPad := dialogInsetPad()
	pad := theme.Padding()

	tmin := r.title.MinSize()
	mmin := r.message.MinSize()
	bmin := r.okButton.MinSize()
	imin := r.progressIndicator.MinSize()

	min := fyne.NewSize(0, 0)
	min.Height = insetPad + tmin.Height + pad
	min.Height += mmin.Height + pad
	if r.widget.state == loaderStateLoading {
		min.Width = MaxFloat32(tmin.Width, mmin.Width, imin.Width) + 2*insetPad
		min.Height += imin.Height + insetPad
	} else {
		min.Width = MaxFloat32(tmin.Width, mmin.Width, bmin.Width) + 2*insetPad
		min.Height += bmin.Height + insetPad
	}
	return min
}

func (r *loaderContentRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *loaderContentRenderer) Refresh() {
	r.bg.FillColor = dialogBackgroundColor()
	r.message.SetText(r.widget.message)
	switch r.widget.state {
	case loaderStateLoading:
		r.title.Text = r.widget.loadingTitle
		r.title.TextSize = dialogTitleSize() - 2
		r.title.Color = theme.ForegroundColor()
		r.progressIndicator.Show()
		r.bgIcon.Hide()
		r.okButton.Hide()
	case loaderStateDone:
		r.title.Text = r.widget.doneTitle
		r.title.TextSize = dialogTitleSize()
		r.title.Color = successColor()
		r.bgIcon.Resource = theme.ConfirmIcon()
		r.bgIcon.Show()
		r.okButton.OnTapped = r.widget.onTappedOk
		r.okButton.Hidden = false
		r.okButton.Refresh()
		r.progressIndicator.Hide()
	case loaderStateError:
		r.title.Text = r.widget.errorTitle
		r.title.TextSize = dialogTitleSize()
		r.title.Color = theme.ErrorColor()
		r.bgIcon.Resource = theme.ErrorIcon()
		r.bgIcon.Show()
		r.okButton.OnTapped = r.widget.onTappedOk
		r.okButton.Hidden = false
		r.okButton.Refresh()
		r.progressIndicator.Hide()
	}
	r.bg.Refresh()
	r.title.Refresh()
}
