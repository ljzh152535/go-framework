package connections

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

func DecodedMsgToSSHClient(msg string) (SSHClient, error) {
	client := NewSSHClient()
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(decoded, &client)
	if err != nil {
		return client, err
	}
	return client, nil
}

func formatError(err error, errStr string) string {
	if err != nil {
		return fmt.Sprintf("%s:%s", errStr, err.Error())
	}
	return ""
}

func publicKeyAuthFunc() (sshAuthMethod ssh.AuthMethod, errStr string) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, formatError(err, "get UserHomeDir failed")
	}

	key, err := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
	if err != nil {
		return nil, formatError(err, "ssh key file read failed")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, formatError(err, "ssh key signer failed")
	}
	return ssh.PublicKeys(signer), ""
}

func (this *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)

	if this.SshType == "key" {
		publicKeyAuthFuncRes, str := publicKeyAuthFunc()
		if str == "" {
			auth = []ssh.AuthMethod{publicKeyAuthFuncRes}
		} else {
			return errors.New("publicKeyAuthFunc get fail")
		}
	} else {
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(this.Password))
	}

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    this.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", this.IpAddress, this.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	this.Client = client
	return nil
}

func (this *SSHClient) RequestTerminal(terminal Terminal) *SSHClient {
	session, err := this.Client.NewSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	this.Session = session
	channel, inRequests, err := this.Client.OpenChannel("session", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	this.channel = channel
	go func() {
		for req := range inRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    uint32(terminal.Columns * 8),
		Height:   uint32(terminal.Columns * 8),
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	return this
}

func (this *SSHClient) Connect(ws *websocket.Conn) {
	//这里第一个协程获取用户的输入
	//go func() {
	//	for {
	//		// p为用户输入
	//		_, p, err := ws.ReadMessage()
	//		if err != nil {
	//			return
	//		}
	//		_, err = this.channel.Write(p)
	//		if err != nil {
	//			return
	//		}
	//	}
	//}()
	execute(ws, this.Session)
	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		br := bufio.NewReader(this.channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			defer this.Client.Close()
			defer this.Session.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					log.Println(err)
					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
					ws.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		// 主循环
		for {
			select {
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}

func asyncLog(reader io.ReadCloser, ws *websocket.Conn) error {
	bucket := make([]byte, 1024)
	buffer := make([]byte, 100)
	for {
		num, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return err
		}
		if num > 0 {
			line := ""
			bucket = append(bucket, buffer[:num]...)
			tmp := string(bucket)
			if strings.Contains(tmp, "\n") {
				ts := strings.Split(tmp, "\n")
				if len(ts) > 1 {
					line = strings.Join(ts[:len(ts)-1], "\n")
					bucket = []byte(ts[len(ts)-1]) //不够整行的以后再处理
				} else {
					line = ts[0]
					bucket = bucket[:0]
				}
				err := ws.WriteMessage(websocket.TextMessage, []byte(line))
				if err != nil {
					log.Println(err)
					return err
				}
				fmt.Printf("%s\n", line)
			}

		}
	}
	return nil
}

func execute(ws *websocket.Conn, session *ssh.Session) error {

	//cmd := output
	//stdout, _ := session.StdoutPipe()
	//stderr, _ := session.StderrPipe()
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	err := session.Run("sh ./scripts/curl.sh")
	if err != nil {
		return err
	}
	//if err := cmd.Start(); err != nil {
	//	log.Printf("Error starting command: %s......", err.Error())
	//	return err
	//}

	ws.WriteMessage(websocket.TextMessage, []byte("dsfsdfds"))
	//go asyncLog(session.stdout, ws)
	//go asyncLog(stderr, ws)

	//if err := cmd.Wait(); err != nil {
	//	log.Printf("Error waiting for command execution: %s......", err.Error())
	//	return err
	//}

	return nil
}
