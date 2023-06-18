package cron

import (
	"time"
)

/*
Cron represents the cron schedule
*/
type Cron struct {
	minute    *set[uint8]
	hour      *set[uint8]
	day       *set[uint8]
	month     *set[uint8]
	dayOfWeek *set[uint8]
	utc       bool
}

/*
UseLocal will set the cron schedule calculations to use
the local time of the system it is running on
*/
func (c *Cron) UseLocal() {
	c.utc = false
}

/*
NextFrom accepts a time in which it will calculate the next activation time after
*/
func (c *Cron) NextFrom(from time.Time) time.Time {
	if c.utc {
		from = from.UTC()
	}
	nextTime := from.Add(time.Minute)

	for {
		if c.isTime(nextTime) {
			break
		}

		nextTime = nextTime.Add(time.Minute)
	}

	return nextTime
}

/*
Next will return the next cron activation after now
*/
func (c *Cron) Next() time.Time {
	return c.NextFrom(c.now())
}

/*
Prev will return the previous cron activation before now
*/
func (c *Cron) Prev() time.Time {
	return c.PrevBefore(c.now())
}

/*
PrevBefore accepts a time in which it will calculate the previous activation time before now
*/
func (c *Cron) PrevBefore(before time.Time) time.Time {
	if c.utc {
		before = before.UTC()
	}
	prevTime := before.Add(-1 * time.Minute)

	for {
		if c.isTime(prevTime) {
			break
		}

		prevTime = prevTime.Add(-1 * time.Minute)
	}

	return prevTime
}

/*
Now will tell you it is currently time for a cron schedule to activate
*/
func (c *Cron) Now() bool {
	now := c.now()
	if c.utc {
		now = now.UTC()
	}
	return c.isTime(now)
}

func (c *Cron) isTime(time time.Time) bool {
	if c.minute.Contains(uint8(time.Minute())) &&
		c.hour.Contains(uint8(time.Hour())) &&
		c.day.Contains(uint8(time.Day())) &&
		c.month.Contains(uint8(time.Month())) &&
		c.dayOfWeek.Contains(uint8(time.Weekday())) {
		return true
	}
	return false
}

var timeNow = time.Now

func (c *Cron) now() time.Time {
	now := timeNow().Truncate(1 * time.Minute)
	if c.utc {
		now = now.UTC()
	}
	return now
}
