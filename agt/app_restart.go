package agt

import (
	"ctl"
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetApps() ([]AppInfo, error) {
	var (
		apps []AppInfo
	)
	b, err := ioutil.ReadFile(AppInfoFile)
	if err != nil {
		return apps, err
	}

	err = json.Unmarshal(b, &apps)
	return apps, err
}

func GetApp(name string) (AppInfo, error) {
	var (
		apps []AppInfo
		app  AppInfo
		err  error
	)
	apps, err = GetApps()
	if err != nil {
		return app, err
	}
	for _, v := range apps {
		if v.AppName == name {
			app = v
		}
	}

	if app.AppName == "" || app.AppId == 0 {
		return app, ctl.Errorf("获取app信息异常")
	}
	return app, err
}

func AppRestart() {
	var (
		app AppInfo
		err error
	)
	if len(os.Args) != 3 {
		ctl.Panic(ctl.Errorf(""), "输入参数不正确")
	}
	app, err = GetApp(os.Args[2])
	ctl.FatalErr(err)
	err = app.Restart()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func (app *AppInfo) Restart() error {
	var (
		err error
	)
	err = app.CheckOpt()
	if err != nil {
		return err
	}
	err = app.Kill()
	if err != nil {
		return err
	}
	err = app.Start()
	return err
}
