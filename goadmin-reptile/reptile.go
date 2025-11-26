package goadmin_reptile

import (
	"fmt"
	"io"
	"net/http"
)

func Fetch(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Add("Cookie", "JSESSIONID=_DvL7c5D0RtcON9mGKvPKzhf9yo_XYUpddb2wtIx; Hm_lvt_0882a29bba355751ac6f4fe522f8d6de=1646879724,1646900491,1646973412,1646979544; Hm_lpvt_0882a29bba355751ac6f4fe522f8d6de=1646981342")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}

	return string(body)

}
