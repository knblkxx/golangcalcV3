# Калькулятор 
Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления.
### Как запустить?
```cmd
git clone https://github.com/knblkxx/golang-calc
cd yandex-go
go run main/main.go
```
##### Приложение использует post запросы с json формата {"expression": "ваше выражение"}
### Пример работы:
##### status code 200 (вычислено успешно)
```bash
curl --location http://localhost:8080/api/v1/calculate --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}'
```

##### status code 422 (присутствует деление на 0)
```bash
curl --location http://localhost:8080/api/v1/calculate --header 'Content-Type: application/json' --data '{"expression": "2+2*2/0"}'
```
### Структура программы
```
ExpressionHandler <- Клиент -> Оркестр <-> Агент
```
