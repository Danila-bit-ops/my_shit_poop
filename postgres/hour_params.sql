create table hour_params
(
    id         bigserial      not null,
    val        numeric(10, 3) not null default 0,
    timestamp  timestamp      not null default now(),
    param_id   integer        not null default 0,
    xml_create boolean        not null default false,
    manual     boolean        not null default false,
    change_by  varchar        not null default '',
    comment    text           not null default ''
);

create index hour_params__param_id on hour_params (param_id);
