package hw02unpackstring

import (
	"errors"
	"regexp"
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
	runes, err := validateAndReturnRunes(input)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for i := 0; i < len(runes); i++ {
		curr, next := readCurrNextElements(runes, i)
		var partString string
		if isNotNumber(curr) {
			if isCoefficient(next) {
				coeff, _ := strconv.Atoi(string(next))
				partString = unpackPartString(string(curr), coeff)
			} else {
				partString = string(curr)
			}
		}
		result.WriteString(partString)
	}
	return result.String(), nil
}

func validateAndReturnRunes(input string) ([]rune, error) {
	if input == "" {
		return nil, ErrEmptyInputUnpackError
	}
	twoOrMoreDigitsTogetherRegexp := regexp.MustCompile(`\d{2,}`)

	if twoOrMoreDigitsTogetherRegexp.MatchString(input) {
		return nil, ErrTwoOrMoreDigitsTogetherUnpackError
	}

	runes := []rune(input)
	if unicode.IsDigit(runes[0]) {
		return nil, ErrFirstElementIsDigitUnpackError
	}
	return runes, nil
}

func readCurrNextElements(runes []rune, index int) (curr rune, next rune) {
	curr = runes[index]
	if index+1 < len(runes) {
		next = runes[index+1]
	}
	return curr, next
}

func isNotNumber(r rune) bool {
	return !unicode.IsNumber(r)
}

func isCoefficient(r rune) bool {
	return isNumber(r)
}

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func unpackPartString(partString string, count int) string {
	return strings.Repeat(partString, count)
}
