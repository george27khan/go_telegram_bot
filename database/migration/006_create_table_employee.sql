create table go_bot."employee"
(
    id serial constraint idx_pk_employee primary key,
    first_name varchar(100) not null,
    middle_name varchar(100) not null,
    last_name varchar(100) not null,
    birth_date date not null,
    email varchar(100) not null,
    phone_number varchar(100) not null,
    id_position integer not null constraint fk_emp_pos references go_bot."position",
    hire_date date not null,
    created_dt timestamp default now() not null,
    photo bytea not null
    );

alter table go_bot."employee" owner to postgres;
create index idx_fk_emp_pos on go_bot."employee" (id_position);

---- create above / drop below ----
drop table go_bot."employee"
