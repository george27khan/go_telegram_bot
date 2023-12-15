-- Write your migrate up statements here
create table go_bot."user"
(
    id         bigint primary key,--serial primary key,
    user_name  varchar(100) not null,
    first_name varchar(100),
    last_name  varchar(100),
    phone      varchar(20),
    created_dt timestamp default now() not null
);

alter table go_bot."user"
    owner to postgres;
---- create above / drop below ----
drop table go_bot."user";

