#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Реализация Rest API Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Интеграция Apache Kafka в Календарь»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)
- [Домашнее задание №16 «Мониторинг Календаря»](./docs/16_README.md)

#### Запуск

```bash
make up                  # postgres, kafka, calendar, scheduler, storer, prometheus
make down                # остановить и удалить volume'ы
make integration-tests   # поднять окружение, прогнать интеграционные тесты, всё удалить
```

* API: <http://localhost:8888/>
* Swagger/OpenAPI: [`api/swagger.yaml`](./api/swagger.yaml)
* Prometheus: <http://localhost:9090/>
* Метрики: <http://localhost:8888/metrics>, <http://localhost:9102/metrics>, <http://localhost:9103/metrics>

Если стандартные порты заняты, их можно переопределить при запуске:

```bash
POSTGRES_HOST_PORT=55432 KAFKA_HOST_PORT=39092 PROMETHEUS_HOST_PORT=19090 make up
```

Доступные переменные: `POSTGRES_HOST_PORT`, `ZOOKEEPER_HOST_PORT`, `KAFKA_HOST_PORT`,
`CALENDAR_HOST_PORT`, `CALENDAR_METRICS_HOST_PORT`, `SCHEDULER_METRICS_HOST_PORT`,
`STORER_METRICS_HOST_PORT`, `PROMETHEUS_HOST_PORT`, `KAFKA_UI_HOST_PORT`.

Локальный запуск без Docker:

```bash
make run                 # использует configs/config.yaml
make test                # unit-тесты
make lint                # golangci-lint
```

Метрики и примеры запросов PromQL описаны в [`docs/16_METRICS.md`](./docs/16_METRICS.md).

#### Ветки при выполнении
- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)
- `hw16_calendar` (от `hw15_calendar`) -> Merge Request в `hw15_calendar` (если уже вмержена, то в `master`)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**
