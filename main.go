package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
    "strings"
    "encoding/json"
    "regexp"
)

const (
	ResourcesFILE = "./checkFolder/Test.resx"
	DictionaryFILE = "./checkFolder/Dictionary.txt"
	DictionaryDataJSON = "./checkFolder/Dictionary.json"
	ErrorlogFILE = "./checkFolder/errorLog.txt"
	EnglishDictFILE = "./checkFolder/EnglishDict.txt"
	EnglishDictJSONFILE = "./checkFolder/EnglishDict.json"
)

func main() {
	fmt.Println("Hello, World!")
	//没有字典可用的我 没法子只能自己造轮子了 formatDictionary() 就是造轮子的方法
	//formatDictionary()
	addEnglishWorldList()
	//

	//轮子造好了 就是用了
	CheckWorlds()
}

type Recurlyservers struct {
	XMLName     xml.Name `xml:"root"`
	Version     string   `xml:"version,attr"`
	Svs         []data `xml:"data"`
	Description string   `xml:",innerxml"`
}

type data struct {
	Value  string   `xml:"value"`
}

func CheckWorlds() {
	// //os
	// fd, err := os.Open(ResourcesFILE) //打开文件
	// if err != nil {
	// 	fmt.Println("open file error")
	// }
	// defer fd.Close()
	// len, err := fd.Seek(0, 2) //获取文件长度
	// if err != nil {
	// 	fmt.Println("get len error")
	// }
	// strbyte := make([]byte, len)
	// fd.Seek(0, 0)    //移回指针到文件开头
	// fd.Read(strbyte) //读文件
	// fmt.Println(string(strbyte))

	// //io/ioutil
	// dat, err := ioutil.ReadFile(FILE) //直接读文件
	// if err != nil {
	// 	fmt.Println("Read file err")
	// }
	// fmt.Println(string(dat))


	//open log files
	ErrorLogFILE,err := os.OpenFile(ErrorlogFILE,os.O_APPEND,0777)
	if err != nil{
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
	Rdata, err := ioutil.ReadAll(Rfile)

	
	//open Json Dictionary.json
	Dfile, err := os.Open(DictionaryDataJSON) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	Ddata, err := ioutil.ReadAll(Dfile)
	// fmt.Println(Ddata)

	m := make(map[string]string)
    err = json.Unmarshal(Ddata , &m)
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
	srp := strings.NewReplacer("0" ," " , "1" ," " , "2" ," " , "3" ," " , "4" ," " , "5" ," " , "6" ," " , "7" ," " , "8" ," " , "9" ," " ,"CompAnalyst" ," ", "." ," " , "\\" ," " , "," ," ")
	for _, val := range v.Svs {
		strLabel := val.Value
		// fmt.Println(strLabel)
		strLabel = srp.Replace(strLabel)
		strLabel = compressStr(strLabel)
		strings.Split(strLabel," ")
		for _, strWorld := range strings.Split(strLabel," "){
			// fmt.Println(strWorld)
			// m[strWorld]
			// if m[strWorld] == "" {
			// 	strErrorMessage := m[strWorld] + "\n"
			// 	strErrorMessage += val.Value + "\n"
			// 	error := ioutil.WriteFile(ErrorlogFILE, []byte(strErrorMessage), 0644)
			// 	check(error)
			// }

			if _, ok := m[strings.ToLower(strWorld)]; !ok{
				strErrorMessage := m[strWorld] + "\n"
				strErrorMessage += strWorld + "\n"
				strErrorMessage += val.Value + "\n"
				_,err := ErrorLogFILE.Write([]byte(strErrorMessage))
				if err != nil{
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
    reg := regexp.MustCompile("\\s+")
    return reg.ReplaceAllString(str, " ")
}

func formatDictionary(){
	//open Dictionary
	Dfile, err := os.Open(DictionaryFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	len, err := Dfile.Seek(0, 2) //获取文件长度
	strbyte := make([]byte, len)
	Dfile.Seek(0, 0)    //移回指针到文件开头
	Dfile.Read(strbyte) //读文件
	aryWorldList := strings.Split(string(strbyte),"\r\n")
	var DictionaryMap map[string]string /*创建集合 */
	DictionaryMap = make(map[string]string)
	for _, value := range aryWorldList {
		aryContent := strings.Split(value,"      ")
		// fmt.Println(len(aryContent))
		// fmt.Println(aryContent[0])
		// fmt.Println(index)
		if cap(aryContent) > 1 {
			DictionaryMap[aryContent[0]]=aryContent[1]
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

func addEnglishWorldList(){
	Dfile, err := os.Open(EnglishDictFILE) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer Dfile.Close()
	len, err := Dfile.Seek(0, 2) //获取文件长度
	strbyte := make([]byte, len)
	Dfile.Seek(0, 0)    //移回指针到文件开头
	Dfile.Read(strbyte) //读文件
	aryWorldList := strings.Split(string(strbyte),"\r\n")
	var DictionaryMap map[string]string /*创建集合 */
	DictionaryMap = make(map[string]string)
	for _, value := range aryWorldList {
		if _, ok := DictionaryMap[value]; !ok {
			DictionaryMap[value] = "";
		}
	}
    if data, err := json.Marshal(DictionaryMap); err == nil {
		// fmt.Printf("%s\n", data)
		// d1 := []byte("hello\ngo\n")
		error := ioutil.WriteFile(EnglishDictJSONFILE, data, 0644)
		check(error)
    }
}