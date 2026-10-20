package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPeriodRange(t *testing.T) {
	date := time.Date(2026, 10, 21, 15, 30, 0, 0, time.UTC) // Wednesday

	tests := []struct {
		period Period
		from   time.Time
		to     time.Time
	}{
		{
			period: PeriodDay,
			from:   time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			period: PeriodWeek,
			from:   time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC), // Monday
			to:     time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			period: PeriodMonth,
			from:   time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
			to:     time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(string(test.period), func(t *testing.T) {
			from, to, err := PeriodRange(date, test.period)
			require.NoError(t, err)
			require.Equal(t, test.from, from)
			require.Equal(t, test.to, to)
		})
	}

	_, _, err := PeriodRange(date, Period("year"))
	require.ErrorIs(t, err, ErrInvalidPeriod)
}
