package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	//websocket 的升级接口
	upgrader := websocket.Upgrader{}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//通过 upgrader 将 http 连接升级为 websocket 连接
		connect, err := upgrader.Upgrade(writer, request, nil)
		if nil != err {
			log.Println(err)
			return
		}

		defer connect.Close()

		//定时向客户端发送数据
		go tickWriter(connect)

		//启动数据读取循环，读取客户端发送来的数据
		for {
			//从 websocket 中读取数据
			//messageType 消息类型，websocket 标准
			//messageData 消息数据
			messageType, messageData, err := connect.ReadMessage()
			if nil != err {
				log.Println(err)
				break
			}
			switch messageType {
			case websocket.TextMessage: //文本数据
				fmt.Println(string(messageData))
			case websocket.BinaryMessage: //二进制数据
				fmt.Println(messageData)
			case websocket.CloseMessage: //关闭
			case websocket.PingMessage: //Ping
			case websocket.PongMessage: //Pong
			default:

			}
		}
	})

	err := http.ListenAndServe(":60000", nil)
	if nil != err {
		log.Println(err)
		return
	}
}

func tickWriter(connect *websocket.Conn) {
	for {
		//向客户端发送类型为文本的数据
		err := connect.WriteMessage(websocket.TextMessage, []byte("from server to client"))
		if nil != err {
			log.Println(err)
			break
		}

		time.Sleep(time.Second)
	}
}
