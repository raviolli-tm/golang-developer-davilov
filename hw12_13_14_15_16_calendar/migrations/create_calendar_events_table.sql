DO $$
BEGIN

IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'postgres_user') THEN
    CREATE USER postgres_user WITH PASSWORD 'postgres_password';
    RAISE NOTICE 'Role "postgres_user" created successfully';
ELSE
    RAISE NOTICE 'Role "postgres_user" already exists';
END IF;

END $$;

CREATE SCHEMA if not exists calendar AUTHORIZATION postgres_user;

GRANT USAGE ON SCHEMA calendar TO postgres_user;
GRANT CREATE ON SCHEMA calendar TO postgres_user;


DROP TABLE IF EXISTS calendar.calendar_event;
CREATE TABLE calendar.calendar_event
(
    event_id          uuid primary key default gen_random_uuid(),
    title             varchar(255) not null,
    date_start        timestamp    not null,
    date_end          timestamp    not null,
    event_description text,
    user_id           bigint       not null,
    event_notify_time bigint
);

COMMENT ON TABLE calendar.calendar_event IS 'Таблица событий календаря';
COMMENT ON COLUMN calendar.calendar_event.event_id IS 'Уникальный ID события';
COMMENT ON COLUMN calendar.calendar_event.title IS 'Заголовок';
COMMENT ON COLUMN calendar.calendar_event.date_start IS 'Дата и время начала события';
COMMENT ON COLUMN calendar.calendar_event.date_end IS 'Дата и время окончания события';
COMMENT ON COLUMN calendar.calendar_event.event_description IS 'Описание события, опционально';
COMMENT ON COLUMN calendar.calendar_event.user_id IS 'Уникальный ID пользователя';
COMMENT ON COLUMN calendar.calendar_event.event_notify_time IS 'За сколько времени высылать уведомление, опционально';


ALTER TABLE calendar.calendar_event OWNER TO postgres_user;

DROP INDEX IF EXISTS calendar.idx_calendar_event;
CREATE INDEX idx_calendar_event ON calendar.calendar_event (event_id);

