package hw10programoptimization

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStatEdgeCases(t *testing.T) {
	t.Run("at sign in another field", func(t *testing.T) {
		data := `{"Id":1,"Name":"x@evil.com","Username":"u","Email":"user@good.com"}`

		stat, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"good.com": 1}, stat)
	})

	t.Run("domain with hyphen", func(t *testing.T) {
		data := `{"Id":1,"Email":"user@my-site.com"}`

		stat, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"my-site.com": 1}, stat)
	})

	t.Run("suffix must start with a dot", func(t *testing.T) {
		data := `{"Id":1,"Email":"user@notcom"}`

		stat, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Empty(t, stat)
	})

	t.Run("email is case insensitive", func(t *testing.T) {
		data := `{"Id":1,"Email":"USER@Good.COM"}`

		stat, err := GetDomainStat(bytes.NewBufferString(data), "COM")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"good.com": 1}, stat)
	})

	t.Run("invalid json returns an error", func(t *testing.T) {
		_, err := GetDomainStat(strings.NewReader(`{"Email":`), "com")
		require.Error(t, err)
	})

	t.Run("empty lines are skipped", func(t *testing.T) {
		data := "\n" + `{"Id":1,"Email":"user@good.com"}` + "\n\n"

		stat, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"good.com": 1}, stat)
	})

	t.Run("empty domain", func(t *testing.T) {
		stat, err := GetDomainStat(bytes.NewBufferString(`{"Email":"user@good.com"}`), "")
		require.NoError(t, err)
		require.Empty(t, stat)
	})

	t.Run("reader error", func(t *testing.T) {
		_, err := GetDomainStat(errorReader{}, "com")
		require.Error(t, err)
		require.True(t, errors.Is(err, errTestReader))
	})
}

var errTestReader = errors.New("test reader error")

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errTestReader
}
