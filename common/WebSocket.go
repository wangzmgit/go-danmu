package common

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"kuukaa.fun/danmu-v4/util"
)

var (
	message = make(map[string]chan interface{}) // 消息通道
	client  = make(map[string]*websocket.Conn)  // websocket客户端链接池
	mux     sync.Mutex                          // 互斥锁
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理ws请求
func WsHandler(w http.ResponseWriter, r *http.Request, id string) {
	var conn *websocket.Conn
	var err error
	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 10)
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Logfile(util.ErrorLog, " Failed to set websocket upgrade "+err.Error())
		return
	}
	// 把与客户端的链接添加到客户端链接池中
	addClient(id, conn)
	// 获取该客户端的消息通道
	m, exist := getMessageChannel(id)
	if !exist {
		m = make(chan interface{})
		addMessageChannel(id, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteClient(id)
		// fmt.Println(code)
		return nil
	})

	for {
		select {
		case content, ok := <-m:
			// 从消息通道接收消息，然后推送给前端
			err = conn.WriteJSON(content)
			if err != nil {
				// log.Println(err)
				util.Logfile(util.ErrorLog, " Send MEssage Error"+err.Error())
				if ok {
					go func() {
						m <- content
					}()
				}

				conn.Close()
				deleteClient(id)
				return
			}
		case <-pingTicker.C:
			// 服务端心跳:每20秒ping一次客户端，查看其是否在线
			conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				util.Logfile(util.ErrorLog, " Send ping error "+err.Error())
				conn.Close()
				deleteClient(id)
				return
			}
		}
	}

}

/*********************************************************
** 函数功能: 向用户发送消息
** 日    期: 2022年2月25日09:37:48
**********************************************************/
func SendMsgToUser(id string, msg interface{}) {
	if id != "" {
		_, exist := getMessageChannel(id)
		if !exist {
			return
		}
	}
	setMessage(id, msg)
	return
}

func addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	client[id] = conn
	mux.Unlock()
}

func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = client[id]
	mux.Unlock()
	return
}

func deleteClient(id string) {
	mux.Lock()
	delete(client, id)
	mux.Unlock()
}

func addMessageChannel(id string, m chan interface{}) {
	mux.Lock()
	message[id] = m
	mux.Unlock()
}

func getMessageChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	m, exist = message[id]
	mux.Unlock()
	return
}

func setMessage(id string, content interface{}) {
	mux.Lock()
	if m, exist := message[id]; exist {
		go func() {
			m <- content
		}()
	}
	mux.Unlock()
}

func setMessageAllClient(content interface{}) {
	mux.Lock()
	all := message
	mux.Unlock()
	go func() {
		for _, m := range all {
			m <- content
		}
	}()

}

func deleteMessageChannel(id string) {
	mux.Lock()
	if m, ok := message[id]; ok {
		close(m)
		delete(message, id)
	}
	mux.Unlock()
}
