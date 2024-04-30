package telnet

import (
	"fmt"
	"net"
	"os"
	"time"
)

// [root@k8s-master01 ~]# nc -v -w 1 172.30.2.100 9100
// nc: connect to 172.30.2.100 port 9100 (tcp) timed out: Operation now in progress
// [root@k8s-master01 ~]# ./checkPort 172.30.2.100:9100
// [2023-01-17 09:40:37] 172.30.2.100:9100 端口未开启(fail)! dial tcp 172.30.2.100:9100: i/o timeout
//
// [root@k8s-master01 ~]# nc -v -w 1 172.30.2.100 22
// Connection to 172.30.2.100 22 port [tcp/ssh] succeeded!
// [root@k8s-master01 ~]# ./checkPort 172.30.2.100:22
// [2023-01-17 09:41:17] 172.30.2.100:22 端口已开启(success)!
//
// [root@k8s-master01 tmp]# nc -v -w 1 172.30.2.103 9205
// nc: connect to 172.30.2.103 port 9205 (tcp) failed: Connection refused
// [root@k8s-master01 tmp]# ./checkPort 172.30.2.103:9205
// [2023-01-17 09:46:18] 172.30.2.103:9205 端口未开启(fail)! dial tcp 172.30.2.103:9205: connect: connection refused
// 获取IP和端口
//func GetIpPort() []string {
//	// 根据接收参数个数，定义动态数组，
//	ip_ports := make([]string, len(os.Args)-1)
//	i := 0
//	for index, value := range os.Args {
//		//fmt.Println("value", value)
//		//排除脚本名称
//		if index == 0 {
//			continue
//		}
//		//写入数组
//		ip_ports[i] = value
//		i++
//	}
//	return ip_ports
//}

func GetIpPort() string {
	// 根据接收参数个数，定义动态数组，
	return os.Args[1]
}

// 检测端口
func CheckPort(ip_port string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("ip_port", ip_port)
	if len(ip_port) < 1 {
		return fmt.Sprintln("["+now+"]", "请输入ip和端口")
	}
	// 检测端口
	conn, err := net.DialTimeout("tcp", ip_port, 1*time.Second)
	if err != nil {
		return fmt.Sprintln("["+now+"]", ip_port, "端口未开启(fail)!", err)
	} else {
		if conn != nil {
			return fmt.Sprintln("["+now+"]", ip_port, "端口已开启(success)!")
			conn.Close()
		} else {
			return fmt.Sprintln("["+now+"]", ip_port, "端口未开启(fail)!", err)
		}
	}
	return ""
}

//func CheckPorts(ip_ports []string) string {
//	now := time.Now().Format("2006-01-02 15:04:05")
//	if len(ip_ports) < 1 {
//		return fmt.Sprintln("["+now+"]", "请输入ip和端口")
//	}
//	for _, ip_port := range ip_ports {
//		// 检测端口
//		CheckPort(ip_port)
//	}
//	return ""
//}
