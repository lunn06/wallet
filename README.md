# Wallet Backend
Выполненное тестовое задание для прохождения стажировки в компании ИнфоТекС

## Запуск
Перед запуском необходимо клонировать репозиторий с исходным ходом
и перейти в директорию с проектом:
```bash
$ git clone https://github.com/lunn06/wallet && cd wallet
```
### Docker Compose
#### Зависимости
- Docker Compose

Для запуска программы необходимо скопировать example-deploy.yaml:
```bash
$ cp configs/example-deploy.yaml configs/deploy.yaml
```
И изменить database.password на реальный пароль и продублировать его в POSTGRES_PASSWORD в deployments/docker-compose.yaml.
Для запуска необходимо ввести:
```bash
$ docker compose -f deployments/docker-compose.yaml up -d
```

### Запуск вручную 
#### Зависимости
- Go 1.23.5
- Postgres 17

Установку и запуск Postgresql 17 выполнить в соответствии с [официальной документацией](https://www.postgresql.org/download/).
Данные для запуска продублировать в main.yaml(его скопировать как deploy.yaml в пункте выше)

Скачиваем зависимости и собираем проект:
```bash
$ go mod download && go build -o wallet-backend cmd/app/main.go
```

Для запуска в "release" режиме:
```bash
$ GIN_MODE=release ./wallet-backend
```

Для запуска в "debug" режиме и вывода логов в консоль:
```bash
$ LOG_LEVEL=debug ./wallet-backend
```
