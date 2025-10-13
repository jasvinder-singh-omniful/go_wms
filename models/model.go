package models

import (
	"time"
	"gorm.io/datatypes"
)

type Hub struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`

	TenantID  string         `gorm:"type:text;not null;index:idx_hubs_tenant" json:"tenant_id"`
	Name      string         `gorm:"type:text;not null" json:"name"`
	Location  datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"location"`
	
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

type SKU struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`

	TenantID  string         `gorm:"type:text;not null;" json:"tenant_id"`
	SellerID  string         `gorm:"type:text;not null;" json:"seller_id"`
	SKUCode   string         `gorm:"type:text;not null;" json:"sku_code"`
	Name      string         `gorm:"type:text" json:"name"`
	MetaData  datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'" json:"metadata"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

type Inventory struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`

	TenantID  string    `gorm:"not null;" json:"tenant_id"`
	SellerID  string    `gorm:"not null;" json:"seller_id"`
	HubID     int       `gorm:"not null;" json:"hub_id"`
	SKUID     int       `gorm:"column:sku_id;not null;" json:"sku_id"`
	Quantity  int64     `gorm:"not null;default:0;check:quantity>=0" json:"quantity"`
	
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}


func (Inventory) TableName() string {
    return "inventory"
}