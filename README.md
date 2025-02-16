# avito-shop

## Запуск
Для запуска используйте Docker. \

Запустить сервис

```bash
docker build .
docker compose --profile run up -d
```

Запустить тесты

```bash
docker build .
docker compose --profile test up --build --abort-on-container-exit -d
```
