package ctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func GetJson1(path string) (interface{}, error) {
	var (
		err          error
		data         interface{}
		jsonStr, buf bytes.Buffer
	)

	regStart, _ := regexp.Compile("^[ \t]+")
	regEnd, _ := regexp.Compile("^[ \t]+$")
	regIgnore, _ := regexp.Compile("//.*")
	regIgnore1, _ := regexp.Compile("/\\*.*\\*/")
	err = ReadLine(path, func(line string) error {
		//fmt.Println("line", line)
		line = strings.Replace(line, string([]byte{239, 187, 191}), "", 1)
		line = strings.Replace(line, "\r\n", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, string(9), "", -1)
		line = strings.Replace(line, string(32), "", -1)
		line = regIgnore.ReplaceAllString(line, "")
		line = regStart.ReplaceAllString(line, "")
		line = regEnd.ReplaceAllString(line, "")
		line = strings.Replace(line, "'", "\"", -1)
		line = strings.Replace(line, "\": ", "\":", -1)
		buf.WriteString(line)
		return nil
	})

	jsonStr.WriteString(regIgnore1.ReplaceAllString(buf.String(), ""))
	err = json.Unmarshal(jsonStr.Bytes(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetJson3(path string) (interface{}, error) {
	var (
		err          error
		data         interface{}
		jsonStr, buf bytes.Buffer
	)

	//regStart, _ := regexp.Compile("^[ \t]+")
	//regEnd, _ := regexp.Compile("^[ \t]+$")
	regIgnore, _ := regexp.Compile("//.*")
	regIgnore1, _ := regexp.Compile("/\\*.*\\*/")
	err = ReadLine(path, func(line string) error {
		/*
			line = strings.Replace(line, string([]byte{239, 187, 191}), "", 1)
				line = strings.Replace(line, "\r\n", "", -1)
				line = strings.Replace(line, "\n", "", -1)
				line = strings.Replace(line, string(9), "", -1)
				line = strings.Replace(line, string(32), "", -1)  */
		line = regIgnore.ReplaceAllString(line, "")
		/*		line = regStart.ReplaceAllString(line, "")
				line = regEnd.ReplaceAllString(line, "")
				line = strings.Replace(line, "'", "\"", -1)
				line = strings.Replace(line, "\": ", "\":", -1)  */
		buf.WriteString(line)
		return nil
	})
	jsonStr.WriteString(regIgnore1.ReplaceAllString(buf.String(), ""))
	fmt.Println(jsonStr.String())
	err = json.Unmarshal(jsonStr.Bytes(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetJson(path string) (interface{}, error) {
	var (
		err  error
		b    []byte
		data interface{}
	)
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err

}
func GetJson2(path string) (map[string]interface{}, error) {
	var (
		err                                                             error
		data, bracketMap                                                map[string]interface{}
		key1                                                            string
		jsonStr, bracketBuf, buf, value                                 bytes.Buffer
		key2, key3, key4, key5                                          string
		bracketLevel, bracketType, bracketKeyType, level, nextType, num int
		//bracketNum             int
		sliceString []string
		sliceInt    []int
		sliceMap    []map[string]interface{}
		//sliceMap1                                                                        map[string][]map[string]interface{}
		bracketKey, bracketValue                                                         string
		isNext, quo, isNum, isKey, bracketQuo, bracketNext, isBracketKey, BracketKeyNext bool
	)

	regStart, _ := regexp.Compile("^[ \t]+")
	regEnd, _ := regexp.Compile("^[ \t]+$")
	regNum, _ := regexp.Compile("[0-9]")
	regIgnore, _ := regexp.Compile("//.*")
	regIgnore1, _ := regexp.Compile("/\\*.*\\*/")
	data = make(map[string]interface{})
	err = ReadLine(path, func(line string) error {
		line = strings.Replace(line, string([]byte{239, 187, 191}), "", 1)
		line = strings.Replace(line, "\r\n", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, string(9), "", -1)
		line = strings.Replace(line, string(32), "", -1)
		line = regIgnore.ReplaceAllString(line, "")
		line = regStart.ReplaceAllString(line, "")
		line = regEnd.ReplaceAllString(line, "")
		line = strings.Replace(line, "'", "\"", -1)
		line = strings.Replace(line, "\": ", "\":", -1)
		buf.WriteString(line)
		return nil
	})
	jsonStr.WriteString(regIgnore1.ReplaceAllString(buf.String(), ""))

	isKey = true
	for _, v := range jsonStr.Bytes() {
		if bracketLevel > 0 {
			//		fmt.Println("vvvv", v, "bracketType", bracketType, "isBracketKey", isBracketKey, "bracketKey", bracketKey, "bracketValue", bracketValue, "bracketKeyType", bracketKeyType)
			if bracketNext {
				bracketKey = ""
				bracketValue = ""
				//		fmt.Println("bracketNext", v)
				if regNum.Match([]byte{v}) {
					bracketType = 0
					bracketNext = false
				} else if v == 34 {
					bracketType = 1
					bracketNext = false
					bracketQuo = true
					continue
				} else if v == 91 {
					bracketNext = true
					bracketLevel += 1
					continue
				} else if v == 123 {
					bracketNext = false
					bracketMap = make(map[string]interface{})
					isBracketKey = true
					bracketType = 2
				}
			}
			if bracketLevel == 1 {
				if bracketType == 0 && (v == 93 || v == 44) {
					num, err = strconv.Atoi(bracketBuf.String())
					if err != nil {
						return nil, err
					}
					sliceInt = append(sliceInt, num)
					bracketBuf.Reset()
				} else if v == 34 && bracketType == 1 {
					if bracketQuo == true {

						sliceString = append(sliceString, bracketBuf.String())
						bracketBuf.Reset()
					}
					//			fmt.Println("brackkey", bracketKey, "slice string", sliceString, "lenslice", len(sliceString), "slice0", sliceString[0])
					bracketQuo = !bracketQuo
				} else if bracketType == 2 {
					if v == 123 {
						//		fmt.Println("brackMap2")
						bracketMap = make(map[string]interface{})
						isBracketKey = true
						continue
					}
					if BracketKeyNext {
						if v == 34 {
							bracketKeyType = 34
						} else if regNum.Match([]byte{v}) {
							bracketKeyType = 0
						}
						BracketKeyNext = false
					}
					if v == 34 {
						bracketQuo = !bracketQuo
						if !bracketQuo {

							if !isBracketKey && bracketKeyType == 34 {
								//sliceMap = append(sliceMap, map[string]interface{}{bracketKey: bracketValue})
								bracketMap[bracketKey] = bracketValue
								bracketKey = ""
								bracketValue = ""
							} else if isBracketKey {
								isBracketKey = false
							}

						}
						continue
					}

					if v == 58 && !bracketQuo {
						isBracketKey = false
						BracketKeyNext = true
						continue
					}
					if v == 44 {
						if bracketKeyType == 0 && !isBracketKey {
							num, err = strconv.Atoi(bracketValue)
							if err != nil {
								fmt.Println("errr", err.Error())
								return nil, err
							}
							bracketMap[bracketKey] = num
							bracketKey = ""
							bracketValue = ""
						}
						isBracketKey = true

						continue
					}
					if v == 125 {
						if bracketKeyType == 0 {
							num, err = strconv.Atoi(bracketValue)
							if err != nil {
								return nil, err
							}
							bracketMap[bracketKey] = num
							bracketKey = ""
							bracketValue = ""
							isBracketKey = true
						}
						sliceMap = append(sliceMap, bracketMap)
						//		fmt.Println("sliceMap", sliceMap)
						continue
					}
					if isBracketKey {
						bracketKey += string(v)
					} else {
						bracketValue += string(v)
					}

				} else if v == 93 {

				} else {
					//		fmt.Println("v", v, "buf", bracketBuf.String())
					bracketBuf.WriteByte(v)
				}
			}

			if v == 93 {
				//	fmt.Println("brack level", level)
				switch level {
				case 1:
					if bracketType == 0 {

						data[key1] = sliceInt
						sliceInt = nil
					} else if bracketType == 1 {
						data[key1] = sliceString
						sliceString = nil
					} else if bracketType == 2 {
						data[key1] = sliceMap

					}
				//	fmt.Println("level1", data[key1])
				case 2:
					if bracketType == 0 {
						data[key1].(map[string]interface{})[key2] = sliceInt
					} else if bracketType == 1 {
						//		fmt.Println("slicestring", sliceString)

						//		fmt.Println("key1", key1, "key2", key2)
						/*
							if level == 2 {
								data[key1].(map[string]interface{})[key2] = make(map[string]interface{})
							}
							fmt.Println("type1", data[key1], "key1", key1)
							data[key1].(map[string]interface{})[key2] = []string{"a"}
						*/
						data[key1].(map[string]interface{})[key2] = sliceString
					} else if bracketType == 2 {
						data[key1].(map[string]interface{})[key2] = sliceMap
						sliceMap = nil
					}
				//	fmt.Println("level2", data[key1])
				case 3:
					if bracketType == 0 {
						data[key1].(map[string]map[string]interface{})[key2][key3] = sliceInt
					} else if bracketType == 1 {
						data[key1].(map[string]map[string]interface{})[key2][key3] = sliceString
					} else if bracketType == 2 {
						data[key1] = sliceMap
						sliceMap = nil
					}
				case 4:
					if bracketType == 0 {
						data[key1].(map[string]map[string]map[string]interface{})[key2][key3][key4] = sliceInt
					} else if bracketType == 1 {
						data[key1].(map[string]map[string]map[string]interface{})[key2][key3][key4] = sliceString
					} else if bracketType == 2 {
						data[key1] = sliceMap
						sliceMap = nil
					}
				case 5:
					if bracketType == 0 {
						data[key1].(map[string]map[string]map[string]map[string]interface{})[key2][key3][key4][key5] = sliceInt
					} else if bracketType == 1 {
						data[key1].(map[string]map[string]map[string]map[string]interface{})[key2][key3][key4][key5] = sliceString
					} else if bracketType == 2 {
						data[key1] = sliceMap
						sliceMap = nil
					}
				}
				sliceMap = nil
				sliceInt = nil
				sliceString = nil

				bracketLevel -= 1

			}
			continue
		}
		if true {
			//	fmt.Println("vvvv", v, "level", level, "isKey", isKey, "Key1", key1, "key2", key2)
		}
		if v == 58 && !quo {
			//fmt.Println("iskey1", isKey)
			isKey = false
			isNext = true
			continue
		}
		if v == 123 {
			level += 1
			isKey = true
			switch level {
			case 2:
				data[key1] = make(map[string]interface{})
			}
			//fmt.Println("123{", v)
			continue
		}
		if v == 44 {
			isKey = true
			if isNum {
				num, err = strconv.Atoi(value.String())
				if err != nil {
					return nil, err
				}
				switch level {
				case 1:
					data[key1] = num
				case 2:
					data[key1].(map[string]interface{})[key2] = num
				//	fmt.Println("num", num, "level", level, "key1_value", data[key1])
				case 3:
					data[key1].(map[string]map[string]interface{})[key2][key3] = num
				case 4:
					data[key1].(map[string]map[string]map[string]interface{})[key2][key3][key4] = num
				case 5:
					data[key1].(map[string]map[string]map[string]map[string]interface{})[key2][key3][key4][key5] = num
				}
				value.Reset()
				isNum = false
			}
			continue
		}
		if v == 125 {
			if isNum {
				num, err = strconv.Atoi(value.String())
				if err != nil {
					return nil, err
				}
				switch level {
				case 1:
					data[key1] = num
				case 2:
					data[key1].(map[string]interface{})[key2] = num
				//	fmt.Println("num", num, "level", level, "key1_value", data[key1])
				case 3:
					data[key1].(map[string]map[string]interface{})[key2][key3] = num
				case 4:
					data[key1].(map[string]map[string]map[string]interface{})[key2][key3][key4] = num
				case 5:
					data[key1].(map[string]map[string]map[string]map[string]interface{})[key2][key3][key4][key5] = num
				}
				value.Reset()
				isNum = false
			}
			level -= 1
			continue
		}
		if isNum {
			value.WriteByte(v)
			continue
		}

		if isNext {
			//	fmt.Println("isNext")
			if v == 34 {
				nextType = 34
			} else if v == 123 {
				level += 1
			} else if regNum.Match([]byte{v}) {
				nextType = 1
				value.WriteByte(v)
				isNum = true
			} else if v == 91 {
				nextType = 91
				bracketLevel += 1
				bracketNext = true
			}
			isNext = false
			//continue
		}

		if v == 34 {
			quo = !quo
			if !quo {
				if isKey {
					if level == 1 {
						key1 = value.String()
					}
					if level == 2 {
						key2 = value.String()
					}
					value.Reset()
					isKey = false
				} else {
					if nextType == 34 {
						if level == 1 {
							data[key1] = value.String()
							//	fmt.Println("key1", key1, "value", value)
						}
						if level == 2 {
							//	fmt.Println("key1", key1, "key2", key2, "value", value.String())
							data[key1].(map[string]interface{})[key2] = value.String()
							//	fmt.Println("data", data[key1])
						}
						value.Reset()
					}
				}
			}
			continue
		}
		if quo {
			value.WriteByte(v)
			continue
		}
	}
	return data, err
}
