package main

import (
	"fmt"
	"github.com/abrander/go-supervisord"
)

func main() {
	c, err := supervisord.NewClient("http://10.176.38.152:11911/RPC2", supervisord.ClientOption(supervisord.WithAuthentication("supervisor", "Supervisor@01")))
	defer c.Close()
	if err != nil {
		panic(err.Error())
	}

	info, err := c.GetAllProcessInfo()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(info)
	//err = c.ClearLog()
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//
	//
	//err = c.Restart()
	//if err != nil {
	//	panic(err.Error())
	//}
}
