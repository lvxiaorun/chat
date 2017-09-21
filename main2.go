package main

import (
	"fmt"
	"time"
	"golang.org/x/net/websocket"
	"net/http"
	"html/template"
)

var allData AllData
var allUser map[string]UserInfo
func main() {
	fmt.Println("启动时间")
	fmt.Println(time.Now())

	//初始化
	allData = AllData{}
	allUser = make(map[string]UserInfo)

	//绑定效果页面
	http.HandleFunc("/", h_index2)
	//绑定socket方法
	http.Handle("/webSocket", websocket.Handler(h_webSocket2))
	//开始监听
	http.ListenAndServe(":8080", nil)
}

func h_index2(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	t,_ := template.ParseFiles("index.html")
	t.Execute(w,name)
	//http.ServeFile(w, r, "index2.html")
}

func h_webSocket2(ws *websocket.Conn){
	var msg Message
	var data string
	for {
		msglen := len(allData.Messages)
		fmt.Println("Msgs", msglen, "allUser长度：", len(allUser))
	}

}
type Message struct {
	UserName string
	Msg      string
	DataType string
	ToUser   string
}

type UserInfo struct {
	UserName string
	Conn *websocket.Conn
}

type AllData struct {
	Messages []Message
	UserInfos []UserInfo
}
