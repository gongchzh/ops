package ctl

import (
	//	"fmt"
	"regexp"
	"strconv"
)

var (
	regF2, _ = regexp.Compile("[0-9]+\\.[0-9][0-9]")
)

func bAtof2(s string) (float32, error) {
	s = regF2.FindString(s)
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}
func Atof2(s interface{}) (float32, error) {
	var (
		st  string
		f   float64
		err error
	)
	switch value := s.(type) {
	case string:
		f, err = strconv.ParseFloat(value, 32)
		st = strconv.FormatFloat(f, 'f', 2, 32)
		f, err = strconv.ParseFloat(st, 32)
		return float32(f), err
	case float32:
		st = strconv.FormatFloat(float64(value), 'f', 2, 32)
		f, err := strconv.ParseFloat(st, 32)
		return float32(f), err
	case int:
		st = strconv.FormatFloat(float64(value), 'f', 2, 32)
		f, err := strconv.ParseFloat(st, 32)
		return float32(f), err
	case uint:
		st = strconv.FormatFloat(float64(value), 'f', 2, 32)
		f, err := strconv.ParseFloat(st, 32)
		return float32(f), err
	case float64:
		st = strconv.FormatFloat(value, 'f', 2, 32)
		f, err := strconv.ParseFloat(st, 32)
		return float32(f), err
	}
	return 0, Errorf("type error")

}

func NumIn(nc, ns, ne int) bool {
	if nc >= ns && nc <= ne {
		return true
	}
	return false
}

func NumNotIn(nc, ns, ne int) bool {
	if nc >= ns && nc <= ne {
		return false
	}
	return true
}
