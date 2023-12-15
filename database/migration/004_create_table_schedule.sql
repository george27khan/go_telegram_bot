-- Write your migrate up statements here
create table go_bot."schedule"
(
    id bigserial primary key,
    id_user bigint not null constraint fk_sched_user references go_bot."user",
    visit_dt timestamp not null,
    created_dt timestamp default now() not null
);

alter table go_bot."schedule" owner to postgres;

create index idx_fk_sched_user on go_bot."schedule" (id_user);

create unique index schedule_idx1
    on go_bot."schedule" (id_user, visit_dt);

comment on index go_bot.schedule_idx1 is 'Уникальная запись на время по клиенту';
---- create above / drop below ----
drop table go_bot."schedule"
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
