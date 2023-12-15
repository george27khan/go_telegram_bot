create table go_bot."position"
(
    id serial constraint idx_pk_position primary key,
    position_name varchar(100) not null,
    created_dt timestamp default now() not null
);

alter table go_bot."position" owner to postgres;
create unique index position_idx1
    on go_bot."position" (position_name);

---- create above / drop below ----
drop table go_bot."position"

