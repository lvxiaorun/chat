package main

import (
	"fmt"
	"time"
	"golang.org/x/net/websocket"
	"net/http"
	"html/template"
	"encoding/json"
	"strings"
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
	t, _ := template.ParseFiles("index2.html")
	t.Execute(w, name)
	//http.ServeFile(w, r, "index2.html")
}
//
//func h_webSocket2Login(ws *websocket.Conn){
//	var name string
//	err := websocket.Message.Receive(ws,&name)
//	if err != nil{
//		fmt.Println("Login err name:",name,"err:",err)
//		return
//	}
//}
func h_webSocket2(ws *websocket.Conn) {
	var msg Message
	var data string
J:
	for {
		msglen := len(allData.Messages)
		fmt.Println("Msgs", msglen, "allUser长度：", len(allUser))
		//有消息时,判断发给个人还是群组发送
		if msglen > 0 {
			//b,err := json.Marshal(allData)
			//if err != nil{
			//	fmt.Println("全局信息异常:",err)
			//	break
			//}
			for key, value := range allUser {
				partdata := AllData{}
				partdata.UserInfos = allData.UserInfos
				for _, item := range allData.Messages {
					if item.ToUser == key || item.ToUser == "all" {
						partdata.Messages = append(partdata.Messages, item)
					}
				}
				b, errh := json.Marshal(partdata)
				if errh != nil {
					fmt.Println("序列化用户全局消息失败,user:", key, ":", value)
					break J
				}
				errsend := websocket.Message.Send(value.Conn, string(b))
				if errsend != nil{
					fmt.Println("发送消息失败,user:",key,":",value)
					delete(allUser,key)
					break J
				}
			}
			allData.Messages = make([]Message,0)
		}
		fmt.Println("开始解析数据")
		errr := websocket.Message.Receive(ws,&data)
		if errr != nil{
			for key,value := range allUser{
				if value.Conn == ws{
					delete(allUser,key) //删除错误的连接
				}
			}
			fmt.Println("接收消息失败")
			break J
		}
		fmt.Println("data:",data)
		data = strings.Replace(data, "\n", "", 0)
		errby := json.Unmarshal([]byte(data),&msg)
		if errby != nil{
			fmt.Println("解析数据出错:",errby)
			break J
		}
		fmt.Println("数据类型:",msg.DataType)
		allData.Messages = append(allData.Messages,msg)
		switch msg.DataType {
		case "login":
			fmt.Println(msg.UserName+"上线")
			allData.UserInfos = append(allData.UserInfos,UserInfo{msg.UserName,ws})
			allUser[msg.UserName] = UserInfo{msg.UserName,ws}
		}
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
	Conn     *websocket.Conn
}

type AllData struct {
	Messages  []Message
	UserInfos []UserInfo
}
