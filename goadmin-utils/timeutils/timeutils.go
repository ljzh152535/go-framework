package timeutils

import "time"

// 获取格式化时间
func GetCurrentTime(tmpTime time.Time) string {
	//currentTime := time.Now()
	dateString := tmpTime.Format("2006-01-02 15:04:05.000")
	return dateString
}

// 获取日志
//func GetDate() {
//
//}
