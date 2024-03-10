package config

import "os"

type Config struct {
	HostDB       string
	UsernameDB   string
	PasswordDB   string
	DbNameDB     string
	PortDB       string
	SearchPathDB string
}

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
