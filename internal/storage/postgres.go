package storage

import (
	"context"

	"github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/config"
)

type Postgres struct {
	Cluster *postgres.DbCluster
}

func NewPostgres(ctx context.Context) *Postgres {
	cfg := config.GetConfig()

	dbCluster := postgres.InitializeDBInstance(cfg.Master, &cfg.Slaves)

	log.Info("postgreSQL database initialized successfully",
		"host", cfg.Master.Host,
		"port", cfg.Master.Port,
		"database", cfg.Master.Dbname,
	)

	return &Postgres{
		Cluster: dbCluster,
	}
}
