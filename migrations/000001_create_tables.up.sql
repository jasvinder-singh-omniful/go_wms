create table if not exists hubs (
    id serial primary key,

    tenant_id text not null,
    name text not null,
    location jsonb,

    created_at timestamp with time zone default now()
);

create index if not exists idx_hubs_tenant on hubs(tenant_id);

create table if not exists skus(
    id serial primary key,

    tenant_id text not null,
    seller_id text not null,
    sku_code text not null,
    name text,
    metadata jsonb,
    
    created_at timestamp with time zone default now(),
    unique(tenant_id, seller_id, sku_code)
);

create index if not exists idx_skus_tenant_seller on skus(tenant_id, seller_id);


create table if not exists inventory (
    id bigserial primary key,

    tenant_id text not null,
    seller_id text not null,
    hub_id int not null references hubs(id) on delete cascade,
    sku_id int not null references skus(id) on delete cascade,
    quantity bigint not null default 0,

    updated_at timestamp with time zone default now(),
    unique(sku_id, hub_id)
);

create index if not exists idx_inventory_sku_hub on inventory(sku_id, hub_id);


