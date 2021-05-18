package svalid

import (
	"errors"
	"fmt"
	"net/mail"

	"fyne.io/fyne/v2"
)

// NotEmpty defines not empty validator.
func NotEmpty() fyne.StringValidator {
	return func(s string) error {
		if s == "" {
			return errors.New(errMsgs.NotEmpty)
		}
		return nil
	}
}

// Email defines email validator.
func Email() fyne.StringValidator {
	return func(s string) error {
		_, err := mail.ParseAddress(s)
		if err != nil && errMsgs.Email != "" {
			return errors.New(errMsgs.Email)
		}
		return err
	}
}

// MinLength defines min length validator.
func MinLength(min int) fyne.StringValidator {
	return func(s string) error {
		if len([]rune(s)) < min {
			return fmt.Errorf(errMsgs.MinLength, min)
		}
		return nil
	}
}
