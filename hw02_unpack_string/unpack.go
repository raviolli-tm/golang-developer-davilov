package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(x string) (string, error) {

	var runes = []rune(x)
	var flag bool = false
	var b strings.Builder

	for i := 0; i < len(runes); i++ {
		if unicode.IsNumber(runes[i]) {
			if !flag {
				return "", ErrInvalidString
			} else {
				b.WriteString(strings.Repeat(string(runes[i-1]), int(runes[i]-'0')))
				flag = false
			}
		} else if runes[i] == '\\' {

			if flag {
				b.WriteRune(runes[i-1])
			}
			i++
			if i < len(runes) && (unicode.IsNumber(runes[i]) || runes[i] == '\\') {
				flag = true
			} else {
				return "", ErrInvalidString
			}

		} else {
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
