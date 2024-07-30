package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
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

//type GVA_MODEL struct {
//	ID        uint           `json:"id" gorm:"primarykey" form:"id"` // 主键ID
//	CreatedAt time.Time      `json:"createdAt"`                      // 创建时间
//	UpdatedAt time.Time      `json:"updatedAt"`                      // 更新时间
//	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                 // 删除时间
//}

type GVA_MODEL struct {
	ID        uint      `json:"id" gorm:"primarykey" form:"id"` // 主键ID
	CreatedAt LocalTime `json:"createdAt"`                      // 创建时间
	UpdatedAt LocalTime `json:"updatedAt"`                      // 更新时间
	DeletedAt LocalTime `gorm:"index" json:"-"`                 // 删除时间
}

// LocalTime 别名
type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	// LocalTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = LocalTime(t1)
	return err
}
