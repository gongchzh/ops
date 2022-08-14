package db

import (
	"strings"
)

type Host struct {
	HostId   int    `xorm:"pk"`
	HostName string `xorm:"unique"`
	Local    string
	Remote   string
	//	PdtId    int
	StateSub int
	State    int
}

func (h *Host) GetLocalIp() string {
	return strings.Split(h.Local, ":")[0]
}
func (h *Host) GetRemoteIp() string {
	return strings.Split(h.Remote, ":")[0]
}

type Service struct {
	SrvId      int //`xorm:pk`
	SrvType    string
	NewMd5     string
	OldMd5     string
	Md5Path    string
	CreateTime string `xorm:"DATETIME"`
	BackDir    string
	Num        int
	UpdateCmd  string
	State      int
}
type ServiceMember struct {
	SrvMemberId int //`xorm:pk`
	SrvId       int
	SrvName     string
}
type UpdateService struct {
	SrvId      int    `xorm:"pk"`
	HostId     int    `xorm:"pk"`
	Md5        string `xorm:"pk"`
	UpdateTime string
	BackFile   string
	State      int
}

type GamePdt struct {
	PdtId         int    `xorm:"pk"`
	PdtName       string `xorm:"unique"`
	PdtGame       string
	PdtFullName   string `xorm:"unique"`
	State         int
	SourceDir     string
	LogHostId     int
	Robot         string
	UpdateDir     string
	UpdateNextDir string
}

type HostGroup struct {
	HostGroupId int `xorm:"pk"`
	GroupName   string
}

type HostGroupMember struct {
	HostId      int `xorm:"pk"`
	HostGroupId int `xorm:"pk"`
}
type HostGroupUserMember struct {
	HostGroupId int `xorm:"pk"`
	UserId      int `xorm:"pk"`
}

type HostUser struct {
	UserId   int `xorm:"pk"`
	Name     string
	FullName string `xorm:"unique"`
	Password string
	State    int
}

type HostUserMember struct {
	HostId int `xorm:"pk"`
	UserId int `xorm:"pk"`
}

type Auth struct {
	AuthId     int `xorm:"pk"`
	User       string
	Password   string
	Name       string
	State      uint8
	AuthStatus int
	Phone      string
}

type UserGroup struct {
	UserGroupId int //`xorm:"pk"`
	Name        string
}

type UserGroupMember struct {
	UserGroupId int //`xorm:"pk"`
	AuthId      int //`xorm:"pk"`
}

type SrvOpt struct {
	SrvId      int //`xorm:pk`
	SrvType    string
	SrvName    string
	Md5Path    string
	Md5Name    string
	BasePath   string
	ConfName   string
	StartName  string
	State      int
	Port       int
	UpdateName string
	SourceName string
	BackDir    string
}
type SystemEncConfig struct {
	ConfId   int //`xorm:pk`
	ConfName string
	Value    string
}

func (cnf *SystemEncConfig) GetValue() string {
	return Aes.Decode(cnf.Value)
}
