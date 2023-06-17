package cron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		want     *Cron
		wantErr  bool
	}{
		{
			name:     "base cron",
			schedule: "* * * * *",
			want: &Cron{
				minute:    []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
				hour:      []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
				day:       []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
				month:     []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				dayOfWeek: []uint8{0, 1, 2, 3, 4, 5, 6},
			},
			wantErr: false,
		},
		{
			name:     "single digit cron",
			schedule: "1 1 1 1 1",
			want: &Cron{
				minute:    []uint8{1},
				hour:      []uint8{1},
				day:       []uint8{1},
				month:     []uint8{1},
				dayOfWeek: []uint8{1},
			},
			wantErr: false,
		},
		{
			name:     "double digit cron",
			schedule: "12 12 12 12 *",
			want: &Cron{
				minute:    []uint8{12},
				hour:      []uint8{12},
				day:       []uint8{12},
				month:     []uint8{12},
				dayOfWeek: []uint8{0, 1, 2, 3, 4, 5, 6},
			},
			wantErr: false,
		},
		{
			name:     "simple list cron",
			schedule: "1,12 1,12 1,12 1,12 1,2",
			want: &Cron{
				minute:    []uint8{1, 12},
				hour:      []uint8{1, 12},
				day:       []uint8{1, 12},
				month:     []uint8{1, 12},
				dayOfWeek: []uint8{1, 2},
			},
			wantErr: false,
		},
		{
			name:     "simple step cron",
			schedule: "*/5 */5 */5 */5 */5",
			want: &Cron{
				minute:    []uint8{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
				hour:      []uint8{0, 5, 10, 15, 20},
				day:       []uint8{5, 10, 15, 20, 25, 30},
				month:     []uint8{5, 10},
				dayOfWeek: []uint8{0, 5},
			},
			wantErr: false,
		},
		{
			name:     "simple range cron",
			schedule: "1-4 1-4 1-4 1-4 1-4",
			want: &Cron{
				minute:    []uint8{1, 2, 3, 4},
				hour:      []uint8{1, 2, 3, 4},
				day:       []uint8{1, 2, 3, 4},
				month:     []uint8{1, 2, 3, 4},
				dayOfWeek: []uint8{1, 2, 3, 4},
			},
			wantErr: false,
		},
		{
			name:     "range with step cron",
			schedule: "1-4/2 1-4/2 1-4/2 1-4/2 1-4/2",
			want: &Cron{
				minute:    []uint8{2, 4},
				hour:      []uint8{2, 4},
				day:       []uint8{2, 4},
				month:     []uint8{2, 4},
				dayOfWeek: []uint8{2, 4},
			},
			wantErr: false,
		},
		{
			name:     "range with step cron",
			schedule: "1-4/2 1-4/2 1-4/2 1-4/2 1-4/2",
			want: &Cron{
				minute:    []uint8{2, 4},
				hour:      []uint8{2, 4},
				day:       []uint8{2, 4},
				month:     []uint8{2, 4},
				dayOfWeek: []uint8{2, 4},
			},
			wantErr: false,
		},
		{
			name:     "lists with range with step cron",
			schedule: "1-2,*/5 1-2,*/5 1-2,*/5 1-2,*/5 1-2,*/5",
			want: &Cron{
				minute:    []uint8{0, 1, 2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
				hour:      []uint8{0, 1, 2, 5, 10, 15, 20},
				day:       []uint8{1, 2, 5, 10, 15, 20, 25, 30},
				month:     []uint8{1, 2, 5, 10},
				dayOfWeek: []uint8{0, 1, 2, 5},
			},
			wantErr: false,
		},
		{
			name:     "lists with range with step cron - inverse",
			schedule: "*/5,1-2 */5,1-2 */5,1-2 */5,1-2 */5,1-2",
			want: &Cron{
				minute:    []uint8{0, 1, 2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55},
				hour:      []uint8{0, 1, 2, 5, 10, 15, 20},
				day:       []uint8{1, 2, 5, 10, 15, 20, 25, 30},
				month:     []uint8{1, 2, 5, 10},
				dayOfWeek: []uint8{0, 1, 2, 5},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.schedule)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.False(t, tt.wantErr)
				fmt.Println(fmt.Sprintf("%#v", got))
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}
