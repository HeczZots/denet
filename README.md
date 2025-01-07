## Шаги запуска
```
git clone https://github.com/HeczZots/denet.git
cd denet
docker-compose up --build
```
## Перед началом
```
Сначала нужно зарегестрироваться по ендпоинту 

localhost:8080/registration

с телом запроса в формате

{
  "login": "user",
  "password": "password"
}

Далее получаем токен его перед каждым запросом передаем в header Authorization

Теперь все ендпоинты доступны
```
## Описание задания

```

Есть задание с подпиской на мой телеграмм чат: https://t.me/hft_alerts_station

Далее вы стучитесь в этот бот @getmyid_bot

получаете id

ендпоинт localhost:8080/users/{id}/task/complete

в теле запроса отправляете например 

{
    "telegram_user_id":1834623444
}

```


## Примеры 
1. Регистрация
```
curl -X POST http://localhost:8080/registration \
-d '{
  "login": "user123",
  "password": "password123"
}'
```
2. Логин
```
curl -X POST http://localhost:8080/login \
-d '{
  "login": "user123",
  "password": "password123"
}'
```
3. Завершение задания
```
curl -X POST http://localhost:8080/users/1/task/complete \
-H "Authorization: токен" \
-d '{
  "telegram_user_id": 1834623444
}'
```
4. Таблица лучших
```
curl -X GET http://localhost:8080/users/leaderboard \
-H "Authorization: токен" \
```