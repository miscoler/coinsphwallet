package postgres

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ListenPort int             `json:"port"`
	Postgres   *PostgresConfig `json:"postgres"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func Parse(configFile string) (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
