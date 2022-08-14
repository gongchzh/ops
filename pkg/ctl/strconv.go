package ctl

import (
	"strconv"
)

func Itoa(i int) string {
	return strconv.Itoa(i)
}
func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}
func AtoiNe(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
