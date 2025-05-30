# test-task-2905

#### Для запуска понадобится:
Установленный Docker и Docker Compose

### 1. Клонируйте репозиторий 
```
git clone <ваш-репозиторий>
cd <папка-проекта>
```
### 2. Запустите сервисы через Docker Compose 
```
docker compose up -d --remove-orphans 
```
### P.S
Для настройки переменных окружения:

Файл docker-compose.yml уже содержит настройки:

PostgreSQL:
  Порт: 5432
  Пользователь: postgres
  Пароль: qwerty
  База данных: test-task-db
Сервер:
  Порт: 8080
  Подключение к БД: postgresql://postgres:qwerty@postgres:5432/test-task-db?sslmode=disable
  
Если нужно изменить параметры, отредактируйте docker-compose.yml и перезапустите сервисы.

### Команды согласно ТЗ:
curl -X POST http://localhost:8080/quotes -H "Content-Type: application/json" -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'

curl http://localhost:8080/quotes curl http://localhost:8080/quotes/random

curl http://localhost:8080/quotes?author=Confucius


