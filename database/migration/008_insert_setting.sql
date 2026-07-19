-- Write your migrate up statements here
insert into go_bot.setting(setting_code, setting_describe, number_value)
values ('session_time_minute', 'Продолжительность приема в минутах', 15),
       ('time_keyboar_width', 'Ширина клавиатуры выбора времени', 3),
       ('days_in_schedule', 'Число дней доступных для бронирования', 40);

insert into go_bot.setting(setting_code, setting_describe, json_value)
values ('start_hour_schedule', 'График начала рабочего дня', '{"Mon": "18:00", "Tue": "18:00", "Wed": "18:00", "Thu": "18:00", "Fri": "18:00", "Sat": "18:00", "Sun": "18:00"}'),
       ('end_hour_schedule', 'График конца рабочего дня', '{"Mon": "9:00", "Tue": "9:00", "Wed": "9:00", "Thu": "9:00", "Fri": "9:00", "Sat": "9:00", "Sun": "9:00"}');

---- create above / drop below ----
delete from go_bot.setting;
