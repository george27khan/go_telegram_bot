-- Write your migrate up statements here
create table go_bot."user"
(
    id         serial primary key,
    telegram   varchar(100) not null,
    phone      varchar(20),
    user_name  varchar(20),
    created_dt timestamp default now() not null
);

alter table go_bot."user"
    owner to go_bot;
---- create above / drop below ----
drop table go_bot."user";

