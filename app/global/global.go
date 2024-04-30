package global

import (
	"github.com/ljzh152535/go-framework/app/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_CONFIG config.Config
	GVA_VP     *viper.Viper
	GVA_LOG    *logrus.Entry
	GVA_DB     *gorm.DB
)
