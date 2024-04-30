package websocketLog

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func LogWs(c *gin.Context) {
	// 开启socket管理器
	//go manager.start()
	//// 监控文件
	//go monitor("filePath")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	// 开启httpServer
	//go server(port)

	//var err error
	//msg := c.DefaultQuery("msg", "")
	//cols := c.DefaultQuery("cols", "150")
	//rows := c.DefaultQuery("rows", "35")
	////sshType := c.DefaultQuery("sshType", "key")
	//col, _ := strconv.Atoi(cols)
	//row, _ := strconv.Atoi(rows)
	//terminal := connections.Terminal{
	//	Columns: uint32(col),
	//	Rows:    uint32(row),
	//}
	//sshClient, err := connections.DecodedMsgToSSHClient(msg)
	//if err != nil {
	//	c.Error(err)
	//	return
	//}
	//if sshClient.IpAddress == "" || sshClient.Username == "" {
	//	c.Error(&ApiError{Message: "IP地址或用户不能为空", Code: 400})
	//	return
	//}
	//if sshClient.SshType == "" {
	//	c.Error(&ApiError{Message: "sshType 不能为空", Code: 400})
	//	return
	//}
	//conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	//if err != nil {
	//	c.Error(err)
	//	return
	//}
	//err = sshClient.GenerateClient()
	//if err != nil {
	//	conn.WriteMessage(1, []byte(err.Error()))
	//	conn.Close()
	//	return
	//}
	////sshClient.RequestTerminal(terminal)
	//
	//sshClient.Connect(conn)
	//fmt.Println(terminal)
}

// 创建client对象
func genConn(ws *websocket.Conn) {
	//client := &client{time.Now().String(), ws, make(chan []byte, 1024)}
	//manager.register <- client
	//go client.read()
	//client.write()
}
