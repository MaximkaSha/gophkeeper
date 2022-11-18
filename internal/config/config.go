// Package config manage server and client configurations.
package config

import "github.com/spf13/viper"

// ServerConfig - Server configuration structure.
type ServerConfig struct {
	// Host:port of server.
	Addr string
	// DSN of database.
	DSN string
	// Path to certificate file.
	CertFile string
	// Path to certificate key file.
	CertKey string
	// JWT key string.
	JWTSecret string
}

// NewServerConfig - ServerConfig constructor.
// Gets data from config file.
func NewServerConfig() *ServerConfig {

	viper.AddConfigPath("./")
	viper.SetConfigName("config_server")
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

// ClientConfig - CLient config file
type ClientConfig struct {
	// Host:port of server
	Addr string
	// Path to certificate file.
	CertFile string
}

// NewClientConfig ClientConfig constructor.
// Gets data from config file.
func NewClientConfig() *ClientConfig {

	viper.AddConfigPath("./")
	viper.SetConfigName("config_client")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	return &ClientConfig{
		Addr:     viper.GetString("addr"),
		CertFile: viper.GetString("certfile"),
	}
}
