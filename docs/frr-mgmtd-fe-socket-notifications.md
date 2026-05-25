# FRR mgmtd frontend socket notifications

Документ про `/run/frr/mgmtd_fe.sock`: что это за socket, как через него читать
уведомления, как это делает `fe_client.py`, как это делает наш Go watcher, что
реально прилетает, что не прилетает, и почему важно держать одно соединение
открытым.

## Короткая суть

`/run/frr/mgmtd_fe.sock` - это Unix stream socket frontend API к `mgmtd`.

Через него frontend-клиент может:

- создать frontend session;
- читать данные через `GET_DATA`;
- подписаться на datastore/YANG notifications через `NOTIFY_SELECT`;
- получать async `NOTIFY`;
- делать edit/commit/lock/RPC, если реализовать соответствующие message codes.

Для watcher важны три правила:

1. Открыть socket один раз.
2. Создать одну mgmtd session.
3. Дальше читать и писать сообщения по тому же socket, пока connection живой.

Если на каждый poll делать новый `connect -> session -> get -> close`, `mgmtd`
начинает логировать много строк вида:

```txt
FE-ADAPTER-CONN: mgmt_msg_read: got EOF/disconnect
```

Это не полезная диагностика, а следствие того, что клиент постоянно открывает и
закрывает frontend connection.

## Где это живет в FRR

Главные файлы в `third_party/frr`:

- `tests/topotests/lib/fe_client.py` - эталонный Python frontend client.
- `lib/mgmt_msg_native.h` - native protocol structs, message codes, formats.
- `mgmtd/mgmt_fe_adapter.c` - frontend adapter в `mgmtd`, обработка frontend sessions.
- `mgmtd/mgmt_be_adapter.c` - backend adapter, прием notify от daemon/backend clients.
- `lib/mgmt_be_client.c` - backend client side, отправка notify/get-data replies.

Наш код:

- `go-watcher/main.go` - Go watcher для `make go-watch`.
- `Makefile`, target `go-watch` - собирает binary под архитектуру контейнера,
  копирует в контейнер и запускает внутри него.

## Socket path

В контейнере обычно используется:

```txt
/run/frr/mgmtd_fe.sock
```

В Python-клиенте FRR default указан как:

```txt
/var/run/frr/mgmtd_fe.sock
```

На Linux `/var/run` часто является ссылкой на `/run`, поэтому оба пути обычно
ведут к одному месту.

Проверить:

```sh
docker exec -it frr-playground ls -l /run/frr/mgmtd_fe.sock /var/run/frr/mgmtd_fe.sock
```

## Общая модель

`mgmtd` стоит посередине:

```txt
frontend client
  |
  | /run/frr/mgmtd_fe.sock
  v
mgmtd frontend adapter
  |
  | internal mgmtd transactions / backend fanout
  v
backend clients / FRR daemons
```

Frontend client - это наш watcher, `fe_client.py`, `vtysh` frontend или другой
клиент.

Backend clients - это FRR daemons и mgmtd backend clients, которые публикуют
config/state/notify capabilities.

Важно: frontend socket не является текстовым CLI. Это бинарный stream protocol.
Нельзя просто сделать `cat /run/frr/mgmtd_fe.sock` и увидеть JSON. JSON лежит
внутри native messages как payload.

## Wire format

Каждое сообщение в socket имеет внешний frame:

```txt
4 bytes marker
4 bytes total message size
N bytes native message body
```

Для native protocol marker:

```txt
01 23 23 23
```

В Python это записано как:

```py
MGMT_MSG_MARKER_NATIVE = b"\001###"
```

В Go это читается little-endian uint32:

```go
const MGMT_MSG_MARKER_NATIVE = 0x23232301
```

После marker идет `uint32` длина всего сообщения, включая первые 8 bytes.

То есть read loop всегда такой:

1. прочитать 8 bytes;
2. проверить marker;
3. взять length;
4. дочитать `length - 8` bytes body;
5. распарсить native header внутри body.

## Native message header

Все native messages начинаются с `mgmt_msg_header`.

Из `lib/mgmt_msg_native.h`:

```c
struct mgmt_msg_header {
	uint16_t code;
	uint16_t resv;
	uint32_t vsplit;
	uint64_t refer_id;
	uint64_t req_id;
};
```

Размер header - 24 bytes.

Поля:

- `code` - тип сообщения.
- `resv` - reserved, обычно zero.
- `vsplit` - если payload состоит из двух строковых частей, тут граница первой.
- `refer_id` - session id, transaction id или другой связанный id.
- `req_id` - request id, по нему связывается request и reply.

Python формат:

```py
MSG_FMT_HDR = "=H2xIQQ"
```

Go читает так:

```go
code := binary.LittleEndian.Uint16(buf[8:10])
vsplit := binary.LittleEndian.Uint32(buf[12:16])
msgSessID := binary.LittleEndian.Uint64(buf[16:24])
reqID := binary.LittleEndian.Uint64(buf[24:32])
```

Почему offsets начинаются с `8`: первые 8 bytes это внешний frame
`marker + length`.

## Основные message codes

Для watcher важны:

| Code | Name | Direction | Для чего |
| ---: | --- | --- | --- |
| `0` | `ERROR` | mgmtd -> client | Ошибка на request |
| `2` | `TREE_DATA` | mgmtd -> client | Ответ на `GET_DATA` |
| `3` | `GET_DATA` | client -> mgmtd | Запрос config/state по XPath |
| `4` | `NOTIFY` | mgmtd -> client | Async notification |
| `9` | `NOTIFY_SELECT` | client -> mgmtd | Подписка на selectors |
| `10` | `SESSION_REQ` | client -> mgmtd | Создать или удалить session |
| `11` | `SESSION_REPLY` | mgmtd -> client | Ответ на session request |

В `lib/mgmt_msg_native.h` эти codes помечены как `Public API`.

## Session lifecycle

Первое сообщение после connect - `SESSION_REQ`.

Создание session:

```txt
code = SESSION_REQ
refer_id = 0
req_id = client-side request id
notify_format = JSON или 0/default
client_name = null-terminated string
```

`refer_id = 0` означает create.

`mgmtd` отвечает `SESSION_REPLY`:

```txt
code = SESSION_REPLY
refer_id = new session id
req_id = request id from SESSION_REQ
created = true
```

После этого frontend client должен использовать session id в дальнейших
requests.

Корректное удаление session в Python:

```py
# sending session_req with a non-zero session ID destroys the session.
mdata, _ = self.get_native_msg_header(MSG_CODE_SESSION_REQ)
mdata += struct.pack(MSG_FMT_SESSION_REQ, MSG_FORMAT_JSON)
self.send_native_msg(mdata)
self.sock.close()
```

Там `get_native_msg_header()` кладет текущий `self.sess_id` в `refer_id`, поэтому
для `mgmtd` это destroy session.

Наш Go watcher сейчас проще: при остановке процесса закрывает socket. Это
допустимо для watcher, но если делать библиотеку, лучше добавить clean close
через `SESSION_REQ` с `refer_id = sessID`.

## Как работает `fe_client.py`

Файл:

```txt
third_party/frr/tests/topotests/lib/fe_client.py
```

Подключение:

```py
sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
sock.connect_ex(str(spath))
```

Дальше создается `Session(sock)`.

В `Session.__init__()`:

1. `self.sess_id = 0`;
2. собирается native header `SESSION_REQ`;
3. добавляется fixed body `MSG_FMT_SESSION_REQ`;
4. добавляется `client_name`;
5. сообщение отправляется;
6. клиент ждет `SESSION_REPLY`;
7. `self.sess_id` становится равен `refer_id` из reply.

Скелет:

```py
mdata, _ = self.get_native_msg_header(MSG_CODE_SESSION_REQ)
mdata += struct.pack(MSG_FMT_SESSION_REQ, MSG_FORMAT_JSON)
mdata += "test-client".encode("utf-8") + b"\x00"
self.send_native_msg(mdata)

mhdr, mfixed, mdata = self.recv_native_msg()
self.sess_id = mhdr[HDR_FIELD_SESS_ID]
```

### `GET_DATA` в Python

`get_data()` собирает:

```txt
header:
  code = GET_DATA
  refer_id = session id
  req_id = next request id

fixed body:
  result_type = JSON
  flags = STATE/CONFIG/EXACT

payload:
  xpath + "\0"
```

Флаги:

```py
GET_DATA_FLAG_STATE = 0x1
GET_DATA_FLAG_CONFIG = 0x2
GET_DATA_FLAG_EXACT = 0x4
```

Пример:

```py
result = sess.get_data(
    query,
    data=not args.config_only,
    config=(args.both or args.config_only),
)
```

Ответ приходит как `TREE_DATA`, внутри payload находится JSON string с
null-terminator.

### Notifications в Python

`fe_client.py --listen` работает так:

```py
if args.listen is not None:
    if args.listen:
        sess.add_notify_select(True, args.listen)
    while ...:
        result_type, op, xpath, notif = sess.recv_notify()
```

Тонкий момент:

- `--listen` без XPath не отправляет `NOTIFY_SELECT`;
- `--listen /some/xpath` отправляет `NOTIFY_SELECT` с selector;
- `recv_notify()` просто ждет `MSG_CODE_NOTIFY`.

Это важно, потому что session без selectors работает по особому правилу для
YANG-defined notifications.

## `NOTIFY_SELECT`

`NOTIFY_SELECT` нужен, чтобы сказать `mgmtd`: "отправляй мне notifications,
которые соответствуют этим XPath prefixes".

Структура:

```c
struct mgmt_msg_notify_select {
	struct mgmt_msg_header;
	uint8_t replace;
	uint8_t get_only;
	uint8_t resv2[6];
	alignas(8) char selectors[];
};
```

В Python:

```py
mdata, _ = self.get_native_msg_header(MSG_CODE_NOTIFY_SELECT)
mdata += struct.pack(MSG_FMT_NOTIFY_SELECT, replace)

for xpath in notif_xpaths:
    mdata += xpath.encode("utf-8") + b"\x00"
```

Смысл:

- `replace = true` - заменить старые selectors;
- `replace = false` - добавить selectors;
- `selectors[]` - список null-terminated XPath strings.

В `mgmt_fe_adapter.c` selectors валидируются через northbound nodes. Если selector
не резолвится в node, mgmtd вернет error.

## Что значит "без selectors"

Если frontend session вообще не отправляла `NOTIFY_SELECT`, то
`session->notify_xpaths == NULL`.

В `mgmt_fe_adapter.c` есть отдельная логика:

```c
/*
 * Send all YANG defined notifications to all sesisons with *no*
 * selectors as well (i.e., original NETCONF/RESTCONF notification
 * scheme).
 */
if (CHECK_FLAG(nb_node->snode->nodetype, LYS_NOTIF)) {
	LIST_FOREACH (adapter, &fe_adapters, link) {
		LIST_FOREACH (session, &adapter->sessions, link) {
			if (session->notify_xpaths)
				continue;
			...
		}
	}
}
```

Перевод:

- если прилетела настоящая YANG notification;
- и session не имеет selectors;
- `mgmtd` отправит эту notification такой session.

Поэтому для режима "слушать все YANG notifications" правильное поведение - не
посылать `NOTIFY_SELECT` вообще.

Но это касается именно YANG-defined notifications, то есть узлов `notification`
в YANG schema.

## Что прилетает в `NOTIFY`

`NOTIFY` message:

```c
struct mgmt_msg_notify_data {
	struct mgmt_msg_header;
	uint8_t result_type;
	uint8_t op;
	uint8_t resv2[6];
	alignas(8) char data[];
};
```

`data[]` содержит две части:

```txt
xpath\0json-or-empty\0
```

Граница между `xpath` и `json` задается `vsplit`.

В Python:

```py
vsplit = mhdr[HDR_FIELD_VSPLIT]
xpath = mdata[: vsplit - 1].decode("utf-8")
notif = mdata[vsplit:-1].decode("utf-8")
```

В Go:

```go
payload := buf[40 : 8+bodyLen]
xpath, data := splitNotificationPayload(payload, int(vsplit))
```

Операции:

| Op | Name | Значение |
| ---: | --- | --- |
| `0` | `NOTIFICATION` | Обычная YANG notification |
| `1` | `DS_REPLACE` | Datastore node replace |
| `2` | `DS_DELETE` | Datastore node delete |
| `3` | `DS_PATCH` | Datastore patch |
| `4` | `DS_GET_SYNC` | Initial sync dump при selector subscription |

В выводе `fe_client.py --datastore` это отображается как:

```txt
#OP=PATCH: /some/xpath
{...}
```

## Что не прилетает в `NOTIFY`

Не каждое изменение FRR config становится async `NOTIFY`.

Чтобы что-то прилетело как настоящая YANG notification, должны совпасть условия:

1. В YANG module должен быть `notification` node.
2. Backend/client должен реально отправить notify.
3. `mgmtd` должен сматчить notification с sessions/selectors.

Если модуль просто меняет running config, но не имеет YANG `notification`, то
session без selectors ничего не получит.

Пример из нашего кейса:

- `workflow-07-route-maps` меняет route-map config;
- `frr-route-map.yang` не содержит `notification` для изменения route-map;
- `lib/routemap_northbound.c` вызывает internal dependency notifications,
  но это не то же самое, что frontend `MSG_CODE_NOTIFY` для `/run/frr/mgmtd_fe.sock`.

Итог: route-map config change не прилетает как async YANG notification на
frontend socket.

Поэтому наш watcher отдельно делает `GET_DATA` polling для:

```txt
/frr-route-map:lib
```

и сравнивает JSON snapshots.

## Как работает наш Go watcher

Файл:

```txt
go-watcher/main.go
```

Target:

```sh
make go-watch
```

Что делает target:

1. собирает Go binary под Linux arch контейнера;
2. копирует binary в контейнер;
3. запускает `/usr/local/bin/frr-watcher` внутри контейнера.

Watcher внутри контейнера:

1. делает `net.Dial("unix", "/run/frr/mgmtd_fe.sock")`;
2. отправляет `SESSION_REQ`;
3. читает `SESSION_REPLY`;
4. не отправляет `NOTIFY_SELECT`, чтобы слушать all YANG-defined notifications;
5. держит socket открытым;
6. периодически отправляет `GET_DATA` по тому же socket;
7. читает из того же socket и async `NOTIFY`, и replies `TREE_DATA`.

### Почему важно одно соединение

Правильная модель:

```txt
connect once
SESSION_REQ once
loop:
  send GET_DATA when needed
  read TREE_DATA replies
  read NOTIFY async messages
disconnect only on exit/error
```

Неправильная модель:

```txt
loop:
  connect
  SESSION_REQ
  GET_DATA
  read reply
  close
```

Неправильная модель засоряет `frr.log`, потому что `mgmtd` видит постоянные
frontend disconnects:

```txt
MGMTD: FE-ADAPTER-CONN: mgmt_msg_read: got EOF/disconnect
```

После фикса Go watcher больше не открывает отдельную session на каждый poll.
Он отправляет `sendGetData(conn, sessID, reqID, ...)` в существующий socket.

### Pending requests

Поскольку socket один, ответы приходят в общем потоке сообщений.

Поэтому watcher хранит:

```go
pendingRequests := map[uint64]string{}
```

Когда отправляем `GET_DATA`, кладем:

```go
pendingRequests[reqID] = xpath
```

Когда приходит `TREE_DATA`, смотрим `reqID`, находим исходный XPath и понимаем,
какой parser/comparison применить.

Это важно, потому что в том же stream могут прийти:

- `TREE_DATA` на route-map polling;
- `TREE_DATA` на backend state polling;
- async `NOTIFY`;
- `ERROR`.

## Что сейчас ловит `make go-watch`

### 1. Настоящие YANG notifications

Так как watcher не отправляет `NOTIFY_SELECT`, он слушает все YANG-defined
notifications, которые `mgmtd` рассылает sessions без selectors.

Вывод:

```txt
=== NOTIFICATION RECEIVED ===
Session ID: ...
Req ID: ...
Result Type: ...
Op: ...
#OP=...
...
```

### 2. Route-map config snapshots

Для route-map watcher делает polling:

```txt
/frr-route-map:lib
datastore = running
flags = config
```

Когда JSON отличается от предыдущего snapshot, watcher печатает:

```txt
=== ROUTE-MAP CONFIG CHANGE ===
#OP=REPLACE: /frr-route-map:lib
{...}
```

или:

```txt
=== ROUTE-MAP CONFIG CHANGE ===
#OP=DELETE: /frr-route-map:lib
```

Это не native async notification. Это наш derived event из `GET_DATA` polling.
Он нужен именно потому, что route-map changes не приходят как frontend notify.

### 3. Backend state snapshots

Watcher также читает:

```txt
/frr-backend:clients/client
datastore = operational
flags = state
```

Это помогает видеть backend config versions, например изменение
`running-config-version` у daemon clients.

Вывод:

```txt
=== BACKEND STATE CHANGE ===
#OP=PATCH: /frr-backend:clients/client
{...}
```

Это тоже derived event из polling, не async notify.

## Как руками проверить socket через Python

В контейнере или окружении, где есть FRR source:

```sh
python3 third_party/frr/tests/topotests/lib/fe_client.py \
  --server /run/frr/mgmtd_fe.sock \
  --query /frr-route-map:lib \
  --config-only
```

Слушать notifications без selectors:

```sh
python3 third_party/frr/tests/topotests/lib/fe_client.py \
  --server /run/frr/mgmtd_fe.sock \
  --listen \
  --notify-count 0
```

Слушать с selector:

```sh
python3 third_party/frr/tests/topotests/lib/fe_client.py \
  --server /run/frr/mgmtd_fe.sock \
  --listen /some/valid/xpath \
  --notify-count 0
```

Слушать datastore-style output:

```sh
python3 third_party/frr/tests/topotests/lib/fe_client.py \
  --server /run/frr/mgmtd_fe.sock \
  --listen /some/valid/xpath \
  --datastore \
  --notify-count 0
```

Если selector не резолвится в northbound node, `mgmtd` вернет error.

## Как руками проверить наш watcher

Запустить watcher:

```sh
make go-watch
```

В другом терминале:

```sh
docker exec frr-playground vtysh \
  -c "conf t" \
  -c "route-map WATCH_TEST permit 10"
```

Ожидаемый вывод watcher:

```txt
=== ROUTE-MAP CONFIG CHANGE ===
#OP=REPLACE: /frr-route-map:lib
{"frr-route-map:lib":{"route-map":[...]}}
```

Cleanup:

```sh
docker exec frr-playground vtysh \
  -c "conf t" \
  -c "no route-map WATCH_TEST"
```

Ожидаемый вывод:

```txt
=== ROUTE-MAP CONFIG CHANGE ===
#OP=DELETE: /frr-route-map:lib
```

Проверить, что нет EOF spam:

```sh
docker exec frr-playground sh -lc \
  'tail -n 120 /var/log/frr/frr.log | grep "FE-ADAPTER-CONN: mgmt_msg_read: got EOF/disconnect" || true'
```

Если watcher работает правильно, во время обычного polling эта команда не должна
показывать поток новых строк.

## Как читать сообщения самому

Минимальный алгоритм для своего клиента:

1. `socket(AF_UNIX, SOCK_STREAM)`.
2. `connect("/run/frr/mgmtd_fe.sock")`.
3. Send native `SESSION_REQ`.
4. Read `SESSION_REPLY`, сохранить `sessID`.
5. Если нужны filtered datastore notifications, send `NOTIFY_SELECT`.
6. Если нужны all YANG-defined notifications, не send `NOTIFY_SELECT`.
7. В одном read loop:
   - read frame header;
   - parse native header;
   - switch по `code`;
   - для `NOTIFY` разобрать `vsplit`;
   - для `TREE_DATA` сопоставить `req_id` с pending request;
   - для `ERROR` вывести error и удалить pending request.
8. Не закрывать socket между poll requests.

Псевдокод:

```txt
conn = connect(unix_socket)
send SESSION_REQ
sessID = read SESSION_REPLY

while running:
  if time_to_poll:
    reqID = next()
    pending[reqID] = xpath
    send GET_DATA(sessID, reqID, xpath)

  msg = read_next_message(conn)

  switch msg.code:
    NOTIFY:
      print notification
    TREE_DATA:
      xpath = pending[msg.req_id]
      delete pending[msg.req_id]
      process data for xpath
    ERROR:
      print error
```

## Что можно улучшить дальше

Для production-quality клиента полезно добавить:

- clean session destroy через `SESSION_REQ` с `refer_id = sessID`;
- backoff на reconnect;
- bounded pending request map;
- metrics по request latency;
- отдельный parser для `ERROR`;
- настройку poll interval через env/flag;
- allowlist/denylist для derived polling output;
- отдельный режим "только native NOTIFY без polling";
- отдельный режим "poll selected XPaths";
- unit tests для framing/parser functions.

## Практический вывод

`mgmtd_fe.sock` - это не лог-файл и не CLI. Это stateful binary API.

Правильная работа с ним:

- один long-lived Unix socket;
- одна frontend session;
- request/reply correlation через `req_id`;
- async messages в том же stream;
- selectors только если реально нужен filtered subscription;
- без selectors можно слушать all YANG-defined notifications;
- config changes без YANG notification надо читать через `GET_DATA`/polling или
  другим audit-механизмом.

Для нашего `make go-watch` это означает:

- native notifications читаются из socket как есть;
- route-map workflow ловится через `GET_DATA /frr-route-map:lib`;
- socket не открывается заново на каждый poll;
- `frr.log` не должен получать EOF spam от watcher.

