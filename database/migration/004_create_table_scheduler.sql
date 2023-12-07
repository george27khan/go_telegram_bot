-- Write your migrate up statements here
create table go_bot."scheduler"
(
    id serial constraint idx_pk_scheduler primary key,
    id_user integer not null constraint fk_sched_user references go_bot."user",
    visit_dt date not null,
    created_dt timestamp default now() not null
);

alter table go_bot."scheduler" owner to postgres;

create index idx_fk_sched_user on go_bot."scheduler" (id_user);

---- create above / drop below ----
drop table go_bot."scheduler"
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
