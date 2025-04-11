package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		require.Equal(t, TestExpected, Top(Text, false))
		require.Equal(t, TestAsteriskExpected, Top(Text, true))
	})
}
