drop index if exists idx_inventory_sku_hub;
drop index if exists idx_skus_tenant_seller;
drop index if exists idx_hubs_tenant;

drop table if exists inventory;
drop table if exists skus;
drop table if exists hubs;
