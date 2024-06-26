package goadmin_os

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
)

func LocalIP() string {

	ipList := []string{"114.114.114.114:80", "8.8.8.8:80"}
	for _, ip := range ipList {
		conn, err := net.Dial("udp", ip)
		if err != nil {
			continue
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
		return localAddr.IP.String()
	}
	// 从网卡中获取
	return ""

}

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func GetHomePath() string {
	var homeDir string
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("尝试获取文件路径失败!")
	}
	homeDir = filepath.Dir(f)
	return homeDir
}
