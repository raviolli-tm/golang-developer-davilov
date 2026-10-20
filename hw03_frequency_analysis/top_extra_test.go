package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTop10ShortTexts(t *testing.T) {
	t.Run("two words", func(t *testing.T) {
		require.Equal(t, []string{"a", "b"}, Top10("a b"))
	})

	t.Run("frequency and lexicographic order", func(t *testing.T) {
		require.Equal(t, []string{"a", "b", "c"}, Top10("b a a c b a"))
	})

	t.Run("less than ten words", func(t *testing.T) {
		require.Equal(t, []string{"one", "three", "two"}, Top10("one two three"))
	})

	t.Run("punctuation is part of a word", func(t *testing.T) {
		expected := []string{"and", "one", "cat", "cats", "dog,", "dog,two", "man"}
		require.Equal(t, expected, Top10("cat and dog, one dog,two cats and one man"))
	})
}

func TestTop10AsteriskRules(t *testing.T) {
	t.Run("case and edge punctuation", func(t *testing.T) {
		require.Equal(t, []string{"нога", "-------"}, Top10Asterisk("Нога нога! 'нога' ------- -"))
	})

	t.Run("internal punctuation is significant", func(t *testing.T) {
		require.Equal(t, []string{"dog,cat", "dog...cat", "dogcat"}, Top10Asterisk("dog,cat dog...cat dogcat"))
	})

	t.Run("hyphen inside a word", func(t *testing.T) {
		require.Equal(t, []string{"какой-то", "какойто"}, Top10Asterisk("какой-то какойто"))
	})
}
