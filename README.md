# Серверное p2p приложение

Выполнено в рамках курсовой работы в 2017 году. 
Приложение использует т.н. хордовый алгоритм для 
распределенного хранения истории сообщений и базы
пользователей.

## Использование

```
chat -host="192.168.56.101" -remote="192.168.56.2" -port=7777
```

Параметр `host` отвечает за адрес данного сервера 
(какой адаптер будет исопльзовать). По умолчанию равен `127.0.0.1`

Параметр `remote` отвечает за удаленный хост, к которому будет
присоединен данный. По умолчанию равен пустому значению, т.е.
не подключается ни к кому.

Параметр `port` отвечает за порт, который будет слушать сервер.
По умолчанию используется порт `7777`