package ctl

import (
	"os"
	"syscall"
)

func Fatal() {
	var (
		err error
		f   *os.File
	)
	f, err = os.OpenFile("./fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		Log.Debug(err.Error())
	}
	f.WriteString(Now() + "\n")
	syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))

}
