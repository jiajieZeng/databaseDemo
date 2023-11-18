package global

import (
	"database/sql"
	"databaseDemo/config"

	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	Log         *zap.Logger
	DB          *gorm.DB
	RDB         *sql.DB
	Redis       *redis.Client
	Tx          *sql.Tx
}

var App = new(Application)
