package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var res strings.Builder
	runes := []rune(s)
	n := len(runes)

	if n == 0 {
		return "", nil
	}

	var prev rune
	for i := 0; i < n; i++ {
		r := runes[i]

		if unicode.IsDigit(r) {
			if i == 0 || unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}

			count := int(r - '0')
			res.WriteString(strings.Repeat(string(prev), count))
		} else if i != 0 && !unicode.IsDigit(prev) {
			res.WriteRune(prev)
		}
		prev = r
	}

	if !unicode.IsDigit(prev) {
		res.WriteRune(prev)
	}

	return res.String(), nil
}
