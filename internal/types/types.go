package types

import (

	"github.com/omniful/go_commons/db/sql/postgres"
)

type ServerConfig struct {
	Host string
	Port string
}

type AWSConfig struct {
	Region    string
	Account   string
	ShouldLog bool
	SQS       struct {
		Prefix   string
		Endpoint string
	}
}


type AppConfig struct {
	Environment string
	Server      ServerConfig
	Master      postgres.DBConfig
	Slaves      []postgres.DBConfig
	RedisAddr   string
	KafkaBroker string
	AWSConfig   AWSConfig
}
