package ctl

import (
	"strconv"
)

func StrToU16(str string) uint16 {
	var (
		intn int
	)
	intn, _ = strconv.Atoi(str)
	return uint16(intn)
}
