package sh

type AppServer struct {
	AppServerId    int `xorm:"pk autoincr"`
	AppId          int
	HostId         int
	AppServerLocal string
}
