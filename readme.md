Калькулятор на го

Всем привет, реализован максимальный потенциал из возможного для меня


одна страница html является единстенной страницой данного проекта



cmd/main.go просто запуск сервера

html страница html лежит

server.go открывает страницу с минимальным функционалом где считается выражения, пока не ввёл выражение правильное не посчитает
calculator.go считает то, что вы набрали


запрос с сервера -> server.go -> calculator.go -> server.go(плюс запись в базе данных, надеюсь работает) -> пользователь
