package goadmin_config

import "fmt"

type System struct {
	Host            string `yaml:"host"`                                          // ipInfo
	Port            int    `yaml:"port"`                                          // 端口值
	Env             string `yaml:"env"`                                           // 环境值
	DbType          string `mapstructure:"db-type" json:"db-type" yaml:"db-type"` // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	JWT_SIGN_KEY    string `yaml:"JWT_SIGN_KEY"`
	JWT_EXPIRE_TIME int64  `yaml:"JWT_EXPIRE_TIME"`
}

func (s System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
