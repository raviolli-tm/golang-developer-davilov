

CREATE SCHEMA if not exists calendar;

DROP TABLE IF EXISTS calendar.calendar_event;
CREATE TABLE calendar.calendar_event (
                                event_id uuid primary key default gen_random_uuid(),
                                title varchar(255) not null,-- Заголовок
                                date_start timestamp not null, --Дата и время события
                                date_end timestamp not null,  --Длительность события (или дата и время окончания)
                                event_description text, --Описание события - длинный текст, опционально
                                user_id bigint not null, --ID пользователя, владельца события
                                event_notify_time bigint --За сколько времени высылать уведомление, опционально.
);

DROP INDEX IF EXISTS calendar.idx_calendar_event;
CREATE INDEX idx_calendar_event ON calendar.calendar_event (event_id);

