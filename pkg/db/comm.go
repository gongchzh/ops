package db

import (
	"bytes"
)

type HtmlFile struct {
	Data     bytes.Buffer
	FileName string
	Size     int64
	Dst      string
	PathName string
}
