package ctl

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
)

const (
	ConstHeader       = "testHeader"
	ConstHeaderLength = 10
	ConstMLength      = 4
)

type Msg struct {
	Meta    map[string]interface{} `json:"meta"`
	Content interface{}            `json:"content"`
}

func Depack(buffer []byte) []byte {
	length := len(buffer)
	var i int
	data := make([]byte, 32)
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstMLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstMLength])
			if length < i+ConstHeaderLength+ConstMLength+messageLength {
				break
			}
			data = buffer[i+ConstHeaderLength+ConstMLength : i+ConstHeaderLength+ConstMLength+messageLength]
		}
	}
	if i == length {
		return make([]byte, 0)
	}
	return data
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

var routers [][2]interface{}

type Controller interface {
	Excute(message Msg) []byte
}

func TaskDeliver(postdata []byte, conn net.Conn) error {
	var (
		err error
	)
	for _, v := range routers {
		pred := v[0]
		act := v[1]
		var entermsg Msg
		err = json.Unmarshal(postdata, &entermsg)
		if err != nil {
			return err
		}
		if pred.(func(entermsg Msg) bool)(entermsg) {
			result := act.(Controller).Excute(entermsg)
			conn.Write(result)
			return err
		}
	}
	return err
}
