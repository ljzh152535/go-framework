package goadmin_viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper //
// 参数 a interface 绑定的结构体
func InitViper(a interface{}) *viper.Viper {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&a); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&a); err != nil {
		fmt.Println(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	//global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
