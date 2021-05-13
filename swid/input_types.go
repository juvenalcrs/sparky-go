package swid

import (
	"strings"
	"unicode"
)

// RestrictInput defines input type for Restricted TextFields.
type RestrictInput int

// RestrictInput options
const (
	RestrictInputText RestrictInput = iota
	RestrictInputInteger
	RestrictInputFloat
	RestrictInputEmail
)

func acceptChar(input RestrictInput, text string, newChar rune, colPos int) bool {
	if input == RestrictInputText {
		return true
	}
	if input == RestrictInputEmail {
		return newChar != 'ñ' && (unicode.IsLetter(newChar) ||
			unicode.IsDigit(newChar) ||
			newChar == '-' || newChar == '_' ||
			newChar == '.' || newChar == '@')
	}
	if input == RestrictInputInteger || input == RestrictInputFloat {
		isNegative := strings.ContainsRune(text, '-')
		if isNegative && newChar == '-' {
			return false
		}
		if isNegative && colPos == 0 {
			return false
		}
		if newChar == '-' && colPos == 0 {
			return true
		}
		if input == RestrictInputFloat && newChar == '.' && !strings.ContainsRune(text, '.') {
			return true
		}
		return unicode.IsDigit(newChar)
	}

	return true
}
