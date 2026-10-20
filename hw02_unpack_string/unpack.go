package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func Unpack(x string) (string, error) {
	runes := []rune(x)
	var pending bool
	var b strings.Builder

	for i := 0; i < len(runes); i++ {
		switch {
		case isDigit(runes[i]):
			if !pending {
				return "", ErrInvalidString
			}
			b.WriteString(strings.Repeat(string(runes[i-1]), int(runes[i]-'0')))
			pending = false

		case runes[i] == '\\':
			if pending {
				b.WriteRune(runes[i-1])
			}
			i++
			if i < len(runes) && (isDigit(runes[i]) || runes[i] == '\\') {
				pending = true
			} else {
				return "", ErrInvalidString
			}

		default:
			if pending {
				b.WriteRune(runes[i-1])
			}
			pending = true
		}
	}
	if pending {
		b.WriteRune(runes[len(runes)-1])
	}
	return b.String(), nil
}
