package curl

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func CurlGet(url string) string {
	client := http.Client{
		Timeout: time.Second * 1, // 设置超时时间为5秒
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Sprintf("curl req fail with error:%s", err.Error())
	}
	req.Header.Set("acl", "56FFCBF8CDBC174DC6E05DE3DC493D6801D49137A029CBC6960B028F202D6805") // 设置请求头
	req.SetBasicAuth("lemonActuator", "lemonActuator123#")                                    // 设置basicAuth访问
	res, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("curl res fail with error:%s", err.Error())
	}

	//command, err := http2curl.GetCurlCommand(req)
	//if err != nil {
	//	return fmt.Sprintf("curl fail with error:%s", err.Error())
	//}
	//fmt.Println(command)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("curl fail with error:%s", err.Error())
	}

	return string(body)
}
