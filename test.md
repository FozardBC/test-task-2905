Мини-сервис “Цитатник”
Цель: Проверить базовые навыки кандидата:
работа с HTTP, структура кода, умение разбираться с задачей и писать чистый код.
Задание:
- Реализуйте REST API-сервис на Go для хранения и управления цитатами.
- Разместите решение в публичном Github репозитории
- Ссылку на репозиторий отправьте менеджеру
Функциональные требования:
1. Добавление новой цитаты (POST /quotes)
2. Получение всех цитат (GET /quotes)
3. Получение случайной цитаты (GET /quotes/random)
4. Фильтрация по автору (GET /quotes?author=Confucius)
5. Удаление цитаты по ID (DELETE /quotes/{id})
Проверочные команды (curl):
curl -X POST http://localhost:8080/quotes \ -H "Content-Type: application/json" \ -d
'{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
curl http://localhost:8080/quotes curl http://localhost:8080/quotes/random
curl http://localhost:8080/quotes?author=Confucius
curl -X DELETE http://localhost:8080/quotes/1
Технические требования:
- Хранить данные можно в памяти.
- Использовать только стандартные библиотеки Go (максимум gorilla/mux)
- Обязательно: README.md с инструкцией запуска
- Желательно: unit-тесты.
Критерии приемки:
1. Код собирается без ошибок
2. Все функции работают через curl (см. выше).
3. Присутствует README.
4. Структура проекта читаема.
5. (Опционально) Наличие тестов