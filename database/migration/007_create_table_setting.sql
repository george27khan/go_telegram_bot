create table go_bot."setting"
(
    id serial constraint idx_pk_setting primary key,
    setting_code varchar(100) not null,
    setting_describe varchar(100) not null,
    number_value numeric,
    string_value varchar(1000),
    json_value json,
    date_value timestamp,
    created_dt timestamp default now() not null
    );
alter table go_bot."setting" owner to postgres;

---- create above / drop below ----
drop table go_bot."setting"