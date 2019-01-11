package {{Name}}

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBName        string
	DBUser        string
	Addr          string
	IsDevelopment bool
	SessionKey    string
	SessionName   string
}

func NewConfigFromEnv() *Config {
	viper.SetDefault("dbname", "{{Name}}")
	viper.SetDefault("addr", ":4040")
	viper.SetDefault("dbuser", "bnadland")
	viper.SetDefault("isdevelopment", true)
	viper.SetDefault("sessionname", "{{Name}}")
	viper.SetDefault("sessionkey", "super-secret-key")
	viper.AutomaticEnv()
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Print(err)
		return nil
	}
	return &config
}
