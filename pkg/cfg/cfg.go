package cfg

import (
	"ops/pkg/ctl"
	"os"

	"ops/pkg/db"
	"ops/pkg/sh"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gogf/gf/os/glog"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	GameState        = make(map[int]*State)
	NormGameState    = make(map[int]*State)
	PdtState         = make(map[string]*State)
	HeaderState      []HeaderStateName
	NormHeaderState  []HeaderStateName
	PreHeaderState   []HeaderStateName
	Server           ServerCfg
	Platform         ServerCfg
	Task             RunTask
	Deploy           Dep
	GameHosts        = make(map[int]*sh.SSHHost)
	PltHosts         = make(map[int]*sh.SSHHost)
	OthHosts         = make(map[int]*sh.SSHHost)
	AllHosts         = make(map[int]*sh.SSHHost)
	GameStateHosts   = make(map[int]map[int]*sh.SSHHost)
	ScriptDir        = "/home/app/newgame/"
	GameFile         = "/opt/sh/list/games"
	UpdateLock       = &sync.Mutex{}
	GameAddLock      = &sync.Mutex{}
	UpdateQueueLock  = make(map[string]*sync.Mutex)
	AppUpdateLock    = &sync.Mutex{}
	App              AppConf
	RegGb, _         = regexp.Compile("^GB")
	RegLine, _       = regexp.Compile(".*[1|2|3|4|5|6|7|8|9|0]+服.*")
	RegNum, _        = regexp.Compile("[0-9]+")
	RegInt, _        = regexp.Compile("^[0-9]+$")
	RegIso, _        = regexp.Compile("^ISO-8859")
	RegBak, _        = regexp.Compile("^backup_file.*")
	RegSuc, _        = regexp.Compile("^升级.*成功.*")
	RegMinute, _     = regexp.Compile("^[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9] [0-9][0-9]:[0-9][0-9]$")
	RegAppSuc, _     = regexp.Compile("升级[^\n]*成功")
	RegBakSuc, _     = regexp.Compile(".*回退.*成功.*")
	RegMd5Val, _     = regexp.Compile("^[0-z]{32}  /.+")
	RegMd5, _        = regexp.Compile("^[0-z]{32}$")
	RegSvnVersion, _ = regexp.Compile("\n最后修改的版本: [0-9]+")
	RegConfString, _ = regexp.Compile("\".+\"")
	RegConfInt, _    = regexp.Compile("[0-9]+")
	RegConfFloat, _  = regexp.Compile("\\.+")
	//GameStates     = []int{2, 3, 4, 6, 7, 8}
	GameStates     []int
	RegSvr, _      = regexp.Compile("[1|2|3|4|5|6|7|8|9|0]+服")
	RegGame, _     = regexp.Compile(`^[\t| ]+"game":`)
	RegPort, _     = regexp.Compile(`"port":.[0-9]+`)
	SrvMap         = make(map[string]uint16)
	SrvLock        = &sync.Mutex{}
	GameHeaderHtml HeaderHtml

)

type HeaderHtml struct {
	DftState  string
	State     []HeaderStateName
	NormState []HeaderStateName
}
type HeaderStateName struct {
	StateName string
	FullName  string
}
type State struct {
	State     int
	Db        *xorm.Engine
	PortLike  string
	DftPdt    int
	LogHostId int
	StateName string
	IsPre     bool
	//	IsGame     bool
	IsServerId bool
	FullName   string
}

type Gdb struct {
	Host      string
	Port      int
	User      string
	Passwd    string
	Db        string
	State     int
	StateName string
	FullName  string
	PortLike  string
	IsPre     bool
	IsGame    bool
	DftPdt    int
	LogHostId int
}

type Envab struct {
	Ip      string
	System  string
	IsLinux bool
	OutIp   string
	Mysql   string
}

type RunTask struct {
	WaitUpdateAuto bool
	NoUpdateAuto   bool
	Online         bool
	DdosCheckPort  bool
	TabRecv        bool
	CheckDdosUsed  bool
	HostLoadData   bool
	ConfUpdate     bool
}

var (
	PdtOnline   = make(map[int]*PdtOnlines)
	StateOnline = make(map[int]*PdtOnlines)
)

type Online struct {
	Game   string
	Online int
}
type OnlineGame struct {
	State      int
	Online     []Online
	SrvOnline  []Online
	UpdateTime time.Time
	Total      int
}
type PdtOnlines struct {
	OnlineGame OnlineGame
}
type Cfg struct {
	LocalMysql string
	Db         Gdb
	Gdb        []Gdb
	Server     ServerCfg
	Platform   ServerCfg
	Dep        Dep
	App        AppConf
	Task       RunTask
}
type AppConf struct {
	Info         string
	GetInfo      string
	DftScript    string
	NginxUpDir   string
	NginxInfo    string
	NginxGetInfo string
}
type Dep struct {
	SourceDir string
	DepDir    string
	DelDir    string
	DftState  int
}
type ServerCfg struct {
	Listen  string
	Port    int
	SSL     bool
	IsLinux bool
}

func test1() {
	glog.Debug("abc")
	//	gcfg.GetContent("log")

}
func InitState() {
	for _, v := range GameState {
		PdtState[v.StateName] = v
	}
}

type Test1 struct {
	Id   int
	Name string
}

func InitDb(ldb Gdb) (*xorm.Engine, error) {
	var (
		ndb *xorm.Engine
		err error
		t1  []Test1
	)
	//	d1db, err = xorm.NewEngine("mysql", "test:test@("+env.Mysql+":3312)/test?charset=utf8")
	ctl.Debug(ldb.User + ":" + ldb.Passwd + "@(" + ldb.Host + ":" + strconv.Itoa(ldb.Port) + ")/" + ldb.Db + "?charset=utf8")
	//ndb, err = xorm.NewEngine("mysql", ldb.User+":"+ldb.Passwd+"@("+ldb.Host+":"+strconv.Itoa(ldb.Port)+")/"+ldb.Db+"?charset=utf8")
	ndb, err = xorm.NewEngine("sqlite3", "./conf/ops_test.db")
	ctl.Debug(ndb)
	ndb.Find(&t1)
	ctl.Debug(t1)
	//	os.Exit(1)
	ctl.Debug(ldb)
	if ldb.State == 0 {
		return ndb, err
	}
	ctl.Debug(ndb)
	ctl.Debug(err)
	GameStates = append(GameStates, ldb.State)
	GameState[ldb.State] = &State{}
	GameState[ldb.State].Db = ndb
	GameState[ldb.State].LogHostId = ldb.LogHostId
	GameState[ldb.State].PortLike = ldb.PortLike
	GameState[ldb.State].DftPdt = ldb.DftPdt
	GameState[ldb.State].State = ldb.State
	GameState[ldb.State].StateName = ldb.StateName
	GameState[ldb.State].FullName = ldb.FullName

	if ldb.IsGame {
		NormGameState[ldb.State] = GameState[ldb.State]
		NormHeaderState = append(NormHeaderState, HeaderStateName{StateName: ldb.StateName, FullName: ldb.FullName})
	}
	if ldb.IsPre {
		PreHeaderState = append(PreHeaderState, HeaderStateName{StateName: ldb.StateName, FullName: ldb.FullName})
	} else {
		HeaderState = append(HeaderState, HeaderStateName{StateName: ldb.StateName, FullName: ldb.FullName})
	}
	return ndb, err
}
func InitConf() {
	var (
		cfg Cfg
		err error
		dt  toml.MetaData
	)
	ctl.Debug(os.Stdout)
	dt, err = toml.DecodeFile("conf/ops.toml", &cfg)

	ctl.FatalErr(err)
	ctl.Debug(dt)
	ctl.Debug(cfg)
	db.Db, err = InitDb(cfg.Db)
	ctl.FatalErr(err)
	Deploy = cfg.Dep
	Task = cfg.Task
	for _, v := range cfg.Gdb {
		_, err = InitDb(v)
		ctl.FatalErr(err)
	}

	InitState()
	GameHeaderHtml = HeaderHtml{State: HeaderState, NormState: NormHeaderState}
	Server = cfg.Server
	Platform = cfg.Platform
	App = cfg.App
	go openHosts()

}

func GetSysValue(key string) string {
	var (
		cnf db.SystemEncConfig
	)
	db.Db.Where("conf_name=?", key).Get(&cnf)
	return cnf.GetValue()
}
