//解析Json数据
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
)

//定义与json对应的结构体，数组对应slice，字段名对应Json里面的Key
type Server struct {
	ServerName string `json:"serverName"`
	ServerIP   string `json:"serverIp"`
}

type Serverslice struct {
	Servers []Server `json:"servers"`
}

func Run() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	/*
		解析匹配过程
		1. 首先查找tag含有Foo的可导出的struct字段(首字母大写)
		2. 其次查找字段名是Foo的导出字段
		3. 最后查找类似FOO或者FoO这样的除了首字母之外其他大小写不敏感的导出字段

		*能够被赋值的字段必须是可导出字段(即首字母大写）。同时JSON解析的时候只会解析能找得到的字段，找不到的字段会被忽略*
		*当你接收到一个很大的JSON数据结构而你却只想获取其中的部分数据的时候，你只需将你想要的数据对应的字段名大写，即可轻松解决这个问题。*
	*/
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)

	//解析到interface
	/*
		主要解析未知的Json结构
		JSON包中采用map[string]interface{}和[]interface{}结构来存储任意的JSON对象和数组。
		1. bool 代表 JSON booleans,
		2. float64 代表 JSON numbers,
		3. string 代表 JSON strings,
		4. nil 代表 JSON null.
	*/
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f)

	//通过断言访问
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array")
			for i, u := range vv {
				fmt.Println("[", i, "]", u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	//第三方包解析
	js, err := simplejson.NewJson([]byte(`{
	"test": {
		"array": [1, "2", 3],
		"int": 10,
		"float": 5.150,
		"bignum": 9223372036854775807,
		"string": "simplejson",
		"bool": true
		}
	}`))

	arr, _ := js.Get("test").Get("array").Array()
	i, _ := js.Get("test").Get("int").Int()
	ms := js.Get("test").Get("string").MustString()

	fmt.Println(arr)
	fmt.Println(i)
	fmt.Println(ms)

	//生成Json
	//func Marshal(v interface{}) ([]byte, error)
	/*
		定义struct tag的时候需要注意的几点:
		1.字段的tag是"-"，那么这个字段不会输出到JSON
		2.tag中带有自定义名称，那么这个自定义名称会出现在JSON的字段名中，例如上面例子中serverName
		3.tag中如果带有"omitempty"选项，那么如果该字段值为空，就不会输出到JSON串中
		4.如果字段类型是bool, string, int, int64等，而tag中带有",string"选项，那么这个字段在输出到JSON的时候会把该字段对应的值转换成JSON字符串
	*/
	var s2 Serverslice
	s2.Servers = append(s2.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s2.Servers = append(s2.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err = json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
}
