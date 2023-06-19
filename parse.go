package cron

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type PartType string

const (
	Minute  PartType = "Minute"
	Hour    PartType = "Hour"
	Day     PartType = "Day"
	Month   PartType = "Month"
	Weekday PartType = "Weekday"
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
				cron.minute, err = parseCronPart(cronParts[i], 0, 59, Minute)
			case 1:
				cron.hour, err = parseCronPart(cronParts[i], 0, 23, Hour)
			case 2:
				cron.day, err = parseCronPart(cronParts[i], 1, 31, Day)
			case 3:
				cron.month, err = parseCronPart(cronParts[i], 1, 12, Month)
			case 4:
				cron.weekday, err = parseCronPart(cronParts[i], 0, 6, Weekday)
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

func parseCronPart(cronPart string, min, max uint8, part PartType) (set[uint8], error) {
	timeSet := newSet[uint8](int(max))
	offset := 0
	if part == Day || part == Month {
		offset = -1
	}

	if cronPart == "" {
		return timeSet, InvalidCronSchedule
	}

	if cronPart == "*" {
		timeSet.add(rangeSlice(min, max, 1, offset)...)
		return timeSet, nil
	}

	var err error
	list := strings.Split(cronPart, ",")
	for _, item := range list {
		steps := strings.Split(item, "/")
		step := uint8(1)
		if len(steps) == 2 {
			step, err = aToi8(steps[1], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
		}
		if steps[0] == "*" {
			timeSet.add(rangeSlice(min, max, step, offset)...)
			continue
		}

		ranges := strings.Split(steps[0], "-")
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
			timeSet.add(rangeSlice(localMin, localMax, step, offset)...)
			continue
		}

		toi8, err = aToi8(ranges[0], min, max)
		if err != nil {
			return set[uint8]{}, err
		}
		timeSet.add(toi8)
	}

	return timeSet, nil
}

// rangeSlice takes a start, end, step, and offset value to
// create a slice of uint8
func rangeSlice(start, end, step uint8, offset int) []uint8 {
	// Add the offset
	start = uint8(int(start) + offset)
	end = uint8(int(end) + offset)

	// Able to calculate worst-case for the capacity of range slice
	length := ((end - start) / step) + 1
	result := make([]uint8, 0, length)

	for i := start; i <= end; i++ {
		// If i is divisible by the step, add to return slice
		if i%step == 0 {
			// subtract the offset
			val := uint8(int(i) - offset)

			result = append(result, val)
		}
	}
	return result
}

// aToi8 attempts to convert a string into uint8 with vaidation
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
