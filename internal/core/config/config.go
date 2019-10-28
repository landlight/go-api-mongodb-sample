package config

import (
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	CF = &Configs{}
)

type Configs struct {
	Release    bool `mapstructrue:"release"`
	Port       int  `mapstructrue:"port"`
	HTTPServer struct {
		ReadTimeout       time.Duration `mapstructrue:"read_timeout"`
		WriteTimeout      time.Duration `mapstructrue:"write_timeout"`
		ReadHeaderTimeout time.Duration `mapstructrue:"read_header_timeout"`
	} `mapstructure: "http_server"`
	MongoDB struct {
		AddrWithPort string        `mapstructure:"hosts"`
		Timeout      time.Duration `mapstructure:"timeout"`
		Username     string        `mapstructure:"username"`
		Password     string        `mapstructure:"password"`
		Schema       struct {
			DBName string `mapstructure:"db_name"`
		} `mapstructure:"schema"`
		DialInfo struct {
			DBName *mgo.DialInfo
		} `mapstructure:"-"`
	} `mapstructure:"mongo_db"`
}

func InitConfig(configPath string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("config file changed:", e.Name)
		if err := bindingConfig(v, CF); err != nil {
			logrus.Error("binding error:", err)
		}
		logrus.Infof("Initial 'Configuration'. %+v", CF)
	})
	return nil
}

func bindingConfig(vp *viper.Viper, cfg *Configs) error {
	if err := vp.Unmarshal(&cfg); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}

	md := mgo.DialInfo{
		Addrs:    strings.Split(cfg.MongoDB.AddrWithPort, ","),
		Timeout:  cfg.MongoDB.Timeout,
		Database: cfg.MongoDB.Schema.DBName,
		Username: cfg.MongoDB.Username,
		Password: cfg.MongoDB.Password,
	}
	cfg.MongoDB.DialInfo.DBName = &md

	return nil
}
