package ctl

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TimeFormat     = "2006-01-02 15:04:05"
	TimeFormatT    = "2006-01-02T15:04:05"
	TimeDayFormat  = "2006-01-02"
	TimeFormatFile = "2006-01-02_15-04-05"
)

func Now() string {
	return time.Now().Format(TimeFormat)
}

func NowNum() string {
	n1 := Now()
	n1 = strings.ReplaceAll(n1, "-", "")
	n1 = strings.ReplaceAll(n1, ":", "")
	n1 = strings.ReplaceAll(n1, " ", "")
	return n1
}

func NowUnix() int {
	return int(time.Now().Unix())
}
func Nowf() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

func NowDay() string {
	return time.Now().Format("2006-01-02")
}
func NowYear() string {
	return time.Now().Format("2006")
}
func NowHour() string {
	return time.Now().Format("2006-01-02_15")
}
func NowMinute() string {
	return time.Now().Format("2006-01-02_15_04")
}
func TimeBefor(time1 string, time2 string, intday string) bool {
	regTime, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}")
	regDay, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}")
	var beforTime string
	interval, _ := time.ParseDuration(intday)

	if !regTime.MatchString(time1) && regDay.MatchString(time1) {
		time1 = time1 + " " + time.Now().Format("15:04:05")
	}
	if time2 == "" {
		beforTime = time.Now().Format(TimeFormat)
	} else {
		if !regTime.MatchString(time2) && regDay.MatchString(time1) {
			time2 = time2 + " " + time.Now().Format("15:04:05")
		}
		beforTime = time2
	}

	if false {
		fmt.Println("t1", time1)
		fmt.Println("bf", beforTime)
	}
	t1, _ := time.Parse(TimeFormat, beforTime)

	t1 = t1.Add(interval)
	t2, _ := time.Parse(TimeFormat, time1)
	if t1.Before(t2) {
		return false
	} else {
		return true
	}
}

func UnixToTime(i int64) (time.Time, error) {
	var (
		t      time.Time
		str    string
		i1, i2 int64
	)
	str = strconv.FormatInt(i, 10)
	if len(str) < 11 {
		return t, errors.New("unix time error")
	}
	i1, _ = strconv.ParseInt(str[0:10], 10, 64)
	i2, _ = strconv.ParseInt(str[10:], 10, 64)
	return time.Unix(i1, i2), nil
}

func SecondTime(t1 time.Time) time.Time {
	tf := t1.Format(TimeFormat)
	t2, _ := time.ParseInLocation(TimeFormat, tf, time.Local)
	return t2
}
