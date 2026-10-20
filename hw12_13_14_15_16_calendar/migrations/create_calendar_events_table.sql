CREATE SCHEMA IF NOT EXISTS calendar;

DROP TABLE IF EXISTS calendar.calendar_event;
CREATE TABLE calendar.calendar_event
(
    event_id          uuid primary key default gen_random_uuid(),
    title             varchar(255) not null,
    date_start        timestamp    not null,
    date_end          timestamp    not null,
    event_description text,
    user_id           bigint       not null,
    event_notify_time bigint       not null default 0
);

DROP INDEX IF EXISTS calendar.idx_calendar_event;
CREATE INDEX idx_calendar_event ON calendar.calendar_event (event_id);
CREATE INDEX idx_calendar_event_date_start ON calendar.calendar_event (date_start);
