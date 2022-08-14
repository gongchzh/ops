package db

import (
	"ops/pkg/ctl"

	"xorm.io/xorm"
)

var (
	Db       *xorm.Engine
	eccKey11 = `49cc69680f52f0910c9b6e9d3e4db23a`
	Aes      = ctl.NewAes(eccKey11)
)

type Result struct {
	Result   string
	HostId   int
	HostName string
	Err      error
}

type AppResult struct {
	Result      string
	HostId      int
	HostName    string
	AppId       int
	AppServerId int
	Err         error
}
