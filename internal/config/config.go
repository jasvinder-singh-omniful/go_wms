package config

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	"github.com/singhJasvinder101/go_wms/internal/types"
)

var AppConf *types.AppConfig

func InitConfig(ctx context.Context) {
	log.InfofWithContext(ctx, "initializing the local config yaml")

	if err := config.Init(15 * time.Second); err != nil {
		log.ErrorfWithContext(ctx, "error when initializing configs")
		panic(fmt.Errorf("failed to initialize configuration %w", err))
	}

	masterBD := postgres.DBConfig{
			Host:                   config.GetString(ctx, "postgres.master.host"),
			Port:                   config.GetString(ctx, "postgres.master.port"),
			Username:               config.GetString(ctx, "postgres.master.user"),
			Password:               config.GetString(ctx, "postgres.master.password"),
			Dbname:               config.GetString(ctx, "postgres.master.db"),
			MaxOpenConnections:     config.GetInt(ctx, "postgres.master.max_open_connections"),
			MaxIdleConnections:     config.GetInt(ctx, "postgres.master.max_idle_connections"),
			ConnMaxLifetime:        config.GetDuration(ctx, "postgres.master.conn_max_lifetime"),
			DebugMode:              config.GetBool(ctx, "postgres.master.debug_mode"),
			PrepareStmt:            config.GetBool(ctx, "postgres.master.prepare_stmt"),
			SkipDefaultTransaction: config.GetBool(ctx, "postgres.master.skip_default_transaction"),
		}

	slaves := loadSlavesConfig(ctx)

	AppConf = &types.AppConfig{
		Environment: config.GetString(ctx, "env"),
		Server: types.ServerConfig{
			Host: config.GetString(ctx, "http_server.host"),
			Port: config.GetString(ctx, "http_server.port"),
		},
		Master: masterBD,
		Slaves: slaves,
		RedisAddr:   config.GetString(ctx, "redis_addr"),
		KafkaBroker: config.GetString(ctx, "kafka_broker"),
		AWSConfig: types.AWSConfig{
			Region:    config.GetString(ctx, "aws.region"),
			Account:   config.GetString(ctx, "aws.account"),
			ShouldLog: config.GetBool(ctx, "aws.shouldLog"),
			SQS: struct {
				Prefix   string
				Endpoint string
			}{
				Prefix:   config.GetString(ctx, "aws.prefix"),
				Endpoint: config.GetString(ctx, "aws.endpoint"),
			},
		},
	}
}
func loadSlavesConfig(ctx context.Context) []postgres.DBConfig {
    slaves := make([]postgres.DBConfig, 0)
    
    slaveCount := config.GetInt(ctx, "postgres.slaves.count")
    
    for i := 0; i < slaveCount; i++ {
        slavePrefix := fmt.Sprintf("postgres.slaves.slave_%d", i+1)

        slave := postgres.DBConfig{
            Host:                   config.GetString(ctx, slavePrefix+".host"),
            Port:                   config.GetString(ctx, slavePrefix+".port"),
            Username:               config.GetString(ctx, slavePrefix+".user"),
            Password:               config.GetString(ctx, slavePrefix+".password"),
            Dbname:               config.GetString(ctx, slavePrefix+".db"),
            MaxOpenConnections:     config.GetInt(ctx, slavePrefix+".max_open_connections"),
            MaxIdleConnections:     config.GetInt(ctx, slavePrefix+".max_idle_connections"),
            ConnMaxLifetime:        config.GetDuration(ctx, slavePrefix+".conn_max_lifetime"),
            DebugMode:              config.GetBool(ctx, slavePrefix+".debug_mode"),
            PrepareStmt:            config.GetBool(ctx, slavePrefix+".prepare_stmt"),
            SkipDefaultTransaction: config.GetBool(ctx, slavePrefix+".skip_default_transaction"),
        }
        
        if slave.Host != "" {
            slaves = append(slaves, slave)
            log.InfofWithContext(ctx, "slaves dbs configured", 
                "slave_index", i+1, 
                "host", slave.Host, 
                "port", slave.Port)
        }
    }
    
    return slaves
}

func GetConfig() *types.AppConfig {
	return AppConf
}

func InitSQS(ctx context.Context) (*sqs.Queue, *sqs.Publisher) {
	cfg := GetConfig().AWSConfig
	sqsConfig := sqs.GetSQSConfig(ctx, cfg.ShouldLog, cfg.SQS.Prefix, cfg.Region, cfg.Account, cfg.SQS.Endpoint)

	queue, err := sqs.NewStandardQueue(ctx, "bulk-order-queue", sqsConfig)
	if err != nil {
		log.ErrorfWithContext(ctx, "failed to create SQS queue %v", err)
		panic(err)
	}

	publisher := sqs.NewPublisher(queue)
	return queue, publisher
}

func InitKafka(ctx context.Context) *kafka.ProducerClient {
	cfg := GetConfig()
	producer := kafka.NewProducer(
		kafka.WithBrokers([]string{cfg.KafkaBroker}),
		kafka.WithClientID("oms-service"),
	)

	return producer
}
