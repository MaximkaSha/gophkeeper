package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Addr      string
	DSN       string
	CertFile  string
	CertKey   string
	JWTSecret string
}

func NewServerConfig() *ServerConfig {

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	return &ServerConfig{
		Addr:      viper.GetString("addr"),
		DSN:       viper.GetString("dsn"),
		CertFile:  viper.GetString("certfile"),
		CertKey:   viper.GetString("certkey"),
		JWTSecret: viper.GetString("jwtsecret"),
	}
}

type ClientConfig struct {
	Addr     string
	CertFile string
}

func NewClientConfig() *ClientConfig {

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	return &ClientConfig{
		Addr:     viper.GetString("addr"),
		CertFile: viper.GetString("certfile"),
	}
}
