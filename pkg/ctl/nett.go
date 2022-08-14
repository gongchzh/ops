package ctl

import (
	"io/ioutil"
	"net/http"
)

func GetRemoteIp(url string) (string, error) {
	var (
		ip  string
		err error
	)
	res, err := http.Get(url)
	if err != nil {
		return ip, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ip, err
	}
	if !RegIpOn.Match(b) {
		return ip, Errorf("获取ip异常")
	}
	return string(b), err
}
