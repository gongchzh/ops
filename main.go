package main

import (
	"ops/pkg/ctl"
	"ops/pkg/cfg"

	"runtime/debug"
	"time"
)

func main() {
	var (
		logFile string
	)
	cfg.InitConf()
	logFile = "log/"
	ctl.InitLog("", logFile, "ops")
	ctl.Fatal()
	defer func() {
		if except := recover(); except != nil {
			ctl.Log.Error("[崩溃]", except, string(debug.Stack()))
		}
	}()
	go ctl.Rotate(time.Minute * 10)
	time.Sleep(time.Second * 6)
	time.Sleep(time.Second * 4)

	go router()

	time.Sleep(time.Hour * 240)
}
