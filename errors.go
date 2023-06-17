package cron

import "errors"

var (
	EmptyCronSchedule   = errors.New("cron schedule is empty")
	InvalidCronSchedule = errors.New("invalid cron schedule")
)
