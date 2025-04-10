package hw02unpackstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde", err: nil},
		{input: "abccd", expected: "abccd", err: nil},
		{input: "", expected: "", err: nil},
		{input: "aaa0b", expected: "aab", err: nil},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`, err: nil},
		{input: `qwe\45`, expected: `qwe44444`, err: nil},
		{input: `qwe\\5`, expected: `qwe\\\\\`, err: nil},
		{input: `qwe\\\3`, expected: `qwe\3`, err: nil},
		// My unit-tests
		{input: `	qw	e4\5`, expected: `	qw	eeee5`, err: nil},
		{input: `a\4\12\3`, expected: `a4113`, err: nil},
		{input: `	3	2`, expected: `					`, err: nil},
		{input: ` 3 \3	\1t`, expected: `    3	1t`, err: nil},
		// invalid string test-cases
		{input: "3abc", expected: "", err: ErrInvalidString},
		{input: "45", expected: "", err: ErrInvalidString},
		{input: "aaa10b", expected: "", err: ErrInvalidString},
		{input: `\a4aaa`, expected: "", err: ErrInvalidString},
		{input: `abcd\`, expected: "", err: ErrInvalidString},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			if tc.err == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}

// func TestUnpackInvalidString(t *testing.T) {
//	invalidStrings := []string{"3abc", "45", "aaa10b"}
//	for _, tc := range invalidStrings {
//		tc := tc
//		t.Run(tc, func(t *testing.T) {
//			_, err := Unpack(tc)
//			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
//		})
//	}
// }
