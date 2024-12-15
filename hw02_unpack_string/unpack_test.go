package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "aaa0b", expected: "aab"},
		{input: "abc0d0e", expected: "abe"},
		{input: "abc0d0eйц4ы5", expected: "abeйццццыыыыы"},
		{input: "!2 3", expected: "!!   "},
		{input: "日2本3", expected: "日日本本本"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	_, err := Unpack("3abc")
	require.Truef(t, errors.Is(err, errors.New(FIRST_ELEMET_IS_DIGIT_UNPACK_ERROR)), "actual error %q", err)

	_, err = Unpack("aaa10b")
	require.Truef(t, errors.Is(err, errors.New(TWO_OR_MORE_DIGITS_TOGETHER_UNPACK_ERROR)), "actual error %q", err)

	_, err = Unpack("")
	require.Truef(t, errors.Is(err, errors.New(TWO_OR_MORE_DIGITS_TOGETHER_UNPACK_ERROR)), "actual error %q", err)

}
