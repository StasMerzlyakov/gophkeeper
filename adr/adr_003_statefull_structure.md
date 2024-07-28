# ADR 003

## Структура базы данных
- статус: proposed
- 2024-07-17

## Контекст
В базе данных содержится только постоянная информация.

## Принятое решение  (Alt-D)
```plantuml
@startuml
class userInfo{
    userId bigserial                   //  PK
	email text not null                // email, уникальный
	pass_hash text not null            // Хэш от пароля
	pass_salt text not null            // Солья для вычисления хэша от пароля
	otp_key text not null              // Зашифрованный на ServerSecret OTP пароль пользователя
	master_key text not null           // Зашифрованный на MasterKeyPass ключи шифрования данных пользователя
	master_hint text not null          // Напоминалка для пользователя для восстановления MasterKey
	hello_encrypted text not null      // Зашифрованная на MasterKey 'Hello from GophKeeper!!!'. Используется для проверки правильности ввода MasterKeyPass
	primary key(userId)
    unique (email)
}
@enduml
```

## Ссылки

[https://habr.com/ru/articles/747348/](https://habr.com/ru/articles/747348/)