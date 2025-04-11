package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTop10(t *testing.T) {
	var flag bool

	flag = t.Run("no words in empty string", func(t *testing.T) {
		assert.Len(t, Top("", false), 0)
	})

	flag = flag && t.Run("no words in empty string asterisk", func(t *testing.T) {
		assert.Len(t, Top("", true), 0)
	})

	flag = flag && t.Run("positive test", func(t *testing.T) {
		assert.Equal(t, TestExpected, Top(Text, false))
	})
	flag = flag && t.Run("positive test asterisk", func(t *testing.T) {
		assert.Equal(t, TestAsteriskExpected, Top(Text, true))
	})

	if !flag {
		t.FailNow()
	}
}
