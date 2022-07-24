package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	PgSQL      *PgSQLConfig
	Server     *Listener
	Services   *Services
	HttpClient *HttpClient
}

type HttpClient struct {
	Timeout time.Duration
}

type Services struct {
	UsersBaseUrl string
}

type PgSQLConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Listener struct {
	Port         string
	ShutdownTime time.Duration
}

func NewConfig() *Config {
	vpr := viper.New()
	vpr.AutomaticEnv()
	vpr.SetDefault(PgDriverName, "postgres")
	vpr.SetDefault(Port, ":8080")
	vpr.SetDefault(ShutdownTime, 10*time.Second)
	vpr.SetDefault(HttpClientTimeout, 10*time.Second)

	return &Config{
		PgSQL: &PgSQLConfig{
			Driver:   vpr.GetString(PgDriverName),
			Host:     vpr.GetString(PgHost),
			Port:     vpr.GetInt(PgPort),
			User:     vpr.GetString(PgUser),
			Password: vpr.GetString(PgPwd),
			DBName:   vpr.GetString(PgDbName),
		},
		Server: &Listener{
			Port:         vpr.GetString(Port),
			ShutdownTime: vpr.GetDuration(ShutdownTime),
		},
		Services: &Services{
			UsersBaseUrl: vpr.GetString(UsersServiceUrl),
		},
		HttpClient: &HttpClient{
			Timeout: vpr.GetDuration(HttpClientTimeout),
		},
	}
}
