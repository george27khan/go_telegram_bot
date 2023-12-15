-- Write your migrate up statements here
insert into go_bot.setting(setting_code, setting_describe, number_value)
values ('session_time_hour', 'Продолжительность приема в часах', 0.25),
       ('time_keyboar_width', 'Ширина клавиатуры выбора времени', 3),
       ('days_in_schedule', 'Число дней доступных для бронирования', 40);

insert into go_bot.setting(setting_code, setting_describe, json_value)
values ('start_hour_schedule', 'График начала рабочего дня', '{"Mon": 9, "Tue": 9, "Wed": 9, "Thu": 9, "Fri": 9, "Sat": 9, "Sun": 9}'),
       ('end_hour_schedule', 'График конца рабочего дня', '{"Mon": 18, "Tue": 18, "Wed": 18, "Thu": 18, "Fri": 18, "Sat": 16, "Sun": 16}');

---- create above / drop below ----
delete from go_bot.setting;
