package cron

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

/*
Parse takes a standard cron schedule (* * * * *) and returns
a Cron object if the schedule is valid; otherwise, it returns an error.
It uses goroutines to parse each part of the schedule concurrently, resulting
in faster parsing.

The schedule follows the standard format:

* [0-59] (* , / -)

* [0-23] (* , / -)

* [1-31] (* , / -)

* [1-12] (* , / -)

* [0-6]  (* , / -)
*/
func Parse(schedule string) (*Cron, error) {
	if schedule == "" {
		return nil, EmptyCronSchedule
	}

	cronParts := strings.Split(schedule, " ")
	if len(cronParts) != 5 {
		return nil, InvalidCronSchedule
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 5)
	cron := &Cron{
		utc: true,
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var err error
			switch i {
			case 0:
				cron.minute, err = parseCronPart(cronParts[i], 0, 59)
			case 1:
				cron.hour, err = parseCronPart(cronParts[i], 0, 23)
			case 2:
				cron.day, err = parseCronPart(cronParts[i], 1, 31)
			case 3:
				cron.month, err = parseCronPart(cronParts[i], 1, 12)
			case 4:
				cron.dayOfWeek, err = parseCronPart(cronParts[i], 0, 6)
			}
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()

	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, errors.Join(InvalidCronSchedule, err)
		}
	}
	return cron, nil
}

func parseCronPart(cronPart string, min, max uint8) (set[uint8], error) {
	timeSet := newSet[uint8](int(max))

	if cronPart == "" {
		return timeSet, InvalidCronSchedule
	}

	var err error
	list := strings.Split(cronPart, ",")
	for _, listItem := range list {
		steps := strings.Split(listItem, "/")
		stepItem := steps[0]
		step := uint8(1)
		if len(steps) == 2 {
			step, err = aToi8(steps[1], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
		} else if stepItem == "*" {
			timeSet.Add(rangeSlice(min, max, step)...)
			continue
		}

		ranges := strings.Split(stepItem, "-")
		var toi8, localMin, localMax uint8
		if len(ranges) == 2 {
			localMin, err = aToi8(ranges[0], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
			localMax, err = aToi8(ranges[1], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
			if localMin > localMax {
				return set[uint8]{}, errors.New("range min cannot be greater than range max")
			}
			timeSet.Add(rangeSlice(localMin, localMax, step)...)
		} else if stepItem == "*" {
			timeSet.Add(rangeSlice(min, max, step)...)
			continue
		} else {
			toi8, err = aToi8(ranges[0], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
			timeSet.Add(toi8)
			continue
		}
	}

	return timeSet, nil
}

func rangeSlice(start, end, step uint8) []uint8 {
	length := ((end - start) / step) + 1
	vals := make([]uint8, 0, length)
	for i := start; i <= end; i++ {
		if i%step == 0 {
			vals = append(vals, i)
		}
	}
	return vals
}

func aToi8(a string, min, max uint8) (uint8, error) {
	parsed, err := strconv.ParseUint(a, 10, 8)
	if err != nil {
		return 0, err
	}
	val := uint8(parsed)
	if val < min || val > max {
		return 0, InvalidCronSchedule
	}
	return val, nil
}
