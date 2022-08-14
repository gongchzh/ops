package agt

import (
	"ctl"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type SrvInfo struct {
	SrvName  string
	NewMd5   string
	LastMd5  string
	CurMd5   string
	CurTime  string
	LastTime string
	NewTime  string
}

type SrvOpt struct {
	SrvId      int //要修改`xorm:pk`
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

type Service struct {
	SrvOpt  SrvOpt
	SrvInfo SrvInfo
}

const (
	srvList = "/opt/sh/list/service.list"
	md5List = "/opt/sh/list/md5.list"
)

var (
	md5Diff bool
	newMd5  = make(map[string]*Md5List)
)

func (s *Service) GetOpt() error {
	b, err := ioutil.ReadFile(srvList)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &s.SrvOpt)
	return err
}

func getSrv() []Service {
	var (
		srv    []Service
		srvOpt []SrvOpt
	)
	b, err := ioutil.ReadFile(srvList)
	if err != nil {
		return nil
	}
	json.Unmarshal(b, &srvOpt)
	for _, v := range srvOpt {
		srv = append(srv, Service{SrvOpt: v})
	}
	return srv
}

type Md5List struct {
	Md5  string
	Time string
}

func (s *Service) GetInfo(ch chan SrvInfo, m map[string]*Md5List) error {
	var (
		err      error
		info     Md5List
		lastTime time.Time
		lastName string
		dir      []os.FileInfo
	)
	defer func() {
		ch <- s.SrvInfo
	}()
	if s.SrvOpt.SrvName == "" {
		return ctl.Errorf("service name error")
	}
	s.SrvInfo.SrvName = s.SrvOpt.SrvName
	if info, err = getMd5(m, s.SrvOpt.SrvName+"Cur", ctl.PathUinx(s.SrvOpt.Md5Path)+"/"+s.SrvOpt.Md5Name); err != nil {
		return err
	}
	s.SrvInfo.CurMd5 = info.Md5
	s.SrvInfo.CurTime = info.Time

	if info, err = getMd5(m, s.SrvOpt.SrvName+"New", s.SrvOpt.SourceName); err != nil {
		return err
	}
	s.SrvInfo.NewMd5 = info.Md5
	s.SrvInfo.NewTime = info.Time
	dir, err = ioutil.ReadDir(s.SrvOpt.BackDir)
	if err != nil {
		return err
	}
	for _, v := range dir {
		if !v.IsDir() {
			continue
		}
		if v.ModTime().UnixNano() > lastTime.UnixNano() {
			lastTime = v.ModTime()
			lastName = ctl.PathUinx(s.SrvOpt.BackDir) + "/" + v.Name() + "/" + s.SrvOpt.Md5Name
		}
	}
	if len(dir) == 0 {
		return nil
	}
	if lastName == "" {
		return ctl.Errorf("获取备份md5失败")
	}
	if info, err = getMd5(m, s.SrvOpt.SrvName+"Last", lastName); err != nil {
		return err
	}
	s.SrvInfo.LastMd5 = info.Md5
	s.SrvInfo.LastTime = info.Time
	return err
}
func getMdtList() map[string]*Md5List {
	var (
		err  error
		list map[string]*Md5List
	)
	b, err := ioutil.ReadFile(md5List)
	if err != nil {
		return nil
	}
	json.Unmarshal(b, &list)
	return list
}
func getMd5(m map[string]*Md5List, name, path string) (Md5List, error) {
	var (
		info os.FileInfo
		list Md5List
	)
	info, err = os.Stat(path)
	if err != nil {
		return list, err
	}
	list.Time = info.ModTime().Format(ctl.TimeFormatFile)
	if _, ok := m[name]; ok {
		if m[name].Time != list.Time {
			md5Diff = true
		} else {
			newMd5[name].Md5 = m[name].Md5
			newMd5[name].Time = m[name].Time
			return *m[name], nil
		}
	}
	md5Diff = true
	list.Md5, err = ctl.Md5sum(path)
	newMd5[name].Md5 = list.Md5
	newMd5[name].Time = list.Time
	return list, err
}

func serviceGetInfo() {
	var (
		srv  []Service
		m1   map[string]*Md5List
		m2   map[string]*SrvInfo
		ch   = make(chan SrvInfo)
		info SrvInfo
	)

	//	b, err := ioutil.ReadFile(srvList)
	srv = getSrv()
	if srv == nil {
		ctl.Panic(ctl.Errorf("没有数据"), "获取service")
	}
	for _, v := range srv {
		newMd5[v.SrvOpt.SrvName+"Last"] = &Md5List{}
		newMd5[v.SrvOpt.SrvName+"Cur"] = &Md5List{}
		newMd5[v.SrvOpt.SrvName+"New"] = &Md5List{}
	}
	m1 = getMdtList()
	for k, _ := range srv {
		go srv[k].GetInfo(ch, m1)
	}
	m2 = make(map[string]*SrvInfo)
	for range srv {
		info = SrvInfo{}
		info = <-ch
		if info.SrvName == "" {
			continue
		}

		m2[info.SrvName] = &info
	}
	if md5Diff {
		b, err := json.Marshal(&newMd5)
		ctl.Panic(err, "转换md5信息为json")
		ioutil.WriteFile(md5List, b, 0644)
		ctl.Panic(err, "写入md5文件")
	}
	b, err := json.Marshal(&m2)
	ctl.Panic(err, "转换服务信息为json")
	fmt.Print(string(b))
}
