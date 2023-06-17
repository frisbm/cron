package cron

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Parse(schedule string) (*Cron, error) {
	if schedule == "" {
		return nil, EmptyCronSchedule
	}

	cronParts := strings.Split(schedule, " ")
	if len(cronParts) != 5 {
		return nil, InvalidCronSchedule
	}

	var wg sync.WaitGroup
	var mu *sync.RWMutex
	var werr error
	cron := &Cron{
		location: time.UTC,
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			switch i {
			case 0:
				minute, err := parseCronPart(cronParts[i], 0, 59)
				if err != nil {
					mu.Lock()
					werr = errors.Join(werr, err)
					mu.Unlock()
					return
				}
				cron.minute = minute
			case 1:
				hour, err := parseCronPart(cronParts[i], 0, 23)
				if err != nil {
					mu.Lock()
					werr = errors.Join(werr, err)
					mu.Unlock()
					return
				}
				cron.hour = hour
			case 2:
				day, err := parseCronPart(cronParts[i], 1, 31)
				if err != nil {
					mu.Lock()
					werr = errors.Join(werr, err)
					mu.Unlock()
					return
				}
				cron.day = day
			case 3:
				month, err := parseCronPart(cronParts[i], 1, 12)
				if err != nil {
					mu.Lock()
					werr = errors.Join(werr, err)
					mu.Unlock()
					return
				}
				cron.month = month
			case 4:
				dayOfWeek, err := parseCronPart(cronParts[i], 0, 6)
				if err != nil {
					mu.Lock()
					werr = errors.Join(werr, err)
					mu.Unlock()
					return
				}
				cron.dayOfWeek = dayOfWeek
			}
		}(i)
	}
	wg.Wait()

	if werr != nil {
		return nil, errors.Join(InvalidCronSchedule, werr)
	}
	return cron, nil
}

func splitList(cronPart string) []string {
	return strings.Split(cronPart, ",")
}

func splitStep(cronPart string) []string {
	return strings.Split(cronPart, "/")
}

func splitRange(cronPart string) []string {
	return strings.Split(cronPart, "-")
}

func parseCronPart(cronPart string, min, max uint8) ([]uint8, error) {
	timeSet := NewSet[uint8]()
	var err error
	list := splitList(cronPart)
	for _, listItem := range list {
		steps := splitStep(listItem)
		stepItem := steps[0]
		var step uint8 = 1
		if len(steps) == 2 {
			step, err = aToi8(steps[1])
			if err != nil {
				return nil, err
			}
		} else if stepItem == "*" {
			timeSet.Add(rangeSlice(min, max, step)...)
			continue
		}

		ranges := splitRange(stepItem)
		var localMin, localMax = min, max
		if len(ranges) == 2 {
			localMin, err = aToi8(ranges[0])
			if err != nil {
				return nil, err
			}
			localMax, err = aToi8(ranges[1])
			if err != nil {
				return nil, err
			}
		} else if stepItem == "*" {
			timeSet.Add(rangeSlice(min, max, step)...)
			continue
		} else {
			toi8, err := aToi8(ranges[0])
			if err != nil {
				return nil, err
			}
			timeSet.Add(toi8)
			continue
		}
		timeSet.Add(rangeSlice(localMin, localMax, step)...)
		continue
	}

	return timeSet.Values(), nil
}

func rangeSlice(start, end, step uint8) []uint8 {
	if end < start {
		return nil
	}
	length := end - start
	vals := make([]uint8, 0, length)
	for i := start; i <= end; i++ {
		if i%step == 0 {
			vals = append(vals, i)
		}
	}
	return vals
}

func aToi8(a string) (uint8, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, nil
	}
	return uint8(i), nil
}
