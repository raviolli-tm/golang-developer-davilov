DROP TABLE IF EXISTS calendar.notification;
CREATE TABLE calendar.notification
(
    notification_id             uuid primary key default gen_random_uuid(),
    notification_event_title    varchar(255) not null,
    notification_event_date     timestamp    not null,
    notification_date           timestamp    not null default now(),
    notification_user_id        bigint       not null
);


COMMENT ON TABLE calendar.notification IS 'Таблица Уведомлений';
COMMENT ON COLUMN calendar.notification.notification_id IS 'Уникальный ID уведомления';
COMMENT ON COLUMN calendar.notification.notification_event_title IS 'Заголовок события';
COMMENT ON COLUMN calendar.notification.notification_event_date IS 'Дата и время события';
COMMENT ON COLUMN calendar.notification.notification_user_id IS 'ID получателя (пользователя)';
COMMENT ON COLUMN calendar.notification.notification_date IS 'Дата и время уведомления';



ALTER TABLE calendar.notification OWNER TO notify_user;

DROP INDEX IF EXISTS calendar.idx_notification;
CREATE INDEX idx_notification ON calendar.notification (notification_id);

GRANT SELECT ON TABLE calendar.notification TO postgres_user
