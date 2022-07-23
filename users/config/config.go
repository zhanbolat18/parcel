package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Jwt            *JwtConfig
	PasswordHasher *PasswordHasherConfig
	PgSQL          *PgSQLConfig
	Server         *Listener
}

type JwtConfig struct {
	Ttl           time.Duration
	BaseTimeDelta time.Duration
	SignKey       []byte
}

type PasswordHasherConfig struct {
	Cost int
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
	vpr.SetDefault(JwtTokenTtl, 1*time.Hour)
	vpr.SetDefault(JwtBaseTimeDelta, 0)
	vpr.SetDefault(PasswordHashCost, 13)
	vpr.SetDefault(PgDriverName, "postgres")
	vpr.SetDefault(Port, ":8080")
	vpr.SetDefault(ShutdownTime, 10*time.Second)

	return &Config{
		Jwt: &JwtConfig{
			Ttl:           vpr.GetDuration(JwtTokenTtl),
			BaseTimeDelta: vpr.GetDuration(JwtBaseTimeDelta),
			SignKey:       []byte(vpr.GetString(JwtSignKey)),
		},
		PasswordHasher: &PasswordHasherConfig{
			Cost: vpr.GetInt(PasswordHashCost),
		},
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
	}
}
