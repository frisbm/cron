package cron

import "time"

type Cron struct {
	minute    []uint8
	hour      []uint8
	day       []uint8
	month     []uint8
	dayOfWeek []uint8
	location  *time.Location
}

func (c *Cron) UseLocal() {
	c.location = time.Local
}
