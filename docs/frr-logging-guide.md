# FRR logging guide

Этот файл про практичные способы читать логи FRR на сервере или в контейнере.
Отдельно: `log commands` логирует введенные CLI-команды, но не является механизмом
получения diff `running-config`.

## Быстрый вариант для проекта

Включить логирование команд и подробный вывод:

```frr
conf t
 log file /var/log/frr/frr.log debugging
 log stdout debugging
 log commands
 debug mgmt frontend
 debug mgmt backend
 debug northbound callbacks configuration
end
write
```

Читать лог из контейнера:

```sh
docker exec -it frr-playground tail -f /var/log/frr/frr.log
```

Или читать stdout контейнера:

```sh
docker logs -f frr-playground
```

## Основные места, где читать FRR logs

### 1. Log file

Самый удобный и предсказуемый вариант:

```sh
tail -f /var/log/frr/frr.log
```

Для Docker-контейнера:

```sh
docker exec -it frr-playground tail -f /var/log/frr/frr.log
```

Настройка во FRR:

```frr
conf t
 log file /var/log/frr/frr.log debugging
end
write
```

### 2. Docker logs

Если в FRR включен `log stdout`, логи можно читать через Docker:

```sh
docker logs -f frr-playground
```

В этом проекте `configs/frr/frr.conf` уже содержит:

```frr
log stdout
log file /var/log/frr/frr.log
log record-priority
log timestamp precision 3
```

### 3. systemd journal

На обычном сервере с systemd:

```sh
journalctl -u frr -f
```

Для supervisor-процесса FRR:

```sh
journalctl -u watchfrr -f
```

Посмотреть последние события:

```sh
journalctl -u frr --since "10 min ago"
```

### 4. syslog

Можно отправлять FRR logs в syslog:

```frr
conf t
 log syslog debugging
end
write
```

Потом читать системные логи, в зависимости от дистрибутива:

```sh
journalctl -f
tail -f /var/log/syslog
tail -f /var/log/messages
```

### 5. vtysh terminal monitor

Можно смотреть live logs прямо внутри `vtysh`:

```sh
docker exec -it frr-playground vtysh
```

Внутри `vtysh`:

```frr
terminal monitor
```

Для конкретного daemon:

```frr
terminal monitor zebra
terminal monitor bgpd
terminal monitor mgmtd
```

Отключить:

```frr
no terminal monitor
```

Важно: `log monitor` в текущем FRR помечен как deprecated и в коде ничего не
делает. Для live-вывода нужен именно `terminal monitor`.

## Как ловить введенные команды

Для логирования CLI-команд нужен `log commands`:

```frr
conf t
 log commands
 log file /var/log/frr/frr.log debugging
end
write
```

После этого читать:

```sh
docker exec -it frr-playground tail -f /var/log/frr/frr.log
```

Для фильтрации только строк с командами можно начать с простого grep:

```sh
docker exec -it frr-playground sh -lc 'tail -f /var/log/frr/frr.log | grep -i command'
```

Отключить логирование команд:

```frr
conf t
 no log commands
end
write
```

## Принудительное log commands при запуске daemon

У FRR есть режим запуска daemon с `--log commands`. В этом режиме команды
логируются принудительно, а `no log commands` из конфигурации уже нельзя
использовать для отключения.

Это настраивается не через `vtysh`, а через параметры запуска daemon, например
через service/unit или `/etc/frr/daemons`, в зависимости от окружения.

## Проверка текущей настройки

Посмотреть, куда FRR сейчас пишет логи:

```sh
docker exec -it frr-playground vtysh -c "show logging"
```

Посмотреть включенные debug-флаги:

```sh
docker exec -it frr-playground vtysh -c "show debugging"
```

Посмотреть активный running config, включая logging/debug:

```sh
docker exec -it frr-playground vtysh -c "show running-config"
```

## Полезные debug-команды для mgmtd/northbound

Для отладки mgmtd и northbound callbacks:

```frr
conf t
 debug mgmt frontend
 debug mgmt backend
 debug mgmt datastore
 debug mgmt transaction
 debug northbound callbacks configuration
 debug northbound callbacks state
 debug northbound callbacks rpc
 debug northbound callbacks notify
end
write
```

Отключение:

```frr
conf t
 no debug mgmt frontend
 no debug mgmt backend
 no debug mgmt datastore
 no debug mgmt transaction
 no debug northbound callbacks configuration
 no debug northbound callbacks state
 no debug northbound callbacks rpc
 no debug northbound callbacks notify
end
write
```

## Важное различие

`log commands` отвечает на вопрос:

> Какие CLI-команды были введены?

Например:

```frr
route-map RM-IN permit 10
set local-preference 200
no route-map RM-IN permit 10
```

Но `log commands` не отвечает полноценно на вопрос:

> Что именно изменилось в `running-config` как structured data?

Для отслеживания изменений конфигурации как данных нужны другие механизмы:

- чтение `running-config`;
- mgmtd/northbound;
- gRPC get/polling;
- notify/subscription, если конкретный YANG-модуль реально отправляет notifications;
- внешний audit вокруг `vtysh` или workflow-команд.

То есть для аудита CLI лучше `log commands`, а для программного отслеживания
состояния FRR лучше читать данные через mgmtd/northbound/gRPC или сравнивать
`running-config`.

