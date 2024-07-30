package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
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

//type GVA_MODEL struct {
//	ID        uint           `json:"id" gorm:"primarykey" form:"id"` // 主键ID
//	CreatedAt time.Time      `json:"createdAt"`                      // 创建时间
//	UpdatedAt time.Time      `json:"updatedAt"`                      // 更新时间
//	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                 // 删除时间
//}

type GVA_MODEL struct {
	ID        uint           `json:"id" gorm:"primarykey" form:"id"` // 主键ID
	CreatedAt CustomTime     `json:"createdAt"`                      // 创建时间
	UpdatedAt CustomTime     `json:"updatedAt"`                      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                 // 删除时间
}

type CustomTime time.Time

// GORM Scanner 接口, 从数据库读取到类型
func (t *CustomTime) Scan(value any) error {

	if v, ok := value.(time.Time); !ok {
		return errors.Errorf("failed to unmarshal CustomTime value: %v", value)
	} else {
		*t = CustomTime(v)
		return nil
	}
}

// GORM Valuer 接口, 保存到数据库
func (t CustomTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

// JSON Marshal接口，CustomTime结构体转换为json字符串
func (t *CustomTime) MarshalJSON() ([]byte, error) {
	t2 := time.Time(*t)
	return []byte(fmt.Sprintf(`"%v"`, t2.Format("2006-01-02 15:04:05"))), nil
}

// fmt.Printf, 【可选方法】
func (t CustomTime) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
