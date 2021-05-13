package swid

import (
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TextField defines a basic editable text widget.
type TextField struct {
	widget.Entry
	MaxLength int

	mask        []rune
	restriction RestrictInput

	focused   bool
	onFocus   func()
	onUnfocus func()
}

// NewTextField creates a new text field.
func NewTextField() *TextField {
	t := &TextField{}
	t.ExtendBaseWidget(t)
	t.Wrapping = fyne.TextTruncate
	return t
}

// NewRestrictTextField creates a new text field that accepts an input
// type.
func NewRestrictTextField(input RestrictInput) *TextField {
	t := &TextField{}
	t.ExtendBaseWidget(t)
	t.restriction = input
	t.Wrapping = fyne.TextTruncate
	return t
}

// NewPasswordTextField creates a new password text field.
func NewPasswordTextField() *TextField {
	t := &TextField{}
	t.ExtendBaseWidget(t)
	t.Password = true
	t.Wrapping = fyne.TextTruncate
	return t
}

// NewMaskedTextField creates a new text field with a mask.
// Mask definitions:
//
//  9: Represents a numeric character (0-9)
//  a: Represents an alpha character (A-Z,a-z)
//  *: Represents an alphanumeric character (A-Z,a-z,0-9)
func NewMaskedTextField(mask, placeHolder string) *TextField {
	t := &TextField{}
	t.ExtendBaseWidget(t)
	t.Wrapping = fyne.TextTruncate
	t.PlaceHolder = placeHolder
	t.mask = []rune(mask)
	return t
}

// ===============================================================
// Implementation
// ===============================================================

// MinSize implements fyne.CanvasObject.
func (t *TextField) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// FocusGained overrides widget.Entry method.
func (t *TextField) FocusGained() {
	t.focused = true
	t.Entry.FocusGained()
	if t.onFocus != nil {
		t.onFocus()
	}
}

// FocusLost overrides widget.Entry method.
func (t *TextField) FocusLost() {
	t.focused = false
	t.Entry.FocusLost()
	if t.onUnfocus != nil {
		t.onUnfocus()
	}
}

// TypedRune overrides widget.Entry method.
func (t *TextField) TypedRune(r rune) {
	if t.Disabled() {
		return
	}
	if t.mask != nil {
		t.maskVerifyOnTypedRune(r, t.CursorColumn)
		return
	}
	if t.MaxLength > 0 && (len([]rune(t.Text))+1) > t.MaxLength {
		return
	}
	if acceptChar(t.restriction, t.Text, r, t.CursorColumn) {
		t.Entry.TypedRune(r)
	}
}

// maskVerifyOnTypedRune verifies mask when user type a rune.
func (t *TextField) maskVerifyOnTypedRune(r rune, colPos int) {
	totalLen := len([]rune(t.Text))
	if (totalLen + 1) > len(t.mask) {
		return
	}

	i := colPos
	for t.mask[i] != 'a' && t.mask[i] != '9' && t.mask[i] != '*' {
		t.Entry.TypedRune(t.mask[i])
		if t.mask[i] == r {
			return
		}
		i++
		totalLen++
		if totalLen >= len(t.mask) {
			return
		}
	}
	if (t.mask[i] == '9' && !unicode.IsDigit(r)) ||
		(t.mask[i] == 'a' && !unicode.IsLetter(r)) ||
		(t.mask[i] == '*' && !unicode.IsDigit(r) && !unicode.IsLetter(r)) {
		return
	}
	t.Entry.TypedRune(r)
}
