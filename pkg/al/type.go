package al

import (
	"ops/pkg/sh"
)

type App struct {
	AppId             int `xorm:"pk autoincr"`
	AppName           string
	Port              int
	Status            int
	AppTypeId         int
	AppNginxId        int
	NgxNum            int
	FunId             int
	BasePath          string
	BackDir           string
	PortCheckInvl     int
	UpdatePath        string
	UpdateProgram     string
	AppProgram        string
	ProgramDir        string
	UpdateScriptPath  string
	RestartScriptPath string
}

type AppGroup struct {
	AppGroupId   int `xorm:"pk autoincr"`
	AppGroupName string
	OpsEarlTime  string
	OpsLateTime  string
	State        int
}

type AppGroupMember struct {
	AppGroupId int
	FunId      int
}
type AppGroupSubMember struct {
	AppGroupId int
	SubGroupId int
}

type AppType struct {
	AppTypeId   int `xorm:"pk autoincr"`
	AppTypeName string
}
type AppInfo struct {
	AppId            int
	AppServerId      int
	HostId           int
	HostName         string
	AppName          string
	Port             int
	RunState         int
	AppZone          int
	AppTypeId        int
	AppNginxId       int
	FunId            int
	BasePath         string
	BackDir          string
	UpdatePath       string
	UpdateProgram    string
	AppProgram       string
	ProgramDir       string
	UpdateScriptPath string
	LogPath          string
	LogFormat        string
	CheckUsedCmd     string
	CheckUsedNot     string
}

const (
	AppStatusNoNginx     = 0
	AppStatusRun         = 1
	AppStatusNginxExpt   = 2
	AppStatusPortExpt    = 3
	AppStatusProcessExpt = 4
	AppStatusUnknowExpt  = 5
)

type NginxServerInfo struct {
	AppName        string
	AppServerId    int
	AppServerLocal string
	RunState       int
}

type NginxInfo struct {
	ServerInfo       []NginxServerInfo
	AppNginx         AppNginx
	AppNginxServerId int
	HostId           int
}

/*
type AppServer struct {
	AppServerId    int `xorm:"pk"`
	AppId          int
	HostId         int
	AppServerLocal string
}*/
type AppServer sh.AppServer

type AppList struct {
	AppId         int
	AppServerId   int
	AppName       string
	AppProgram    string
	AppNginxId    int
	UpdateProgram string
	CurMd5        string
	Port          int
	FunId         int
	HostName      string
	Pid           int
	HostId        int
	RunState      int
	UpdateTime    string
	NewMd5        string
	NewTime       string
}

type AppList1 struct {
	AppId         string
	AppServerId   string
	AppName       string
	AppProgram    string
	UpdateProgram string
	CurMd5        string
	Port          string
	HostName      string
	HostId        string
	RunState      string
	UpdateTime    string
	NewMd5        string
	NewTime       string
}

type AppNginx struct {
	AppNginxId        int `xorm:"pk autoincr"`
	AppNginxName      string
	NginxFirstScript  string
	NginxSecondScript string
	//	NginxViewScript   string
	NginxFile       string
	FunId           int
	NginxFirstFile  string
	NginxSecondFile string
}

type AppNginxServer struct {
	AppNginxServerId int `xorm:"pk autoincr"`
	AppNginxId       int
	HostId           int
}

type AppFun struct {
	FunId        int `xorm:"pk autoincr"`
	FunName      string
	ParentFunId  int
	OpsEarlTime  string `xorm:"datetime"`
	AppZone      int
	SucForce     int
	UpdateState  int
	CheckUsedCmd string
	CheckUsedNot string
	OpsLateTime  string `xorm:"datetime"`
}

type AppUpdateQueue struct {
	UpdateQueueId int `xorm:"pk autoincr"`
	FunId         int
	Md5           string
	Info          string
	UpEarlTime    string
	UpLateTime    string
	Status        int
	StartTime     string
	UpdateTime    string
	EndTime       string
	Back          string
}
