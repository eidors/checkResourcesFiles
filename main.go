package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/klarkxy/gohtml"
)

const (
	ResourcesFILE         = "./checkFolder/Test.resx"
	DictionaryFILE        = "./checkFolder/Dictionary.txt"
	DictionaryDataJSON    = "./checkFolder/Dictionary.json"
	ErrorlogFILE          = "./checkFolder/errorLog.txt"
	EnglishDictFILE       = "./checkFolder/EnglishDict.txt"
	EnglishDictJSONFILE   = "./checkFolder/EnglishDict.json"
	HBWsSiteResourcesFILE = "./checkFolder/hb-ws/"
	WsSiteResourcesFILE   = "./checkFolder/ws/"
)

func main() {
	//没有字典可用的我 没法子只能自己造轮子了 formatDictionary() 就是造轮子的方法
	//formatDictionary()
	//addEnglishWorldList()
	//轮子造好了 就是用了
	// CheckWorlds()

	s := []string{}
	s, _ = GetFormatAllFile(HBWsSiteResourcesFILE, WsSiteResourcesFILE, s)
	fmt.Printf("the number of file is %v,content is:%v\n", len(s), s)
}

func FormatWsWebFile() {
	//open xml
	Rfile, err := os.Open(WsSiteResourcesFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Rfile.Close()
}

func GetFormatAllFile(pathname string, savePathName string, s []string) ([]string, error) {
	fromSlash := filepath.FromSlash(pathname)
	//fmt.Println(fromSlash)
	rd, err := ioutil.ReadDir(fromSlash)
	if err != nil {
		//log.LOGGER("Error").Error("read dir fail %v\n", err)
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		fileSuffix := path.Ext(fi.Name())                         //获取文件后缀
		filenameOnly := strings.TrimSuffix(fi.Name(), fileSuffix) //获取文件名
		if fi.IsDir() {
			fullDir := filepath.Join(fromSlash, fi.Name())
			s, err = GetFormatAllFile(fullDir, savePathName, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			strfullName := filepath.Join(fromSlash, fi.Name())
			s = append(s, strfullName)
			aryWorldList := strings.Split(string(fi.Name()), "-")
			strFloderName := aryWorldList[0]
			exist, err := PathExists(savePathName + strFloderName)
			if err != nil {
				fmt.Printf("get dir error![%v]n", err)
				return s, err
			}
			src, err1 := os.Open(strfullName)
			if err1 != nil {
				fmt.Println(err1)
			}
			defer src.Close()
			if exist {
				dst, err := os.OpenFile(savePathName+strFloderName+"/"+filenameOnly+".tpl", os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					fmt.Println(err)
				}
				defer dst.Close()
				_, err = io.Copy(dst, src)
				if err != nil {
					fmt.Println(err)
				}
				bootstrap := gohtml.NewHtml()
				fmt.Println(bootstrap.String())
			} else {
				err := os.Mkdir(savePathName+strFloderName, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir failed![%v]", err)
				} else {
					dst, err := os.OpenFile(savePathName+strFloderName+"/"+filenameOnly+".tpl", os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						fmt.Println(err)
					}
					defer dst.Close()
					_, err = io.Copy(dst, src)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
	return s, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetAllFile(pathname string, s []string) ([]string, error) {
	fromSlash := filepath.FromSlash(pathname)
	//fmt.Println(fromSlash)
	rd, err := ioutil.ReadDir(fromSlash)
	if err != nil {
		//log.LOGGER("Error").Error("read dir fail %v\n", err)
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := filepath.Join(fromSlash, fi.Name())
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				//log.LOGGER("Error").Error("read dir fail: %v\n", err)
				return s, err
			}
		} else {
			fullName := filepath.Join(fromSlash, fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}

type Recurlyservers struct {
	XMLName     xml.Name `xml:"root"`
	Version     string   `xml:"version,attr"`
	Svs         []data   `xml:"data"`
	Description string   `xml:",innerxml"`
}

type data struct {
	Value string `xml:"value"`
}

func CheckWorlds() {
	ErrorLogFILE, err := os.OpenFile(ErrorlogFILE, os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("open ErrorLogFILE file error")
	}
	defer ErrorLogFILE.Close()

	//open xml
	Rfile, err := os.Open(ResourcesFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Rfile.Close()
	Rdata, _ := ioutil.ReadAll(Rfile)

	//open Json Dictionary.json
	Dfile, err := os.Open(EnglishDictJSONFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	Ddata, _ := ioutil.ReadAll(Dfile)
	// fmt.Println(Ddata)

	m := make(map[string]string)
	err = json.Unmarshal(Ddata, &m)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(m)

	v := Recurlyservers{}
	err = xml.Unmarshal(Rdata, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	// fmt.Println(v.Svs[0].Value)
	srp := strings.NewReplacer("0", " ", "1", " ", "2", " ", "3", " ", "4", " ", "5", " ", "6", " ", "7", " ", "8", " ", "9", " ", "CompAnalyst", " ", ".", " ", "\\", " ", ",", " ")
	for _, val := range v.Svs {
		strLabel := val.Value
		// fmt.Println(strLabel)
		strLabel = srp.Replace(strLabel)
		strLabel = compressStr(strLabel)
		strings.Split(strLabel, " ")
		for _, strWorld := range strings.Split(strLabel, " ") {
			if strWorld == "" {
				continue
			}
			if _, ok := m[strings.ToLower(strWorld)]; !ok {
				//strErrorMessage := "now datetime:%v\n" = m[strWorld] + "\n"
				strErrorMessage := "--Now datetime: " + time.Now().Format("2006-01-02 15:04:05") + "\n"
				strErrorMessage += strWorld + "\n"
				strErrorMessage += val.Value + "\n"
				_, err := ErrorLogFILE.Write([]byte(strErrorMessage))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

//利用正则表达式压缩字符串，去除空格或制表符
func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile(`\\s+`)
	return reg.ReplaceAllString(str, " ")
}

func formatDictionary() {
	//open Dictionary
	Dfile, err := os.Open(DictionaryFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	len, _ := Dfile.Seek(0, 2) //获取文件长度
	strbyte := make([]byte, len)
	Dfile.Seek(0, 0)    //移回指针到文件开头
	Dfile.Read(strbyte) //读文件
	aryWorldList := strings.Split(string(strbyte), "\r\n")
	DictionaryMap := make(map[string]string)
	for _, value := range aryWorldList {
		aryContent := strings.Split(value, "      ")
		// fmt.Println(len(aryContent))
		// fmt.Println(aryContent[0])
		// fmt.Println(index)
		if cap(aryContent) > 1 {
			DictionaryMap[aryContent[0]] = aryContent[1]
		}
		// DictionaryMap["name"]=aryContent[0]
		// DictionaryMap["value"]=aryContent[1]
	}
	if data, err := json.Marshal(DictionaryMap); err == nil {
		// fmt.Printf("%s\n", data)
		// d1 := []byte("hello\ngo\n")
		error := ioutil.WriteFile(DictionaryDataJSON, data, 0644)
		check(error)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addEnglishWorldList() {
	Dfile, err := os.Open(EnglishDictFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	len, _ := Dfile.Seek(0, 2) //获取文件长度
	strbyte := make([]byte, len)
	Dfile.Seek(0, 0)    //移回指针到文件开头
	Dfile.Read(strbyte) //读文件
	aryWorldList := strings.Split(string(strbyte), "\r\n")
	DictionaryMap := make(map[string]string)
	for _, value := range aryWorldList {
		if _, ok := DictionaryMap[value]; !ok {
			DictionaryMap[value] = ""
		}
	}
	if data, err := json.Marshal(DictionaryMap); err == nil {
		// fmt.Printf("%s\n", data)
		// d1 := []byte("hello\ngo\n")
		error := ioutil.WriteFile(EnglishDictJSONFILE, data, 0644)
		check(error)
	}
}
