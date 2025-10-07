package storage

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/models"
	"gorm.io/gorm"
)


type HubRepo struct {
	DB *Postgres
}

func (r *HubRepo) Create(ctx context.Context, hub *models.Hub) error {
	logTag := "[HubRepo][Create]"
	log.InfofWithContext(ctx, logTag+" creating hub iin database ", "hub", hub)

	db := r.DB.Cluster.GetMasterDB(ctx)

	if err := db.Create(&hub).Error; err != nil {
		log.ErrorfWithContext(ctx, logTag+" error when creating hub in db", err)
		return nil
	}

	log.InfofWithContext(ctx, "hub created successfully")
	return nil
}

func (r *HubRepo) GetByID(ctx context.Context, id uint) (*models.Hub, error) {
	logTag := "[HubRepo][GetByID]"
	log.InfofWithContext(ctx, logTag+" geting hub by id in database ", "id", id)

	db := r.DB.Cluster.GetSlaveDB(ctx)

	var hub models.Hub
	if err := db.First(&hub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no hub found with id %d", id)
		}
		log.ErrorfWithContext(ctx, logTag+" error when finding hub in db", err)
		return nil, fmt.Errorf("error when fetching hub by id %v", err)
	}

	log.InfofWithContext(ctx, "hub fetched successfully")
	return &hub, nil
}

func (r *HubRepo) GetAll(ctx context.Context) ([]models.Hub, error) {
	logTag := "[HubRepo][GetAll]"
	log.InfofWithContext(ctx, logTag+" geting all hubs in database ")


	var hubs []models.Hub
	db := r.DB.Cluster.GetSlaveDB(ctx)

	if err := db.Find(&hubs).Error; err != nil {
		log.ErrorfWithContext(ctx, logTag+" error when finding hubs in db", err)
		return nil, fmt.Errorf("error when fetching hubs%v", err)
	}

	return hubs, nil
}


