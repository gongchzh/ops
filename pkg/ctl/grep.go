package ctl

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

func GrepMatch(path interface{}, reg string) []string {
	var (
		regPath1, regPath2, regStr *regexp.Regexp
		isPath                     bool
		tmpByte                    [][]byte
		grepStr                    []string
	)
	regPath1, _ = regexp.Compile("\\\\")
	regPath2, _ = regexp.Compile("/")
	regStr, _ = regexp.Compile(reg)
	switch value := path.(type) {
	case string:
		if len(strings.Split(value, "\n")) <= 1 {
			if regPath1.MatchString(value) || regPath2.MatchString(value) {
				if CheckFile(value) {
					isPath = true
				}
			}
		}

		if isPath {

			f, err := os.Open(value)
			if err != nil {
				return nil
			}
			buf := bufio.NewReader(f)

			for {
				line, err := buf.ReadBytes(10)
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil
				}
				if regStr.Match(line) {
					//grepStr = append(grepStr, regStr.FindAllString(string(line), -1))
					/*for _, v := range regStr.FindAll(line, -1) {
						grepStr = append(grepStr, string(v))
					}*/
					grepStr = append(grepStr, regStr.FindString(string(line)))
				}
			}
			return grepStr

		} else {
			return regStr.FindAllString(value, -1)
		}

	case []byte:
		tmpByte = regStr.FindAll(value, -1)
		for _, v := range tmpByte {
			grepStr = append(grepStr, string(v))
		}
		return grepStr
	}
	return nil
}

func GrepAll(path interface{}, reg string) []string {
	var (
		regPath1, regPath2, regStr *regexp.Regexp
		isPath                     bool
		tmpByte                    [][]byte
		grepStr                    []string
	)
	regPath1, _ = regexp.Compile("\\\\")
	regPath2, _ = regexp.Compile("/")
	regStr, _ = regexp.Compile(reg)
	switch value := path.(type) {
	case string:
		if len(strings.Split(value, "\n")) <= 1 {
			if regPath1.MatchString(value) || regPath2.MatchString(value) {
				if CheckFile(value) {
					isPath = true
				}
			}
		}

		if isPath {

			f, err := os.Open(value)
			if err != nil {
				return nil
			}
			buf := bufio.NewReader(f)

			for {
				line, err := buf.ReadBytes(10)
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil
				}
				if regStr.Match(line) {
					for _, v := range regStr.FindAll(line, -1) {
						grepStr = append(grepStr, string(v))
					}
				}
			}
			return grepStr

		} else {
			return regStr.FindAllString(value, -1)
		}

	case []byte:
		tmpByte = regStr.FindAll(value, -1)
		for _, v := range tmpByte {
			grepStr = append(grepStr, string(v))
		}
		return grepStr
	}
	return nil
}

func Grep(path interface{}, reg string) []string {
	var (
		regPath1, regPath2, regStr *regexp.Regexp
		isPath                     bool
		grepStr                    []string
	)
	regPath1, _ = regexp.Compile("\\\\")
	regPath2, _ = regexp.Compile("/")
	regStr, _ = regexp.Compile(reg)
	switch value := path.(type) {
	case string:
		if len(strings.Split(value, "\n")) <= 1 {
			if regPath1.MatchString(value) || regPath2.MatchString(value) {
				if CheckFile(value) {
					isPath = true
				}
			}
		}

		if isPath {
			f, err := os.Open(value)
			if err != nil {
				return nil
			}
			buf := bufio.NewReader(f)
			for {
				line, err := buf.ReadBytes('\n')

				if err != nil {
					if err == io.EOF {
						break
					}
					return nil
				}
				if regStr.Match(line) {
					grepStr = append(grepStr, string(line))
				}

			}
			return grepStr
		} else {
			for _, line := range strings.Split(value, "\n") {
				if regStr.MatchString(line) {
					grepStr = append(grepStr, line)
				}
			}
			return grepStr
		}
	case []byte:
		for _, line := range bytes.Split(value, []byte("\n")) {
			if regStr.Match(line) {
				grepStr = append(grepStr, string(line))
			}
		}
		return grepStr
	}
	return nil
}

func Fgrep(path interface{}, reg string) []string {
	var (
		regPath1, regPath2 *regexp.Regexp
		isPath             bool
		grepStr            []string
		regByte            []byte
	)
	regPath1, _ = regexp.Compile("\\\\")
	regPath2, _ = regexp.Compile("/")
	regByte = []byte(reg)
	if false {
		bytes.Split([]byte("abcd"), regByte)
	}
	switch value := path.(type) {
	case string:
		if len(strings.Split(value, "\n")) <= 1 {
			if regPath1.MatchString(value) || regPath2.MatchString(value) {
				if CheckFile(value) {
					isPath = true
				}
			}
		}
		if isPath {
			f, err := os.Open(value)
			if err != nil {
				return nil
			}
			buf := bufio.NewReader(f)
			for {
				line, err := buf.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil
				}
				if bytes.Contains(line, regByte) {
					grepStr = append(grepStr, string(line))
				}
			}
			return grepStr
		} else {
			for _, line := range strings.Split(value, "\n") {
				if strings.Contains(line, reg) {
					grepStr = append(grepStr, line)
				}
			}
			return grepStr
		}
	case []byte:
		for _, line := range bytes.Split(value, []byte("\n")) {
			if bytes.Contains(line, regByte) {
				grepStr = append(grepStr, string(line))
			}
		}
		return grepStr
	}
	return nil
}
