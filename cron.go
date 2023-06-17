package cron

import (
	"golang.org/x/exp/slices"
	"time"
)

type Cron struct {
	minute    []uint8
	hour      []uint8
	day       []uint8
	month     []uint8
	dayOfWeek []uint8
	utc       bool
}

func (c *Cron) UseLocal() {
	c.utc = false
}

func (c *Cron) NextFrom(from time.Time) time.Time {
	if c.utc {
		from = from.UTC()
	}
	nextTime := from.Add(time.Minute)

	for {
		if slices.Contains(c.minute, uint8(nextTime.Minute())) &&
			slices.Contains(c.hour, uint8(nextTime.Hour())) &&
			slices.Contains(c.day, uint8(nextTime.Day())) &&
			slices.Contains(c.month, uint8(nextTime.Month())) &&
			slices.Contains(c.dayOfWeek, uint8(nextTime.Weekday())) {
			break
		}

		nextTime = nextTime.Add(time.Minute)
	}

	return nextTime
}

func (c *Cron) Next() time.Time {
	return c.NextFrom(c.now())
}

func (c *Cron) Prev() time.Time {
	return c.PrevBefore(c.now())
}

func (c *Cron) PrevBefore(before time.Time) time.Time {
	if c.utc {
		before = before.UTC()
	}
	prevTime := before.Add(-1 * time.Minute)

	for {
		if slices.Contains(c.minute, uint8(prevTime.Minute())) &&
			slices.Contains(c.hour, uint8(prevTime.Hour())) &&
			slices.Contains(c.day, uint8(prevTime.Day())) &&
			slices.Contains(c.month, uint8(prevTime.Month())) &&
			slices.Contains(c.dayOfWeek, uint8(prevTime.Weekday())) {
			break
		}

		prevTime = prevTime.Add(-1 * time.Minute)
	}

	return prevTime
}

func (c *Cron) Now() bool {
	now := c.now()
	if c.utc {
		now = now.UTC()
	}

	if !slices.Contains(c.minute, uint8(now.Minute())) {
		return false
	}

	if !slices.Contains(c.hour, uint8(now.Hour())) {
		return false
	}

	if !slices.Contains(c.day, uint8(now.Day())) {
		return false
	}

	if !slices.Contains(c.month, uint8(now.Month())) {
		return false
	}

	if !slices.Contains(c.dayOfWeek, uint8(now.Weekday())) {
		return false
	}

	return true
}

var timeNow = time.Now

func (c *Cron) now() time.Time {
	now := timeNow().Truncate(1 * time.Minute)
	if c.utc {
		now = now.UTC()
	}
	return now
}
