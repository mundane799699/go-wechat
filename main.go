package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/eatmoreapple/openwechat"
)

type Message struct {
	Msg []string `json:"msg"`
}

func contains(groupNames []string, groupName string) bool {
	for _, name := range groupNames {
		if name == groupName {
			return true
		}
	}
	return false
}

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的群组
	groups, err := self.Groups()
	var myGroups []*openwechat.Group
	groupNames := []string{"有钱才算自由", "大厂交流群", "小程序学习交流群"}
	for _, group := range groups {
		if contains(groupNames, group.NickName) {
			myGroups = append(myGroups, group)
		}
	}
	fmt.Println(myGroups)

	// /hello 路径的处理函数
	helloHandler := func(w http.ResponseWriter, r *http.Request) {

		// 检查请求方法是否为POST
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析JSON数据
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 处理接收到的消息
		text := strings.Join(message.Msg, "\n\n")
		for _, group := range myGroups {
			_, err = group.SendText(text)
		}
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("Message received successfully"))
		}
	}

	// 注册处理函数并启动服务器
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("11111111111111111111111111")
	http.ListenAndServe(":8089", nil)
	fmt.Println("22222222222222222222222222")

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
