package cron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCron_Next(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		now      time.Time
		want     time.Time
	}{
		{
			name:     "next min",
			schedule: "* * * * *",
			want:     time.Date(2023, 6, 17, 18, 24, 0, 0, time.UTC),
		},
		{
			name:     "next 5th min",
			schedule: "*/5 * * * *",
			want:     time.Date(2023, 6, 17, 18, 25, 0, 0, time.UTC),
		},
		{
			name:     "next top of the hour",
			schedule: "0 * * * *",
			want:     time.Date(2023, 6, 17, 19, 0, 0, 0, time.UTC),
		},
		{
			name:     "next day at midnight",
			schedule: "0 0 * * *",
			want:     time.Date(2023, 6, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "next month on first day at midnight",
			schedule: "0 0 1 * *",
			want:     time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "First day of the next year",
			schedule: "0 0 1 1 *",
			want:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return time.Date(2023, 6, 17, 18, 23, 0, 0, time.UTC)
			}
			cron, err := Parse(tt.schedule)
			assert.NoError(t, err)
			next := cron.Next()
			fmt.Println(tt.want)
			fmt.Println(next)
			assert.True(t, tt.want.Equal(next))
		})
	}
}

func TestCron_Prev(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		now      time.Time
		want     time.Time
	}{
		{
			name:     "prev min",
			schedule: "* * * * *",
			want:     time.Date(2023, 6, 17, 18, 22, 0, 0, time.UTC),
		},
		{
			name:     "prev 5th min",
			schedule: "*/5 * * * *",
			want:     time.Date(2023, 6, 17, 18, 20, 0, 0, time.UTC),
		},
		{
			name:     "prev top of the hour",
			schedule: "0 * * * *",
			want:     time.Date(2023, 6, 17, 18, 0, 0, 0, time.UTC),
		},
		{
			name:     "prev day at midnight",
			schedule: "0 0 * * *",
			want:     time.Date(2023, 6, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "prev month on first day at midnight",
			schedule: "0 0 1 * *",
			want:     time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "First day of the prev year",
			schedule: "0 0 1 1 *",
			want:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return time.Date(2023, 6, 17, 18, 23, 0, 0, time.UTC)
			}
			cron, err := Parse(tt.schedule)
			assert.NoError(t, err)
			prev := cron.Prev()
			assert.True(t, tt.want.Equal(prev))
		})
	}
}

func TestCron_Now(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		now      time.Time
		want     bool
	}{
		{
			name:     "is now",
			schedule: "* * * * *",
			want:     true,
		},
		{
			name:     "is not now - min",
			schedule: "5 * * * *",
			want:     false,
		},
		{
			name:     "is not now - hour",
			schedule: "* 5 * * *",
			want:     false,
		},
		{
			name:     "is not now - day",
			schedule: "* * 5 * *",
			want:     false,
		},
		{
			name:     "is not now - month",
			schedule: "* * * 5 *",
			want:     false,
		},
		{
			name:     "is not now - day of the week",
			schedule: "* * * * 5",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return time.Date(2023, 6, 17, 18, 23, 0, 0, time.UTC)
			}
			cron, err := Parse(tt.schedule)
			assert.NoError(t, err)
			now := cron.Now()
			assert.Equal(t, tt.want, now)
		})
	}
}

func TestCron_UseLocal(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
	}{
		{
			name:     "use local",
			schedule: "* * * * *",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return time.Now().Local()
			}
			cron, err := Parse(tt.schedule)
			assert.NoError(t, err)
			cron.UseLocal()
			assert.Equal(t, timeNow().Location(), cron.now().Location())
			assert.True(t, timeNow().Truncate(1*time.Minute).Equal(cron.now()))
		})
	}
}
