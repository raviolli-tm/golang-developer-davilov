package hw02unpackstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpackExtra(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{name: "zero repeats", input: "a0", expected: "", err: nil},
		{name: "escaped zero", input: `a\0`, expected: "a0", err: nil},
		{name: "escaped backslash", input: `\\`, expected: `\`, err: nil},
		{name: "escaped digit at start", input: `\4`, expected: "4", err: nil},
		{name: "non-ascii digit is a letter", input: "a٤", expected: "a٤", err: nil},
		{name: "non-ascii digit cannot be escaped", input: `\٤`, expected: "", err: ErrInvalidString},
		{name: "digit after escaped digit", input: `a\45`, expected: "a44444", err: nil},
		{name: "lone backslash", input: `\`, expected: "", err: ErrInvalidString},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Unpack(test.input)
			if test.err == nil {
				require.NoError(t, err)
				require.Equal(t, test.expected, result)
				return
			}

			require.ErrorIs(t, err, test.err)
		})
	}
}
