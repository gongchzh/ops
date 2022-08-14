package cfg

import (
	"ops/pkg/ctl"
	"ops/pkg/db"
	"ops/pkg/sh"
	"sync"
	"time"
)

func openHosts() {
	var (
		host []db.Host
		pdt  []db.GamePdt
		lid  []int
	)
	db.Db.In("state", GameStates).Find(&host)
	db.Db.In("state", GameStates).Find(&pdt)
	for _, v := range PdtState {
		GameStateHosts[v.State] = make(map[int]*sh.SSHHost)
		//pdtOnline[v.PdtId] = &PdtOnline{Lock: &sync.Mutex{}}
		StateOnline[v.State] = &PdtOnlines{}
	}
	ctl.Debug(len(host))
	for _, v := range host {
		GameHosts[v.HostId] = &sh.SSHHost{HostName: v.HostName, Lock: &sync.Mutex{}}
		if v.State == 4 {
			ctl.Debug(v.HostName, GameHosts[v.HostId], v.HostId)
		}
		go openGame(v)
		time.Sleep(time.Second / 100)
	}
	if true {
		host = nil
		db.Db.In("state", []int{403, 503, 11}).Find(&host)
		for _, v := range host {
			ctl.Debug(v.HostName, v.HostId)
			PltHosts[v.HostId] = &sh.SSHHost{HostName: v.HostName, Lock: &sync.Mutex{}}

			go openTest1(v)
		}
	}
	host = nil

	for k, v := range GameState {
		lid = append(lid, v.LogHostId)
		ctl.Debug(k)
	}

	if _, ok := GameState[2]; ok {
		lid = append(lid, 100, 101)
		lid = append(lid, 211)
	}
	db.Db.In("host_id", lid).Find(&host)
	for _, v := range host {
		OthHosts[v.HostId] = &sh.SSHHost{HostName: v.HostName, Lock: &sync.Mutex{}}
		go openOth(v)
	}
	time.Sleep(time.Second * 5)
	time.Sleep(time.Hour * 240)
	ctl.Debug("abc")

}

func openGame(h db.Host) {
	var (
		err error
	)

	//	ctl.Debug("start dial")
	err = GameHosts[h.HostId].SshDial()
	//	ctl.Debug(err)
	/*	if h.State == 4 {
		ctl.Debug(err)
		ctl.Debug(GameHosts[h.HostId])
	}*/
	if err != nil {
		ctl.Log.Error(h.HostName, h.HostId, err)
		ctl.Debug(h.HostName, h.HostId, err)
	} else {
		//ctl.Log.Info(h.HostName + " 创建ssh连接成功")
		ctl.Log.Debug(h.HostName + "创建ssh连接成功")

		//		ctl.Debug("创建ssh连接成功")
		go GameHosts[h.HostId].Keep(50)
	}
	if h.State < 20 {
		//	ctl.Debug(h.State)
		GameStateHosts[h.State][h.HostId] = GameHosts[h.HostId]
		if h.StateSub != 0 {
			GameStateHosts[h.StateSub][h.HostId] = GameHosts[h.HostId]
		}
		//	if h.State == 11 {
		//	}
	}
	AllHosts[h.HostId] = GameHosts[h.HostId]
}

func openTest1(h db.Host) {
	var (
		err error
	)
	time.Sleep(time.Second / 10)
	err = PltHosts[h.HostId].SshDial()
	if err != nil {
		ctl.Log.Error(h.HostName, err)
		ctl.Debug(h.HostName, err)
	} else {
		ctl.Log.Info(h.HostName + " 创建ssh连接成功")
		//		ctl.Debug("创建ssh连接成功")
		go PltHosts[h.HostId].Keep(50)
	}
	AllHosts[h.HostId] = PltHosts[h.HostId]
}

func openOth(h db.Host) {
	var (
		err error
	)
	time.Sleep(time.Second / 10)
	err = OthHosts[h.HostId].SshDial()
	if err != nil {
		ctl.Log.Error(h.HostName, err)
		ctl.Debug(h.HostName, err)
	} else {
		ctl.Log.Info(h.HostName + " 创建ssh连接成功")
		//		ctl.Debug("创建ssh连接成功")
		go OthHosts[h.HostId].Keep(50)
	}
	AllHosts[h.HostId] = OthHosts[h.HostId]
}
