package provider

import (
	"math"
	"math/rand"
	"time"
)

const (
	YEAR_BASE  = 1900
	YEAR_MIN   = 1800
	YEAR_MAX   = 2299
	CYCLE_SIZE = 500
)

var GENDER2INT = map[string]int{
	"female": 0,
	"male":   1,
}

var INT2GENDER = map[int]string{
	0: "female",
	1: "male",
}

var MONTHSOFYEAR = [13]int{31, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var WEIGHTS = [10]int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}

func isLeapYear(year int) bool {
	leapFlag := false
	if year%4 == 0 {
		if year%100 == 0 {
			if year%400 == 0 {
				leapFlag = true
			} else {
				leapFlag = false
			}
		} else {
			leapFlag = true
		}
	} else {
		leapFlag = false
	}
	return leapFlag
}

func monthOffset(year int) int {
	calculatedOffset := (year/100 - 4) % 5 * 20
	return calculatedOffset
}

func lastDayOfMonth(year int, month int) int {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return lastDay.Day()
}

func randomMonth(year, day int) int {
	isLeap := isLeapYear(year)

	for {
		month := rand.Intn(12) + 1
		maxDays := MONTHSOFYEAR[month]
		if month == 2 && isLeap {
			maxDays += 1
		}
		if maxDays <= day {
			return month
		}
	}
}

func randomGenderValue(gender int) int {
	randomValue := rand.Intn(10)
	if randomValue%2 == gender {
		return randomValue
	} else {
		return int(math.Min(math.Max(float64(randomValue-1), 0), 9))
	}
}

func checksum(pesel string) int {
	_checksum := 0

	for idx, char := range pesel {
		if idx > 9 {
			break
		}
		_checksum += WEIGHTS[idx] * (int(char) - '0')
	}

	return (10 - (_checksum % 10)) % 10
}
