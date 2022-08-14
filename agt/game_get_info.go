package agt

import (
	"bufio"
	"bytes"
	"ctl"

	"io"
	"os"
	"strconv"
	"strings"
)

var (
	err         error
	cmdOut      bytes.Buffer
	gameMd5File *os.File
)

type Online struct {
	Game   string
	Online int
}
type GameMd5 struct {
	Md5  string
	Time string
}

const (
	PROC_TCP = "/proc/net/tcp"
)

func gameGetInfo() {
	var (
		chs     chan string
		out     bytes.Buffer
		i       int
		md5List map[string]GameMd5
		ch1     chan map[string]GameMd5
		buf     *bufio.Reader
		f       *os.File
		line    string
		games   []string
		game    []Game
	)
	ch1 = make(chan map[string]GameMd5)
	go getMd5File(ch1)
	f, err = os.Open("/opt/sh/list/games")
	ctl.FatalErr(err)
	buf = bufio.NewReader(f)
	for {
		line, err = buf.ReadString(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			ctl.FatalErr(err)
		}
		games = strings.Split(line, "___")
		game = append(game, Game{GameNote: games[0], Path: games[1] + "/" + games[2], Port: games[3], PortType: games[4]})
	}
	f.Close()
	chs = make(chan string)
	go getOnline(chs, game)
	md5List = <-ch1
	gameMd5File, err = os.OpenFile("/home/app/newgame/md5", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		ctl.FatalErr(err)
	}
	for i = 1; i <= len(game); i++ {
		go getInfo(chs, game[i-1], md5List)
	}
	for i = 0; i <= len(game); i++ {
		out.WriteString(<-chs)
	}
	out.WriteTo(os.Stdout)
}

func getMd5File(ch chan map[string]GameMd5) {
	var (
		gameMd5 []string
		md5List map[string]GameMd5
		buf     *bufio.Reader
		line    string
	)
	gameMd5File, err = os.Open("/home/app/newgame/md5")
	if err != nil {
		ch <- nil
		return
	}
	md5List = make(map[string]GameMd5)
	buf = bufio.NewReader(gameMd5File)
	for {
		line, err = buf.ReadString(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			ctl.Log.Error(err)
			ch <- nil
			return
		}
		gameMd5 = strings.Split(line, "___")
		md5List[gameMd5[0]] = GameMd5{Md5: gameMd5[1], Time: gameMd5[2][:len(gameMd5[2])-1]}
	}
	gameMd5File.Close()
	ch <- md5List
}

type Game struct {
	GameNote string
	Port     string
	PortType string
	Path     string
}

func getOnline(ch chan string, games []Game) {
	var (
		ports      map[int64]*Online
		out        bytes.Buffer
		line_array []string
		line       string
		portInt    int64
		state      bool
		port       int64
		buf        *bufio.Reader
		f          *os.File
	)
	ports = make(map[int64]*Online)
	f, err := os.Open("/proc/net/tcp")
	if err != nil {
		ch <- ""
		return
	}
	buf = bufio.NewReaderSize(f, 16384)
	defer f.Close()
	for {
		line, err = buf.ReadString(10)
		if state == false {
			state = true
			continue
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		line_array = removeEmpty(strings.Split(strings.TrimSpace(line), " "))
		if line_array[3] != "01" {
			continue
		}
		port = hexToDec(strings.Split(line_array[1], ":")[1])
		if _, ok := ports[port]; !ok {
			ports[port] = &Online{}
		}
		ports[port].Online += 1

	}
	f.Close()
	for _, v := range games {
		if v.PortType == "0" {
			out.WriteString(v.GameNote)
			out.WriteString("___online___")
			portInt, _ = strconv.ParseInt(v.Port, 10, 64)
			if _, ok := ports[portInt]; !ok {
				ports[portInt] = &Online{}
			}
			out.WriteString(strconv.Itoa(ports[portInt].Online))
			out.WriteByte(10)
		}
	}

	ch <- out.String()

}

func getInfo(ch chan string, game Game, md5List map[string]GameMd5) {
	var (
		md5Ch    chan string
		err      error
		info     os.FileInfo
		str      bytes.Buffer
		fileTime string
		md5Val   string
	)
	info, err = os.Stat(game.Path)
	if err != nil {
		ch <- ""
		return
	}
	fileTime = info.ModTime().Format("2006-01-02 15:04:05")
	str.WriteString(game.GameNote)
	str.WriteString("___time___")
	str.WriteString(fileTime)
	str.WriteByte(10)
	if md5List[game.GameNote].Time == fileTime && len(md5List) > 0 {
		str.WriteString(game.GameNote)
		str.WriteString("___md5___")
		md5Val = md5List[game.GameNote].Md5
		str.WriteString(md5Val)
		str.WriteByte(10)
	} else {
		md5Ch = make(chan string)
		go md5Sum2(md5Ch, game.Path)
	}
	if md5List[game.GameNote].Time != fileTime {
		str.WriteString(game.GameNote)
		str.WriteString("___md5___")
		md5Val = <-md5Ch
		str.WriteString(md5Val)
		str.WriteByte(10)
	}
	gameMd5File.WriteString(game.GameNote + "___" + md5Val + "___" + fileTime + "\n")
	ch <- str.String()
}

func hexToDec(h string) int64 {
	d, _ := strconv.ParseInt(h, 16, 32)
	return d
}

func removeEmpty(array []string) []string {
	var new_array []string
	for _, i := range array {
		if i != "" && len(new_array) < 4 {
			new_array = append(new_array, i)
		}
	}
	return new_array
}
func md5Sum2(ch chan string, path string) {
	var (
		md5 string
	)
	md5, _ = ctl.Md5sum(path)
	ch <- md5
}
