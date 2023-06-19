package cron

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type partType string

const (
	minute  partType = "minute"
	hour    partType = "hour"
	Day     partType = "day"
	month   partType = "month"
	weekday partType = "weekday"
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
	// If schedule is empty, return error
	if schedule == "" {
		return nil, EmptyCronSchedule
	}

	cronParts := strings.Split(schedule, " ")
	// If the length of all the parts after splitting is not 5, return error
	if len(cronParts) != 5 {
		return nil, InvalidCronSchedule
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 5)
	cron := &Cron{
		utc: true,
	}

	// Using sync.WaitGroup, we can parse the 5 parts independently and concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var err error
			switch i {
			case 0:
				cron.minute, err = parseCronPart(cronParts[i], 0, 59, minute)
			case 1:
				cron.hour, err = parseCronPart(cronParts[i], 0, 23, hour)
			case 2:
				cron.day, err = parseCronPart(cronParts[i], 1, 31, Day)
			case 3:
				cron.month, err = parseCronPart(cronParts[i], 1, 12, month)
			case 4:
				cron.weekday, err = parseCronPart(cronParts[i], 0, 6, weekday)
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

// parseCronPart does all the heavy lifting of turning a cron part
// into an set of values to use in the Cron struct
func parseCronPart(cronPart string, min, max uint8, part partType) (set[uint8], error) {
	var err error
	offset := 0
	// Day & month start with 1 instead of 0
	if part == Day || part == month {
		offset = -1
	}

	timeSet := newSet[uint8](int(max))

	// Simple Validation for empty cron part
	if cronPart == "" {
		return timeSet, InvalidCronSchedule
	}

	// Easiest case, if the cron part is only '*' that means get all values for that part
	if cronPart == "*" {
		timeSet.add(rangeSlice(min, max, 1, offset)...)
		return timeSet, nil
	}

	// 1. Find & Split cron part list components, these are independent of each other
	list := strings.Split(cronPart, ",")

	// 2. Cycle through the list components
	for _, item := range list {

		// 3. Find and split Step Components
		steps := strings.Split(item, "/")
		step := uint8(1)

		// 4. If part is a step component, save in step
		if len(steps) == 2 {
			step, err = aToi8(steps[1], min, max)
			if err != nil {
				return set[uint8]{}, err
			}
		}
		// 5. If first part of split is * (i.e. */5) then we can create range slice and continue
		if steps[0] == "*" {
			timeSet.add(rangeSlice(min, max, step, offset)...)
			continue
		}

		// 6. Find and split range component
		ranges := strings.Split(steps[0], "-")
		var toi8, localMin, localMax uint8

		// 7. If part is a range component, find local min/max of the component,
		// validate, and create range slice using the saved step from earlier
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

		// 8. If part is simply an integer, convert to uint8 and add to timeSet
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
	result := make([]uint8, length)
	idx := 0
	for i := start; i <= end; i++ {
		// If 'i' is divisible by the step, add to return slice
		if i%step == 0 {
			// subtract the offset
			result[idx] = uint8(int(i) - offset)
			idx++
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
