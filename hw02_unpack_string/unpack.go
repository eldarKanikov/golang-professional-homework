package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrEmptyInputUnpackError              = errors.New("input string is empty")
	ErrTwoOrMoreDigitsTogetherUnpackError = errors.New("two or more digits together")
	ErrFirstElementIsDigitUnpackError     = errors.New("first element is digit")
)

func Unpack(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyInputUnpackError
	}
	runes := []rune(input)

	if unicode.IsDigit(runes[0]) {
		return "", ErrFirstElementIsDigitUnpackError
	}

	var result strings.Builder
	for i := 0; i < len(runes); i++ {
		curr := runes[i]
		var next rune
		if i+1 < len(runes) {
			next = runes[i+1]
		}
		var partString string
		if !unicode.IsNumber(curr) {
			if unicode.IsNumber(next) {
				coeff, _ := strconv.Atoi(string(next))
				partString = strings.Repeat(string(curr), coeff)
			} else {
				partString = string(curr)
			}
		} else {
			if unicode.IsNumber(next) {
				return "", ErrTwoOrMoreDigitsTogetherUnpackError
			}
		}
		result.WriteString(partString)
	}
	return result.String(), nil
}
