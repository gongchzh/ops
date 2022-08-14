package db

type GameAddQueue struct {
	QueueId      int `xorm:"pk"`
	GameName     string
	Program      string
	PdtId        int
	OnlineGameId int
	GameType     string
	Status       int
	IsRedis      int
	UpdateTime   string `xorm:"DATETIME"`
	ZipPath      string
	Int1         int
	Int2         int
}

type GameConfOpt struct {
	GameId int `xorm:"pk"`
	Int1   int
	Int2   int
}

type GameAddConfInt struct {
	ConfIntId int
	ConfName  string
	FullName  string
	State     int
	Num       int
	Type      int
}

type GameServer struct {
	ServerId uint16 `xorm:"pk"`
	GameId   uint16
	HostId   int
	Server   string
	Remark   string
	BasePath string
	Port     int
	Update   int
}

type GameUpdateQueue struct {
	UpdateQueueId  int `xorm:"pk"`
	GameId         int
	UpdateInfo     string
	ConfPath       string
	Md5            string
	NewTabNum      int
	OldTabNum      int
	TabRecoverInvl int
	IsSameRun      int
	OthUpdateInvl  int
	Status         int
	BackFile       string
	SvnVersion     int
	CreateTime     string `xorm:"DATETIME"`
	EndTime        string `xorm:"DATETIME"`
}

type GamePeople struct {
	Id       int `xom:"pk"`
	ServerId uint16
	Times    int
	People   uint16
}
type GamePeopleTime struct {
	Time  string `xorm:"DATETIME"`
	Times int    `xorm:"pk"`
}

type GameServerOnline struct {
	SrvOnlineId      int `xorm:"pk"`
	ServerId         int
	State            int
	OnlineGameId     int
	OnlineServiceId  int
	OnlineGameName   string
	OnlineGameRemark string
	Status           int
	CreateTime       string `xorm:"DATETIME"`
	UpdateTime       string `xorm:"DATETIME"`
}

type GameConfQueue struct {
	ConfQueueId int //`xorm:pk`
	GameId      int
	Md5         string
	Status      int
	CompServer  int
	BackFile    string
	ConfPath    string
	OthConfInvl int
	CreateTime  string `xorm:"DATETIME"`
}

type GameHostLoad struct {
	SingleHostId int //`xorm:pk`
	DoubleHostId int
	SingleNum    int
	DoubleNum    int
	SinglePeople int
	DoublePeople int
}
