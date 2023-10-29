package global

import (
	"databaseDemo/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/go-redis/redis/v9"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	Log         *zap.Logger
	DB          *gorm.DB
	RDB         *sql.DB
	Redis 		*redis.Client
}

var App = new(Application)
