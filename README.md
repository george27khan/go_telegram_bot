## Описание
Данный telegram bot создан для бронирования времени в календаре.
- Для работы с telegram api использовалась бибилиотека [go-telegram ](https://github.com/go-telegram)
- Для хранения данных о бронировании использовалась СУБД Postgres.
- Для хранения состояний использовалась СУБД Redis (в качестве практики, а не как лучшее решение проблемы).
- Для миграции базы использовалась библиотека [tern](https://github.com/jackc/tern)

Бот со всеми компонентами разворачивается через docker-compose.  
Для корректного подключение бота к базе Postgres необходимо создать файл [db.env](https://github.com/george27khan/go_telegram_bot/blob/main/db.env) и прописать там настройки для подключения

    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=postgres

Для корректной работы миграций нужно по необходимости внести правки в файл [tern.conf](https://github.com/george27khan/go_telegram_bot/blob/main/tern.conf)  
Чтобы запускать миграцию из консоли необходимо прописать перменные окружения

    TERN_CONFIG={путь до проекта}
    TERN_MIGRATIONS={путь до проекта}\database\migration

