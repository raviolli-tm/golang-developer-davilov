package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		assert.Len(t, Top("", false), 0)
	})

	t.Run("no words in empty string asterisk", func(t *testing.T) {
		assert.Len(t, Top("", true), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		assert.Equal(t, TestExpected, Top(Text, false))
	})
	t.Run("positive test asterisk", func(t *testing.T) {
		assert.Equal(t, TestAsteriskExpected, Top(Text, true))
	})
}
