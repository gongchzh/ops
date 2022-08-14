package ctl

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

type Json struct {
	Map    map[string]interface{}
	Path   string
	KeyErr error
}

func JsonSet(file string, key1, key2 string, value interface{}, location int) error {
	var (
		err error
		nf  string
		ln  int
		end string
	)
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return Errorf("%s,%s,%s", err.Error(), key1, key2)
	}
	regKey, err := regexp.Compile(`^[ |\t]*"` + key1 + `":`)
	for _, v := range strings.Split(string(f), "\n") {
		if ln == location && location != 0 {
			if v[len(v)-1] == 13 {
				end = string(13)
			}
			if !strings.Contains(v, key2) {
				//	Debug(v)
				return Errorf("配置文件参数位置异常:%s,%s", key1, key2)
			}
			vn := strings.Split(v, ":")
			if len(vn) != 2 {
				return Errorf("配置文件key格式异常:%s,%s", key1, key2)
			}
			switch v1 := value.(type) {
			case string:
				nf += vn[0] + ": \"" + v1 + "\""
			case int:
				nf += vn[0] + ": " + Itoa(v1) + ""
			}
			if strings.Contains(v, ",") {
				nf += ","
			}
			nf += end + "\n"
			//	ctl.Debug(vn)
			//	ctl.Debug(nf)
			ln++
			continue
		}
		if regKey.MatchString(v) {
			if key2 == "" {
				if v[len(v)-1] == 13 {
					end = string(13)
				}
				vn := strings.Split(v, ":")
				if len(vn) != 2 {
					return Errorf("配置文件key格式异常:%s,%s", key1, key2)
				}
				switch v1 := value.(type) {
				case string:
					nf += vn[0] + ": \"" + v1 + "\""
				case int:
					nf += vn[0] + ": " + Itoa(v1) + ""
				}
				if strings.Contains(v, ",") {
					nf += ","
				}
				nf += end + "\n"
				ln = 100000
				continue
			} else {
				ln++
			}
		} else if ln != 0 {
			ln++
		}
		nf += v + "\n"
	}
	if nf[:len(nf)-1] != string(13) {
		nf = nf[:len(nf)-1]
	}
	if ln == 0 {
		return Errorf("获取配置文件key异常:%s,%s", key1, key2)
	}
	err = ioutil.WriteFile(file, []byte(nf), 0644)
	return err
}

func RemvJsAnno(b []byte) []byte {
	var (
		n      string
		isAnno bool
	)
	for _, v := range strings.Split(string(b), "\n") {
		if isAnno {
			if strings.Contains(v, "*/") {
				n += strings.Split(v, "*/")[1]
				isAnno = false
			} else {
				continue
			}
		} else {
			if strings.Contains(v, "/*") {
				n += strings.Split(v, "/*")[0]
				isAnno = true
			} else {
				if strings.Contains(v, "//") {
					n += strings.Split(v, "//")[0]
				} else {
					n += v
				}
			}
		}

		n += "\n"
	}
	return []byte(n)
}

func (j *Json) ChangeValue(key1, key2 string, vn interface{}) error {
	var (
		vo interface{}
	)
	if key2 != "" {
		vo = j.Get(key1)
	} else {
		vo = j.GetTwoKey(key1, key2)
	}
	if j.KeyErr != nil {
		return j.KeyErr
	}
	switch value := vo.(type) {
	case string:
		Debug(value)
	case int:
	case float64:
	}
	return nil

}

func (j *Json) IsKeyError() bool {
	if j.KeyErr != nil {
		j.KeyErr = nil
		return true
	}
	return false
}

func (j *Json) IsNotKey(key string) bool {
	if _, ok := j.Map[key]; ok {
		j.KeyErr = nil
		return false
	}
	j.KeyErr = Errorf("key is null")
	return true
}

func (j *Json) GetTwoKey(key1, key2 string) interface{} {
	if j.IsNotKey(key1) {
		return nil
	}
	switch reflect.TypeOf(j.Get(key1)).String() {
	case "map[string]interface {}":
		if _, ok := j.Get(key1).(map[string]interface{})[key2]; ok {
			if reflect.TypeOf(j.Get(key1).(map[string]interface{})[key2]).String() == "string" {
				return j.Get(key1).(map[string]interface{})[key2]
			}
		}
	case "map[string]string":
		if _, ok := j.Get(key1).(map[string]string)[key2]; ok {
			return j.Get(key1).(map[string]string)[key2]
		}
	case "map[string]int":
		if _, ok := j.Get(key1).(map[string]int)[key2]; ok {
			return j.Get(key1).(map[string]int)[key2]
		}
	case "map[string]float64":
		if _, ok := j.Get(key1).(map[string]float64)[key2]; ok {
			return j.Get(key1).(map[string]float64)[key2]
		}
	case "map[string]float32":
		if _, ok := j.Get(key1).(map[string]float32)[key2]; ok {
			return j.Get(key1).(map[string]float32)[key2]
		}
	}

	j.KeyErr = Errorf("key type error")
	return nil
}

func (j *Json) Get(key string) interface{} {
	if j.IsNotKey(key) {
		return nil
	}
	return j.Map[key]
}

func (j *Json) GetString(key string) string {
	if j.IsNotKey(key) {
		return ""
	}
	if reflect.TypeOf(j.Map[key]).String() != "string" {
		j.KeyErr = Errorf("key type error")
		return ""
	}
	return j.Map[key].(string)
}

func (j *Json) GetInt(key string) int {
	if j.IsNotKey(key) {
		return 0
	}
	if reflect.TypeOf(j.Map[key]).String() != "int" {
		j.KeyErr = Errorf("key type error")
		return 0
	}
	return j.Map[key].(int)
}

func (j *Json) GetFloat(key string) float64 {
	if j.IsNotKey(key) {
		return 0
	}
	if reflect.TypeOf(j.Map[key]).String() != "float64" {
		j.KeyErr = Errorf("key type error")
		return 0
	}
	return j.Map[key].(float64)
}

func (j *Json) GetKeyString(key1, key2 string) string {
	var (
		value interface{}
	)
	value = j.GetTwoKey(key1, key2)
	if value != nil && reflect.TypeOf(value).String() == "string" {
		return value.(string)
	}
	j.KeyErr = Errorf("key is error")
	return ""
}

func (j *Json) GetKeyInt(key1, key2 string) int {
	var (
		value interface{}
	)
	value = j.GetTwoKey(key1, key2)
	if value != nil && reflect.TypeOf(value).String() == "int" {
		return value.(int)
	}
	j.KeyErr = Errorf("key is error")
	return 0
}

func (j *Json) GetKeyFloat(key1, key2 string) float64 {
	var (
		value interface{}
	)
	value = j.GetTwoKey(key1, key2)
	if value != nil && reflect.TypeOf(value).String() == "float64" {
		return value.(float64)
	}
	j.KeyErr = Errorf("key is error")
	return 0
}

func (j *Json) GetData() ([]byte, error) {
	var (
		err      error
		nf1, nf2 string
		isAnn    bool
	)
	b, err := ioutil.ReadFile(j.Path)
	if err != nil {
		return nil, err
	}
	for _, v := range strings.Split(string(b), "\n") {
		if v == "" {
			continue
		}
		if strings.Contains(v, "//") {
			nf1 += strings.Split(v, "//")[0] + "\n"
		} else {
			nf1 += v + "\n"
		}
	}
	for _, v := range strings.Split(nf1, "\n") {
		if v == "" {
			continue
		}
		if strings.Contains(v, "/*") {
			isAnn = true
			nf2 += strings.Split(v, "/*")[0] + "\n"
			continue
		}
		if isAnn && strings.Contains(v, "*/") {
			isAnn = false
			nf2 += strings.Split(v, "*/")[1] + "\n"
			continue
		}
		if isAnn {
			continue
		}
		nf2 += v + "\n"
	}
	return []byte(nf2), err
}
func (j *Json) GetMap() error {
	dt, err := j.GetData()
	j.Map = nil
	err = json.Unmarshal(dt, &j.Map)
	return err

}
