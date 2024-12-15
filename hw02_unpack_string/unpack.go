package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	EMPTY_INPUT_UNPACK_ERROR                 = "input string is empty"
	TWO_OR_MORE_DIGITS_TOGETHER_UNPACK_ERROR = "two or more digits together"
	FIRST_ELEMET_IS_DIGIT_UNPACK_ERROR       = "first element is digit"
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
		return nil, errors.New(EMPTY_INPUT_UNPACK_ERROR)
	}
	twoOrMoreDigitsTogetherRegexp := regexp.MustCompile(`\d{2,}`)

	if twoOrMoreDigitsTogetherRegexp.MatchString(input) {
		return nil, errors.New(TWO_OR_MORE_DIGITS_TOGETHER_UNPACK_ERROR)
	}

	runes := []rune(input)
	if unicode.IsDigit(runes[0]) {
		return nil, errors.New(FIRST_ELEMET_IS_DIGIT_UNPACK_ERROR)
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
