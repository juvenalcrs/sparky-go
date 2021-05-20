package swid

import (
	"errors"
	"strconv"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fpabl0/sparky-go/svalid"
	"github.com/stretchr/testify/assert"
)

func TestCustomForm(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())

	title := widget.NewLabel("Title")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	name := NewTextFormField("Name", "")
	name.Validator = svalid.NotEmpty()
	name.Hint = "Your name"

	lastName := NewTextFormField("LastName", "")
	lastName.Validator = svalid.NotEmpty()
	lastName.Hint = "Your lastname"

	city := NewTextFormField("City", "Machala")
	city.Validator = svalid.NotEmpty()
	city.Hint = "Your city"

	accountType := NewSelectFormField("Account Type", "Basic", []string{"Basic", "Professional"})
	accountType.Hint = "Your account type"

	address1 := NewTextFormField("Address 1", "")
	address1.Validator = svalid.NotEmpty()
	address1.Hint = "Your address 1"

	address2 := NewTextFormField("Address 2", "")
	address2.Validator = svalid.NotEmpty()
	address2.Hint = "Your address 2"

	code := NewRestrictTextFormField("Code", "", RestrictInputInteger)
	code.Validator = func(s string) error {
		n, _ := strconv.Atoi(s)
		if n < 10 {
			return errors.New("invalid code")
		}
		return nil
	}
	code.Hint = "The sent code"

	favoriteColor := NewSelectEntryFormField("Favorite color", "Black", []string{
		"Yellow", "Blue",
	})
	favoriteColor.Validator = func(s string) error {
		if s == "white" {
			return errors.New("can't be white")
		}
		return nil
	}
	favoriteColor.Hint = "Your favorite color"

	f := NewCustomForm(container.NewVBox(
		title,
		container.NewGridWithColumns(2,
			name, lastName,
			city, accountType,
		),
		address1,
		address2,
		code,
		favoriteColor,
	))

	submitted := false
	submitButton := f.SubmitButton("Create", func() { submitted = true })
	rstButton := f.ResetButton("Reset")

	w := test.NewWindow(container.NewVBox(f, submitButton, rstButton))

	test.AssertImageMatches(t, "form/custom_form_initial.png", w.Canvas().Capture())

	name.SetText("Peter")
	lastName.SetText("Parker")
	city.SetText("Miami")
	accountType.SetSelected("Professional")
	address1.SetText("Av. Big one")
	address2.SetText("a") // just to make it dirty
	address2.SetText("")  // put empty to show the error
	code.SetText("1")
	favoriteColor.SetText("white")

	assert.False(t, f.IsValid())
	test.Tap(submitButton)
	assert.False(t, submitted)
	test.AssertImageMatches(t, "form/custom_form_invalid.png", w.Canvas().Capture())

	name.SetText("Jorge")
	lastName.SetText("Vélez")
	city.SetText("Orlando")
	accountType.SetSelected("Professional")
	address1.SetText("Long avenue")
	address2.SetText("Dept. 203")
	code.SetText("189")
	favoriteColor.SetText("Yellow")

	assert.True(t, f.IsValid())
	test.AssertImageMatches(t, "form/custom_form_valid.png", w.Canvas().Capture())

	assert.False(t, submitted)
	test.Tap(submitButton)
	assert.True(t, submitted)

	test.Tap(rstButton)
	test.AssertImageMatches(t, "form/custom_form_initial.png", w.Canvas().Capture())
}

func TestForm(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())

	name := NewTextFormField("Name", "")
	name.Validator = svalid.NotEmpty()
	name.Hint = "Your name"

	lastName := NewTextFormField("LastName", "Pérez")
	lastName.Validator = func(s string) error {
		if s == "last" {
			return errors.New("can't be last")
		}
		return nil
	}
	lastName.Hint = "Your lastname"

	f := NewForm(2, name, lastName)

	submitted := false
	submitButton := f.SubmitButton("Create", func() { submitted = true })
	rstButton := f.ResetButton("Reset")

	w := test.NewWindow(container.NewVBox(f, submitButton, rstButton))

	test.AssertImageMatches(t, "form/form_initial.png", w.Canvas().Capture())

	name.SetText("Peter")
	lastName.SetText("last")

	assert.False(t, f.IsValid())
	test.Tap(submitButton)
	assert.False(t, submitted)
	test.AssertImageMatches(t, "form/form_invalid.png", w.Canvas().Capture())

	name.SetText("Jorge")
	lastName.SetText("Vélez")

	assert.True(t, f.IsValid())
	test.AssertImageMatches(t, "form/form_valid.png", w.Canvas().Capture())

	assert.False(t, submitted)
	test.Tap(submitButton)
	assert.True(t, submitted)

	test.Tap(rstButton)
	test.AssertImageMatches(t, "form/form_initial.png", w.Canvas().Capture())
}

func TestForm_Save(t *testing.T) {
	var data struct {
		Name     string
		LastName string
		Age      int
	}

	name := NewTextFormField("Name", "")
	name.Validator = svalid.NotEmpty()
	name.Hint = "Your name"
	name.OnSaved = func(s string) { data.Name = s }

	lastName := NewTextFormField("LastName", "Pérez")
	lastName.Validator = svalid.NotEmpty()
	lastName.Hint = "Your lastname"
	lastName.OnSaved = func(s string) { data.LastName = s }

	age := NewRestrictTextFormField("Age", "", RestrictInputInteger)
	age.Validator = svalid.NotEmpty()
	age.Hint = "Your age"
	age.OnSaved = func(s string) { data.Age, _ = strconv.Atoi(s) }

	f := NewForm(1, name, lastName, age)

	name.SetText("Peter")
	lastName.SetText("Parker")
	age.SetText("45")

	assert.Zero(t, data)

	f.Save()
	assert.NotZero(t, data)
	assert.Equal(t, "Peter", data.Name)
	assert.Equal(t, "Parker", data.LastName)
	assert.Equal(t, 45, data.Age)
}
