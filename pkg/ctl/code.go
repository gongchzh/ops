package ctl

import (
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/encoding/gbase64"
)

type Aes struct {
	key []byte
}

func (a *Aes) Encode(s string) string {
	a1, err := gaes.Encrypt([]byte(s), a.key)
	if err != nil {
		Debug(err)
		return ""
	}
	return string(gbase64.Encode(a1))

}

func (a *Aes) Decode(s string) string {
	a1, err := gbase64.Decode([]byte(s))
	if err != nil {
		Debug(err)
		return ""
	}
	a2, err := gaes.Decrypt(a1, a.key)
	if err != nil {
		return ""
	}
	return string(a2)
}

func NewAes(key string) *Aes {
	a := new(Aes)
	a.key = []byte(key)
	return a
}
