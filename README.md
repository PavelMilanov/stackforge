# stackforge

example.config.yml

```yaml
portainer:
  realm: "http://localhost"
  token: "123"
  teams: [1]

git:
  realm: "http://localhost"
  token: "123"

cors:
  origin: "*"
```


example.metadata.yaml

```
title: Go API Starter
category: Backend
purpose: Быстрый запуск dev-стенда для backend-разработки.
fit: Подходит для тестирования API, миграций и интеграций.

parameters:
  - имя стенда
  - namespace/project
  - branch/tag
  - endpoint

services:
  - name: go-api
    note: HTTP service, port 1323
  - name: postgres
    note: Database, internal
  - name: redis
    note: Cache, internal

```
