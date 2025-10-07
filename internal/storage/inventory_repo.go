package storage

import (
	"context"

	"github.com/singhJasvinder101/go_wms/models"
)

type InventoryRepo struct {
	DB *Postgres
}

func (r *InventoryRepo) Create(ctx context.Context, inventory *models.Inventory)

