package config

import "os"

type (
	Config struct {
		HostDB       string
		UsernameDB   string
		PasswordDB   string
		DbNameDB     string
		PortDB       string
		SearchPathDB string
	}

	AuthConfig struct {
		AccessKey  string
		RefreshKey string
	}
)

var (
	AuthCfg AuthConfig
)

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) SetFromEnv() {
	cfg.HostDB = os.Getenv("HOST_DB")

	cfg.UsernameDB = os.Getenv("USERNAME_DB")

	cfg.PasswordDB = os.Getenv("PASSWORD_DB")

	cfg.DbNameDB = os.Getenv("DBNAME_DB")

	cfg.PortDB = os.Getenv("PORT_DB")

	cfg.SearchPathDB = os.Getenv("SEARCHPATH_DB")
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{}
}

func (authCfg *AuthConfig) SetFromEnv() {
	authCfg.AccessKey = os.Getenv("ACCESS_TOKEN_KEY")
	authCfg.RefreshKey = os.Getenv("REFRESH_TOKEN_KEY")
}
