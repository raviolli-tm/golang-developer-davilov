## Метрики сервиса «Календарь»

Сервис экспортирует метрики в формате Prometheus:

* API — `http://localhost:8888/metrics`;
* Планировщик — `http://localhost:9102/metrics`;
* Хранитель — `http://localhost:9103/metrics`.

При запуске через `make up` поднимается Prometheus (`http://localhost:9090`),
который собирает все три таргета (см. `deployments/prometheus.yml`).

### Метрики API

| Метрика | Тип | Назначение |
|---------|-----|------------|
| `calendar_http_requests_total{method,path,status}` | counter | Количество запросов с разбивкой по методу, маршруту и HTTP-статусу. Позволяет видеть долю ошибок и рост нагрузки. |
| `calendar_http_request_duration_seconds{method,path}` | histogram | Время обработки запросов. По гистограмме считаются p50/p95/p99 и выявляются медленные ручки. |
| `calendar_http_requests_in_flight` | gauge | Количество одновременно обрабатываемых запросов. Помогает заметить перегрузку. |
| `calendar_events_created_total` | counter | Успешно созданные события. |
| `calendar_events_updated_total` | counter | Успешно обновлённые события. |
| `calendar_events_deleted_total` | counter | Успешно удалённые события. |
| `calendar_events_errors_total{operation}` | counter | Ошибки бизнес-операций (`create`, `update`, `delete`, `list`, `list_period`). |

### Метрики планировщика

| Метрика | Тип | Назначение |
|---------|-----|------------|
| `calendar_scheduler_scans_total` | counter | Количество циклов сканирования БД. Если метрика не растёт — планировщик завис. |
| `calendar_scheduler_notifications_sent_total` | counter | Отправленные в Kafka уведомления. |
| `calendar_scheduler_events_deleted_total` | counter | Удалённые устаревшие события (старше года). |
| `calendar_scheduler_errors_total` | counter | Ошибки чтения БД и отправки сообщений. |

### Метрики хранителя

| Метрика | Тип | Назначение |
|---------|-----|------------|
| `calendar_storer_notifications_stored_total` | counter | Уведомления, успешно сохранённые в БД. Сравнение с `calendar_scheduler_notifications_sent_total` показывает потери в брокере. |
| `calendar_storer_notifications_errors_total` | counter | Ошибки разбора сообщений и записи в БД. |

### Примеры запросов

```promql
# RPS по ручкам
sum by (path) (rate(calendar_http_requests_total[1m]))

# Доля 5xx
sum(rate(calendar_http_requests_total{status=~"5.."}[5m]))
  / sum(rate(calendar_http_requests_total[5m]))

# p95 времени ответа
histogram_quantile(0.95, sum by (le) (rate(calendar_http_request_duration_seconds_bucket[5m])))

# Отставание хранителя от планировщика
sum(calendar_scheduler_notifications_sent_total) - sum(calendar_storer_notifications_stored_total)
```

Стандартные `go_*` и `process_*` метрики также доступны на всех эндпоинтах.
