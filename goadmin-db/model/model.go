package model

import (
	"gorm.io/gorm"
	"time"
)

type DBItemConf struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"database"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Config       string `yaml:"config"`                                       // 高级配置
	Timeout      int    `yaml:"timeout"`                                      // connect db timeout , uint ms
	WriteTimeOut int    `yaml:"write_time_out" mapstructure:"write_time_out"` // write data timeout , uint ms
	ReadTimeOut  int    `yaml:"read_time_out" mapstructure:"read_time_out"`   // read data timeout,uint ms
	MaxIdleConns int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"` // 最大的闲置连接数
	MaxOpenConns int    `yaml:"max_open_conns" mapstructure:"max_open_conns"` //最大打开连接数
}

type DBLog struct {
	Enable bool   `yaml:"enable"`
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
}

type DBItem struct {
	Write DBItemConf `yaml:"write"`
	Read  DBItemConf `yaml:"read"`
	Log   DBLog      `yaml:"log"`
}

type GVA_MODEL struct {
	ID        uint           `json:"id" gorm:"primarykey" form:"id"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
