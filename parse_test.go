package cron

import (
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
				minute:  newSet[uint8](60, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59),
				hour:    newSet[uint8](60, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23),
				day:     newSet[uint8](60, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31),
				month:   newSet[uint8](60, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12),
				weekday: newSet[uint8](60, 0, 1, 2, 3, 4, 5, 6),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "single digit cron",
			schedule: "1 1 1 1 1",
			want: &Cron{
				minute:  newSet[uint8](60, 1),
				hour:    newSet[uint8](60, 1),
				day:     newSet[uint8](60, 1),
				month:   newSet[uint8](60, 1),
				weekday: newSet[uint8](60, 1),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "double digit cron",
			schedule: "12 12 12 12 *",
			want: &Cron{
				minute:  newSet[uint8](60, 12),
				hour:    newSet[uint8](60, 12),
				day:     newSet[uint8](60, 12),
				month:   newSet[uint8](60, 12),
				weekday: newSet[uint8](60, 0, 1, 2, 3, 4, 5, 6),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "simple list cron",
			schedule: "1,12 1,12 1,12 1,12 1,2",
			want: &Cron{
				minute:  newSet[uint8](60, 1, 12),
				hour:    newSet[uint8](60, 1, 12),
				day:     newSet[uint8](60, 1, 12),
				month:   newSet[uint8](60, 1, 12),
				weekday: newSet[uint8](60, 1, 2),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "simple step cron",
			schedule: "*/5 */5 */5 */5 */5",
			want: &Cron{
				minute:  newSet[uint8](60, 0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55),
				hour:    newSet[uint8](60, 0, 5, 10, 15, 20),
				day:     newSet[uint8](60, 1, 6, 11, 16, 21, 26, 31),
				month:   newSet[uint8](60, 1, 6, 11),
				weekday: newSet[uint8](60, 0, 5),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "simple range cron",
			schedule: "1-4 1-4 1-4 1-4 1-4",
			want: &Cron{
				minute:  newSet[uint8](60, 1, 2, 3, 4),
				hour:    newSet[uint8](60, 1, 2, 3, 4),
				day:     newSet[uint8](60, 1, 2, 3, 4),
				month:   newSet[uint8](60, 1, 2, 3, 4),
				weekday: newSet[uint8](60, 1, 2, 3, 4),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "range with step cron",
			schedule: "1-4/2 1-4/2 1-4/2 1-4/2 1-4/2",
			want: &Cron{
				minute:  newSet[uint8](60, 2, 4),
				hour:    newSet[uint8](60, 2, 4),
				day:     newSet[uint8](60, 1, 3),
				month:   newSet[uint8](60, 1, 3),
				weekday: newSet[uint8](60, 2, 4),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "lists with range with step cron",
			schedule: "1-2,*/5 1-2,*/5 1-2,*/5 1-2,*/5 1-2,*/5",
			want: &Cron{
				minute:  newSet[uint8](60, 0, 1, 2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55),
				hour:    newSet[uint8](60, 0, 1, 2, 5, 10, 15, 20),
				day:     newSet[uint8](60, 1, 2, 6, 11, 16, 21, 26, 31),
				month:   newSet[uint8](60, 1, 2, 6, 11),
				weekday: newSet[uint8](60, 0, 1, 2, 5),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "lists with range with step cron - inverse",
			schedule: "*/5,1-2 */5,1-2 */5,1-2 */5,1-2 */5,1-2",
			want: &Cron{
				minute:  newSet[uint8](60, 0, 1, 2, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55),
				hour:    newSet[uint8](60, 0, 1, 2, 5, 10, 15, 20),
				day:     newSet[uint8](60, 1, 2, 6, 11, 16, 21, 26, 31),
				month:   newSet[uint8](60, 1, 2, 6, 11),
				weekday: newSet[uint8](60, 0, 1, 2, 5),
				utc:     true,
			},
			wantErr: false,
		},
		{
			name:     "error - empty",
			schedule: "",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - not enough parts",
			schedule: "* * * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - missing last val",
			schedule: "* * * * ",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - missing first val",
			schedule: " * * * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - something else entirely",
			schedule: "quick brown fox",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - value below min",
			schedule: "* * 0 * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - value over max",
			schedule: "* * 100 * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric cron part",
			schedule: "cat cat cat cat cat",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric lower range part",
			schedule: "* * cat-4 * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric upper range part",
			schedule: "* * 4-cat * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric lower step part",
			schedule: "* * cat/4 * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric upper step part",
			schedule: "* 4/cat * * *",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric step part",
			schedule: "* * * * cat,4",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - non numeric step part",
			schedule: "* * * * cat,4",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "error - range min > max",
			schedule: "12-6 * * * *",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.schedule)
			if err != nil {
				assert.True(t, tt.wantErr)
			} else {
				assert.False(t, tt.wantErr)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}
