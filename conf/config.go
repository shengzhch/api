package conf

import (
	"flag"
	"github.com/spf13/viper"
)

var (
	ConfPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&ConfPath, "conf", "./conf/config.yaml", "config path")
}

func Init() error {
	viper.SetConfigFile(ConfPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Http HttpConfig `json:"http"`
}

type HttpConfig struct {
	Port int    `json:"port"`
	Mode string `json:"mode"`
}
