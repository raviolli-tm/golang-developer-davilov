package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(x string) (string, error) {
	runes := []rune(x)
	var flag bool
	var b strings.Builder

	for i := 0; i < len(runes); i++ {
		switch {
		case unicode.IsNumber(runes[i]):
			if !flag {
				return "", ErrInvalidString
			}
			b.WriteString(strings.Repeat(string(runes[i-1]), int(runes[i]-'0')))
			flag = false

		case runes[i] == '\\':
			if flag {
				b.WriteRune(runes[i-1])
			}
			i++
			if i < len(runes) && (unicode.IsNumber(runes[i]) || runes[i] == '\\') {
				flag = true
			} else {
				return "", ErrInvalidString
			}
		default:
			if flag {
				b.WriteRune(runes[i-1])
			}
			flag = true
		}
	}
	if flag {
		b.WriteRune(runes[len(runes)-1])
	}
	return b.String(), nil
}
