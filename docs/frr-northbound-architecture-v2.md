---
name: frr-northbound-architecture
description: "FRRNorthbound architecture, notification mechanisms, MGMTd, sysrepo"
metadata: 
  node_type: memory
  type: project
  originSessionId: 0f4e7fb6-98f4-46fb-b6be-ce15aef46a0f
---

# FRR Northbound Architecture - Полное документирование

## Содержание

1. [Общая архитектура FRR](#общая-архитектура-frr)
2. [Northbound API](#northbound-api)
3. [Механизмы уведомлений](#механизмы-уведомлений)
4. [Межмодульное взаимодействие](#межмодульное-взаимодействие)
5. [Проблемы с BGP](#проблемы-s-bgp)
6. [Технические детали реализаций](#технические-детали-реализаций)
7. [Итоговая сводка решений](#итоговая-сводка-решений)

---

## Общая архитектура FRR

### Компоненты системы

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FRR Architecture                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐   │
│  │   vtysh     │  │   mgmtd     │  │   bgpd      │  │   other daemons │   │
│  │ (CLI shell) │  │ (central   │  │ (BGP)       │  │   (zebra, ripd, │   │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘   │
│         │                │                │                   │              │
│         └────────────────┴────────────────┴───────────────────┘              │
│                            │                                                 │
│                    ┌───────▼───────┐                                         │
│                    │  Northbound   │                                         │
│                    │    Layer      │                                         │
│                    └───────────────┘                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Поток данных конфигурации

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Configuration Flow                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. User input (vtysh CLI)                                                  │
│     ↓                                                                        │
│  2. vtysh sends command to daemon via unix socket                          │
│     ↓                                                                        │
│  3. Daemon applies to in-memory configuration (running_config)             │
│     ↓                                                                        │
│  4. "show running-config" reads from running_config (memory)               │
│     ↓                                                                        │
│  5. "write memory" writes running_config to /etc/frr/*.conf                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Northbound API

### Что такое Northbound API?

Northbound API - это интерфейс для внешнего управления конфигурацией FRR. Он предоставляет стандартизированный способ получения и изменения конфигурации через различные протоколы.

### Архитектура Northbound

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Northbound Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐              │
│  │   Client     │       │  Northbound  │       │   YANG       │              │
│  │  (Go/Python) │──────▶│   Layer      │──────▶│  Modules     │              │
│  └──────────────┘       └──────┬───────┘       └──────────────┘              │
│                                │                                            │
│                     ┌──────────▼──────────┐                                │
│                     │  northbound.c       │                                │
│                     │  - nb_config        │                                │
│                     │  - nb_node          │                                │
│                     │  - callbacks        │                                │
│                     └──────────┬──────────┘                                │
│                                │                                            │
│                     ┌──────────▼──────────┐                                │
│                     │   Daemons           │                                │
│                     │   - bgpd            │                                │
│                     │   - zebra           │                                │
│                     │   - isisd           │                                │
│                     └─────────────────────┘                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Типы northbound модулей

Каждый daemon регистрирует свой `frr_yang_module_info`:

```c
static const struct frr_yang_module_info *const daemon_yang_modules[] = {
    &frr_backend_info,        // Для mgmtd
    &frr_filter_info,
    &frr_interface_info,
    &frr_<daemon>_info,       // Основной модуль (ПОЛНЫЙ)
    &frr_route_map_info,
    &frr_<daemon>_route_map_info, // Модуль для route-map (частичный)
    &frr_vrf_info,
};
```

### Northbound Callbacks

```c
struct nb_callbacks {
    int (*create)(struct nb_cb_create_args *args);     // Создание объекта
    int (*modify)(struct nb_cb_modify_args *args);     // Изменение значения
    int (*destroy)(struct nb_cb_destroy_args *args);   // Удаление
    
    enum nb_error (*get)(const struct nb_node *nb_node, 
                         const void *list_entry, 
                         struct lyd_node *parent);     // Получение данных
    
    int (*rpc)(struct nb_cb_rpc_args *args);          // Вызов RPC
    void (*notify)(struct nb_cb_notify_args *args);   // Обработка уведомлений
};
```

---

## Механизмы уведомлений

### 1. MGMTd (Central Management Daemon)

#### Межмодульные соединения:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MGMTd Architecture                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Frontend (FE)                                   Backend (BE)               │
│  ┌──────────────┐      unix socket       ┌──────────────────┐               │
│  │ mgmt_fe_     │────── mgmtd_fe.sock ──▶│ mgmt_be_         │               │
│  │   adapter.c  │      (port 50051-64)   │   adapter.c      │               │
│  └──────────────┘                        └──────────────────┘               │
│        │                                         │                         │
│        │notifications                            │config/oper                │
│        ▼                                         ▼                         │
│  ┌──────────────┐                        ┌──────────────────┐               │
│  │  gRPC (50051)│                        │   bgpd (port     │               │
│  │   client     │                        │    50052)        │               │
│  └──────────────┘                        └──────────────────┘               │
│                                                                             │
│  Datastores:                                                                │
│  - MGMTD_DS_RUNNING       - Текущая конфигурация                           │
│  - MGMTD_DS_CANDIDATE     - Кандидат на коммит                              │
│  - MGMTD_DS_OPERATIONAL   - Операционные данные                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Модули MGMTd:

- **mgmt_main.c** - Главный демон
- **mgmt_ds.c** - Datastores (running, candidate, operational)
- **mgmt_fe_adapter.c** - Frontend клиенты (gRPC, vtysh)
- **mgmt_be_adapter.c** - Backend клиенты (bgpd, zebra, ripd и т.д.)
- **mgmt_txn.c** - Транзакции коммитов

#### Уведомления в MGMTd:

**Типы уведомлений:**
1. **ON-CHANGE** - Уведомления при изменении
2. **PERIODIC** - Периодический опрос (polling)

**Как работает PERIODIC notify:**

```c
// fe_session_periodic_notify_timer() в mgmt_fe_adapter.c
static void fe_session_periodic_notify_timer(struct event *event)
{
    struct mgmt_fe_session_ctx *session = EVENT_ARG(event);
    
    // Получить clients, которые SUBSCRIBED на этот xpath
    clients = mgmt_be_interested_clients(msg->xpath, 
                                         MGMT_BE_XPATH_SUBSCR_TYPE_OPER, 
                                         "GET-DATA");
    
    // Отправить запрос backend'ам
    mgmt_txn_send_notify_selectors(0, session->session_id, clients, 
                                   false, session->periodic_xpaths);
}
```

**Ключевые структуры:**

```c
// mgmt_be_xpath_map - сопоставление xpath и клиентов
struct mgmt_be_xpath_map {
    char *xpath_prefix;
    uint64_t clients;  // bitmask клиентов
};

static struct mgmt_be_xpath_map *be_cfg_xpath_map;    // config xpaths
static struct mgmt_be_xpath_map *be_oper_xpath_map;   // oper xpaths
static struct mgmt_be_xpath_map *be_notif_xpath_map;  // notify xpaths
```

#### Регистрация backend клиентов:

```c
// В bgpd, ripd, isisd и т.д.
static const struct mgmt_be_client_cbs daemon_be_client_data = {
    .config_xpaths = daemon_config_xpaths,
    .nconfig_xpaths = array_size(daemon_config_xpaths),
    .oper_xpaths = daemon_oper_xpaths,
    .noper_xpaths = array_size(daemon_oper_xpaths),
    .rpc_xpaths = daemon_rpc_xpaths,
    .nrpc_xpaths = array_size(daemon_rpc_xpaths),
    .notify_xpaths = daemon_notify_xpaths,  // Для on-change notify
};

// Создание backend client
mgmt_be_client = mgmt_be_client_create("bgpd", &bgpd_be_client_data, 
                                       0, master);
```

### 2. Sysrepo

#### Архитектура:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Sysrepo Integration                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐       unix socket        ┌──────────────┐                 │
│  │  Sysrepo     │── SR_UNIX_SOCK_PATH ───▶│  FRR Daemon  │                 │
│  │   Server     │                          │   (via SB)   │                 │
│  └──────────────┘                          └──────────────┘                 │
│         │                                         │                         │
│         │subscribe                                │ Yang modules            │
│         ▼                                         ▼                         │
│  ┌──────────────┐                        ┌──────────────────┐               │
│  │  libsysrepo  │                        │  northbound_     │               │
│  │  library     │◀──── SB messages ─────▶│   sysrepo.c      │               │
│  └──────────────┘                        └──────────────────┘               │
│                                                                             │
│  Subscribe types:                                                           │
│  - sr_module_change_subscribe() - config changes                           │
│  - sr_oper_get_items_subscribe() - operational data                        │
│  - sr_rpc_subscribe_tree() - RPC calls                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Уведомления через Sysrepo:

```c
// northbound_sysrepo.c
static int frr_sr_notification_send(const char *xpath, struct list *arguments)
{
    sr_val_t *values = NULL;
    size_t values_cnt = 0;
    int ret;
    
    if (arguments && listcount(arguments) > 0) {
        // Конвертировать FRR данные в sysrepo формат
        values_cnt = listcount(arguments);
        ret = sr_new_values(values_cnt, &values);
        // ... заполнение values
    }
    
    // Отправить уведомление
    ret = sr_notif_send(session, xpath, values, values_cnt, 0, 0);
    
    return NB_OK;
}

// Hook registration
hook_register(nb_notification_send, frr_sr_notification_send);
```

### 3. gRPC Northbound

#### Архитектура:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         gRPC Northbound                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Client (Go/Python/GoRyu)                                                   │
│       │                                                                       │
│       │  gRPC Get(CONFIG, "/frr-bgp:bgp")                                  │
│       ▼                                                                       │
│  ┌────────────────────────────────┐                                         │
│  │  frr-northbound.proto          │                                         │
│  │  service Northbound {          │                                         │
│  │    rpc Get(GetRequest)         │                                         │
│  │      returns (stream          │                                         │
│  │               GetResponse)     │                                         │
│  │    rpc EditCandidate(...)      │                                         │
│  │    rpc Commit(...)             │                                         │
│  │  }                             │                                         │
│  └────────────────────────────────┘                                         │
│       │                                                                       │
│       │  NOTIF: No notification support!                                    │
│       ▼                                                                       │
│  ┌────────────────────────────────┐                                         │
│  │  FRR daemon (bgpd/zebra/etc)   │                                         │
│  │                                │                                         │
│  │  running_config (memory)       │                                         │
│  └────────────────────────────────┘                                         │
│                                                                             │
│  Проблема:                                                                  │
│  - Нет pushed уведомлений                                                    │
│  - Только polling черезGet()                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Протокол gRPC (frr-northbound.proto):

```protobuf
service Northbound {
  rpc Get(GetRequest) returns (stream GetResponse) {}
  rpc CreateCandidate(CreateCandidateRequest) returns (CreateCandidateResponse) {}
  rpc EditCandidate(EditCandidateRequest) returns (EditCandidateResponse) {}
  rpc Commit(CommitRequest) returns (CommitResponse) {}
  // Нет notification RPC!
}
```

**Важное ограничение из документации:**
> There is currently no support for YANG notifications

---

## Межмодульное взаимодействие

### Скрытое взаимодействие между демонами

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Daemon Communication Flow                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. vtysh (CLI)                                                             │
│     │                                                                         │
│     │  configure terminal                                                   │
│     │  router bgp 65001                                                     │
│     │  neighbor 10.0.0.1 remote-as 65002                                  │
│     ▼                                                                         │
│  2. vtysh_client_run()                                                      │
│     - Отправляет команду через unix socket в bgpd                          │
│     │                                                                         │
│  3. bgpd процесс                                                            │
│     │                                                                         │
│     │  command_config_read_one_line()                                       │
│     │    ↓                                                                    │
│     │  cmd_execute_command_strict()                                         │
│     │    ↓                                                                    │
│     │  bgp_vty.c:neighbor_remote_as()                                       │
│     │    ↓                                                                    │
│     │  struct peer *peer = peer_create()                                    │
│     │    ↓                                                                    │
│     │  peer->remote_as = as                                                  │
│     │    ↓                                                                    │
│     │  Изменения в памяти: struct bgp *bm->bgp                             │
│                                                                             │
│  4. "show running-config"                                                    │
│     │                                                                         │
│     │  vty_write_config()                                                   │
│     │    ↓                                                                    │
│     │  bgp_config_write()                                                   │
│     │    ↓                                                                    │
│     │  Читает из struct bgp (ПАМЯТЬ, не файл!)                             │
│     ▼                                                                         │
│  5. "write memory"                                                           │
│     │                                                                         │
│     │  file_write_config()                                                  │
│     │    ↓                                                                    │
│     │  vty_write_config()                                                   │
│     │    ↓                                                                    │
│     │  bgp_config_write()                                                   │
│     │    ↓                                                                    │
│     │  Пишет в /etc/frr/bgpd.conf (ФАЙЛ)                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Диаграмма файлов и сокетов

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FRR File Structure                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  /run/frr/                                                                    │
│  ├─ mgmtd_fe.sock      - Frontend socket (для gRPC, vtysh)                 │
│  ├─ mgmtd_be.sock      - Backend socket (для daemon'ов)                    │
│  └─ sysrepo.sock       - Sysrepo socket                                    │
│                                                                             │
│  /etc/frr/                                                                    │
│  ├─ bgpd.conf          - Конфигурация BGP                                  │
│  ├─ zebra.conf         - Конфигурация Zebra                               │
│  ├─ ripd.conf          - Конфигурация RIP                                  │
│  ├─ isisd.conf         - Конфигурация ISIS                                 │
│  ├─ frr.conf           - Главный конфиг (vtysh integrated)                 │
│  └─ daemons            - Список включенных демонов                         │
│                                                                             │
│  /var/run/frr/                                                                │
│  ├─ mgmtd_fe.sock      - Фактическое расположение (relative к runstatedir) │
│  └─ ...                                                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Цепочка вызовов для изменения конфигурации:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│              Цепочка вызовов для "router bgp 65001"                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. vtysh (пользователь вводит команду)                                     │
│     │                                                                         │
│  2. vtysh_execute()                                                         │
│     │                                                                         │
│  3. vtysh_execute_func()                                                    │
│     │                                                                         │
│  4. cmd_execute() -> cmd_execute_command()                                 │
│     │                                                                         │
│  5. cmd_execute_command_real()                                              │
│     │                                                                         │
│  6. vtysh_client_execute(&vtysh_client[bgpd], "router bgp 65001")         │
│     │                                                                         │
│  7. vtysh_client_run()                                                      │
│     │                                                                         │
│  8. connect() -> send("router bgp 65001")                                  │
│     │    Unix socket /run/frr/mgmtd_fe.sock                                │
│     │                                                                         │
│  9. mgmtd receive -> fe_session_handle_edit()                              │
│     │                                                                         │
│ 10. mgmt_be_adapter_send() -> send to bgpd via mgmtd_be.sock              │
│     │                                                                         │
│ 11. bgpd receive -> mgmt_be_client_process_msg()                          │
│     │                                                                         │
│ 12. bgp_vty.c:router_bgp()                                                 │
│     │    Создает/находит struct bgp                                        │
│     │                                                                         │
│ 13. Изменения в памяти: bm->bgp list                                       │
│                                                                             │
│  Файл НЕ обновлен!仅 изменения в памяти!                                    │
│                                                                             │
│  Для fcтизации через write memory:                                          │
│                                                                             │
│  14. vtysh> write memory                                                    │
│     │                                                                         │
│  15. file_write_config() -> vty_write_config()                            │
│     │                                                                         │
│  16. bgp_config_write() читает из struct bgp                               │
│     │                                                                         │
│  17. Печатает в /etc/frr/bgpd.conf                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Проблемы с BGP

### 1. Отсутствие полного Northbound модуля

#### Сравнение с другими демонами:

| Daemon | frr_*_info | frr_*_route_map_info | mgmt_be_client |
|--------|------------|---------------------|----------------|
| **BGP** | ❌ Нет | ✅ Да | ❌ Нет |
| **ISIS** | ✅ Да | ✅ Да | ✅ Да |
| **ZEBRA** | ✅ Да | ✅ Да | ✅ Да |
| **RIP** | ✅ Да | ✅ Да | ✅ Да |

#### Проверка в коде:

```c
// bgpd/bgp_main.c
static const struct frr_yang_module_info *const bgpd_yang_modules[] = {
    &frr_filter_info,
    &frr_interface_info,
    &frr_route_map_info,
    &frr_vrf_info,
    &frr_bgp_route_map_info,  // ТОЛЬКО route-map!
    // frr_bgp_info ОТСУТСТВУЕТ! -这就是问题!
};
```

```c
// isisd/isis_main.c
static const struct frr_yang_module_info *const isisd_yang_modules[] = {
    &frr_filter_info,
    &frr_interface_info,
    &frr_isisd_info,          // ПОЛНЫЙ модуль!
    &frr_route_map_info,
    &frr_vrf_info,
};
```

#### Корень проблемы:

**BGP не регистрирует `frr_bgp_info`**, который должен содержать callbacks для:
- router bgp (instance)
- neighbor configuration
- AS configuration
- global settings
- address-family configuration
- и другие BGP директивы

**Регистрируется только `frr_bgp_route_map_info`**, который покрывает только route-map совпадения.

#### Последствия:

1. **gRPC Get(CONFIG)** не может получить BGP конфигурацию
2. **Sysrepo** не может получать BGP уведомления
3. **MGMTd** не может управлять BGP через northbound API
4. **EditCandidate/Commit** через gRPC не работают для BGP

### 2. Регистрация через vtysh, не через northbound

BGP использует **классическую CLI vtysh** для конфигурации:

```c
// bgpd/bgp_vty.c
install_element(CONFIG_NODE, &router_bgp_cmd);
install_element(BGP_NODE, &bgp_router_id_cmd);
install_element(BGP_NODE, &neighbor remote_as_cmd);
install_element(BGP_NODE, &neighbor description_cmd);
// ... и т.д.
```

Команды:
- `router bgp 65001`
- `neighbor 10.0.0.1 remote-as 65002`
- `network 192.168.1.0/24`

**Проблема:** Эти команды:
1. Сохраняют данные **только в память** (struct bgp)
2. **НЕ вызывают** northbound callbacks
3. **НЕ вызывают** mgmt_be_client notify
4. **НЕ вызывают** sysrepo notifications

### 3. Отсутствие oper_xpaths для BGP

Даже если бы был `frr_bgp_info`, для работы с MGMTd нужен `oper_xpaths`:

```c
// Пример из ripd
static const char *const ripd_oper_xpaths[] = {
    "/frr-backend:clients",
    "/frr-ripd:ripd",  // Для oper data
    "/ietf-key-chain:key-chains",
};
```

**BGP не регистрирует oper_xpaths**, поэтому:
- MGMTd не знает как получить oper data от BGP
- Periodic notify не может работать для BGP
- Operational data запросы не работают для BGP

### 4. Отсутствие RPC для BGP

| Function | Available? | Comment |
|----------|-----------|---------|
| `clear bgp *` | ❌ | CLI only, no northbound RPC |
| `debug bgp updates` | ❌ | CLI only, no northbound RPC |
| `show bgp summary` | ❌ | CLI only, no northbound RPC |

---

## Технические детали реализаций

### northbound_notif.c - Buffering и delivery уведомлений

```c
// lib/northbound_notif.c

// Buffer для изменений
static struct op_changes nb_notif_adds = RB_INITIALIZER(&nb_notif_adds);
static struct op_changes nb_notif_dels = RB_INITIALIZER(&nb_notif_dels);

// Timer для batch delivery (10msec)
void nb_notif_set_walk_timer(void)
{
    event_add_timer_msec(nb_notif_master, timer_walk_start, NULL, 
                         NB_NOTIF_TIMER_MSEC, &nb_notif_timer);
}

// Отправка уведомлений
void mgmt_be_send_ds_replace_notification(const char *path, 
                                          const struct lyd_node *tree,
                                          uint64_t refer_id);

void mgmt_be_send_ds_delete_notification(const char *path);
```

**Как работает:**
1. `nb_notif_add(path)` - добавляет путь в буфер
2. Timer ждет 10msec и собирает все изменения
3. Отправляет batch уведомлений через mgmt_be_client

### mgmt_be_client.c - Backend client library

```c
// lib/mgmt_be_client.c

struct mgmt_be_client_cbs {
    void (*client_connect_notify)(struct mgmt_be_client *client, 
                                  uintptr_t usr_data, bool connected);
    const char *const *config_xpaths;
    uint nconfig_xpaths;
    const char *const *oper_xpaths;
    uint noper_xpaths;
    const char *const *notify_xpaths;
    uint nnotify_xpaths;
    const char *const *rpc_xpaths;
    uint nrpc_xpaths;
};

// Создание client
struct mgmt_be_client *mgmt_be_client_create(
    const char *name, 
    struct mgmt_be_client_cbs *cbs,
    uintptr_t user_data, 
    struct event_loop *event_loop);

// Отправка уведомлений
int mgmt_be_send_ds_replace_notification(const char *path, 
                                         const struct lyd_node *tree,
                                         uint64_t refer_id);

int mgmt_be_send_ds_delete_notification(const char *path);
```

### northbound_sysrepo.c - Sysrepo integration

```c
// lib/northbound_sysrepo.c

// Инициализация
static int frr_sr_init(void)
{
    sr_connect(SR_CONN_DEFAULT, &connection);
    sr_session_start(connection, SR_DS_RUNNING, &session);
    
    // Subscribe на config changes
    sr_module_change_subscribe(session, module->name, NULL, 
                               frr_sr_config_change_cb, NULL, 0,
                               &module->sr_subscription);
    
    // Subscribe на oper data
    sr_oper_get_items_subscribe(session, module->name, xpath,
                                frr_sr_state_cb, NULL, 0,
                                &module->sr_subscription);
}

// Обработка изменений
static int frr_sr_config_change_cb(sr_session_ctx_t *session, uint32_t sub_id,
                                    const char *module_name, const char *xpath,
                                    sr_event_t sr_ev, uint32_t request_id,
                                    void *private_data)
{
    // Получить изменения
    sr_get_changes_iter(session, "//*", &it);
    
    // Применить к candidate config
    frr_sr_process_change(candidate, sr_op, sr_old_val, sr_new_val);
    
    // Commit transaction
    nb_candidate_commit(context, candidate, ...);
}

// Отправка уведомлений
static int frr_sr_notification_send(const char *xpath, struct list *arguments)
{
    // Конвертировать в sr_val_t
    sr_val_t *values;
    sr_new_values(values_cnt, &values);
    
    // Отправить
    sr_notif_send(session, xpath, values, values_cnt, 0, 0);
}
```

---

## Итоговая сводка решений

### Сравнительная таблица механизмов

| Механизм | Протокол | Порты | Notifications | BGP Support | Файл конфиг |
|----------|----------|-------|---------------|-------------|-------------|
| **gRPC** | TCP | 50051-64 | ❌ Нет (polling only) | ❌ Нет | Тот же |
| **Sysrepo** | Unix socket | /var/run/sysrepo.sock | ✅ Да (через nb_notification_send) | ⚠️ Только route-map | Тот же |
| **MGMTd FE** | Unix socket | /run/frr/mgmtd_fe.sock | ✅ ON-CHANGE + PERIODIC | ⚠️ Только с registered xpaths | Тот же |
| **MGMTd BE** | Unix socket | /run/frr/mgmtd_be.sock | - | ✅ Если registered | Тот же |
| **vtysh CLI** | - | - | - | ✅ Полный | Обновляется при write memory |

### Проблемная таблица BGP

| Задача | Доступно? | Причина |
|--------|-----------|---------|
| Получить BGP конфиг через gRPC Get | ❌ | Нет frr_bgp_info в yang_modules |
| Получить BGP конфиг через Sysrepo | ❌ | Нет frr_bgp_info |
| Получить BGP конфиг через MGMTd | ❌ | Нет registered xpaths |
| Уведомления об изменениях BGP | ❌ | Нет registration |
| Управление BGP через gRPC EditCommit | ❌ | Нет northbound |
| Читать BGP конфиг из /etc/frr/bgpd.conf | ✅ | Прямое чтение файла |
| vtysh show running-config | ✅ | Читает из памяти |

### Альтернативные решения для BGP

#### Решение 1: Прямое чтение файла
```go
// Плюсы: Работает, показывает текущее состояние
// Минусы: Нет уведомлений, нужно polling
func readBGPConfig() string {
    data, _ := os.ReadFile("/etc/frr/bgpd.conf")
    return string(data)
}
```

#### Решение 2: vtysh commands
```bash
# Плюсы: Работает, показывает текущее состояние
# Минусы: Нет уведомлений, нужно polling
vtysh -c "show running-config bgpd"
```

#### Решение 3: poll `show running-config`
```go
// Плюсы: Работает, показывает текущее состояние
// Минусы: Нет уведомлений
func pollConfig() {
    out, _ := exec.Command("vtysh", "-c", "show running-config bgpd").Output()
    return string(out)
}
```

#### Решение 4: Добавить frr_bgp_info ( Requires FRR code change )
```c
// Нужно изменить bgpd/bgp_main.c
static const struct frr_yang_module_info *const bgpd_yang_modules[] = {
    &frr_filter_info,
    &frr_interface_info,
    &frr_bgp_info,      // ДОБАВИТЬ - это большой труд
    &frr_route_map_info,
    &frr_vrf_info,
};

// Нужно создать northbound callbacks для всех BGP конфигураций
```

### Рекомендации

**Если нужна полная управляемость BGP:**
1. **Простейший вариант**: Poll `show running-config` или читать файл каждые N секунд
2. **Продвинутый вариант**: Модифицировать FRR добавив `frr_bgp_info`
3. **Компромисс**: Использовать gRPC только для других протоколов (ISIS, OSPF, RIP)

**Для udmin и прочих систем:**
- MGMTd работает хорошо для ISIS/OSPF/RIP/zebra
- BGP требует отдельного подхода из-за отсутствия полного northbound API

### Цепочка файлов FRR

```
/etc/frr/                           - Files
  ├─ bgpd.conf                      - BGP config (после write memory)
  ├─ zebra.conf                     - Zebra config
  ├─ ripd.conf                      - RIP config
  ├─ isisd.conf                     - ISIS config
  ├─ frr.conf                       - Integrated config
  └─ daemons                        - Daemon list

/run/frr/                           - Unix sockets
  ├─ mgmtd_fe.sock                  - Frontend (for gRPC clients)
  ├─ mgmtd_be.sock                  - Backend (for daemons)
  └─ sysrepo.sock                   - Sysrepo socket

/var/run/frr/                       - Runtime files
  ├─ mgmt_*.sock                    - MGMTd sockets (symlinked from /run/frr)
  ├─ frr.log                        - FRR log file
  └─ *.log                          - Daemon logs
```

### Глоссарий терминов

| Термин | Описание |
|--------|----------|
| **Northbound** | API для получения/изменения конфигурации извне |
| **Southbound** | Протоколы для управления устройствами (BGP, OSPF, etc) |
| **MGMTd** | Централизованный management daemon для FRR |
| **Frontend (FE)** | Клиенты MGMTd (gRPC, vtysh) |
| **Backend (BE)** | Daemon'ы (bgpd, zebra, ripd) которые предоставляют данные |
| **Running Config** | Текущая конфигурация в памяти |
| **Startup Config** | Конфигурация в файле `/etc/frr/*.conf` |
| **Operational Data** | Runtime состояние (neighbor state, routes, etc) |
| **Sysrepo** | YANG-based configuration datastore |

---

## Приложение: Проверка в коде

### Проверка наличия frr_bgp_info:

```bash
grep -rn "frr_bgp_info" /Users/shsalibekov/GitHub/FRRouting/frr/bgpd/*.c
# Результат: (пусто - нет такой переменной)

grep -rn "frr_bgp_route_map_info" /Users/shsalibekov/GitHub/FRRouting/frr/bgpd/*.c
# Результат:
# /Users/shsalibekov/GitHub/FRRouting/frr/bgpd/bgp_routemap_nb.c:18:const struct frr_yang_module_info frr_bgp_route_map_info = {
```

### Проверка bgpd_yang_modules:

```bash
grep -A5 "static const struct frr_yang_module_info \*const bgpd_yang_modules" \
    /Users/shsalibekov/GitHub/FRRouting/frr/bgpd/bgp_main.c

# Результат:
# static const struct frr_yang_module_info *const bgpd_yang_modules[] = {
#     &frr_filter_info,
#     &frr_interface_info,
#     &frr_route_map_info,
#     &frr_vrf_info,
#     &frr_bgp_route_map_info,
# };
```

### Проверка ISIS (сравнение):

```bash
grep -A7 "static const struct frr_yang_module_info \*const isisd_yang_modules" \
    /Users/shsalibekov/GitHub/FRRouting/frr/isisd/isis_main.c

# Результат:
# static const struct frr_yang_module_info *const isisd_yang_modules[] = {
#     &frr_filter_info,
#     &frr_interface_info,
#     &frr_isisd_info,              # <-- this exists!
#     &frr_route_map_info,
#     &frr_vrf_info,
# };
```

---

## Дополнение: VTYSH Socket Protocol

### Как работает vtysh protocol

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        VTYSH Protocol Flow                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Client (vtysh)                             Daemon (bgpd)                  │
│  ┌──────────────┐                              ┌──────────────────┐          │
│  │   vtysh      │                              │    bgpd          │          │
│  │              │                              │                  │          │
│  │  - node state│                              │  - node state    │          │
│  │    (local)   │                              │    (local)       │          │
│  │  - buf       │   socket write(line + \0)    │  - buf           │          │
│  │  - hist      │─────────────────────────────▶│  - hist          │          │
│  └──────────────┘                              └──────────────────┘          │
│         │                                              │                    │
│         │   socket read(response + \0)                  │                    │
│        ◀───────────────────────────────────────────────│                    │
│         │                                              │                    │
│  Response handling                             Command parsing               │
│  - display                                     - cmd_execute()               │
│  - callbacks                                   - vty_command()               │
│  - feedback                                    - zlog_notice()               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### VTY SOCKET PROTOCOL

**Формат команды:**
```
<command>\0
```
- Команда завершается **одним нулевым байтом** (`\0`, не `\n`)
- vtysh передает команду как есть, без оберток

**Формат ответа:**
```
<output>\0
```
- Ответ тоже завершается `\0`
- vtysh ждет полный ответ перед отправкой следующей команды

### VTY TYPES

| Type | Описание | daemon side |
|------|----------|-------------|
| **VTY_TERM** | Telnet/stdin/stdout | vty_new_init() |
| **VTY_FILE** | Файл конфигурации | vty_read_file() |
| **VTY_SHELL** | Клиент vtysh | vtysh_main.c |
| **VTY_SHELL_SERV** | Сервер vtysh (daemon) | vtysh_accept() |

### ВАЖНЫЕ ОГРАНИЧЕНИЯ vtysh socket

1. **Одновременный доступ:**
   - vtysh и прямой сокет - **разные соединения**
   - Они НЕ влияют друг на друга
   - Каждое соединение имеет свою vty структуру

2. **Один запрос за раз:**
   ```c
   // vtysh_read(): "counting on vtysh not sending more than 1 command line
   // before waiting on the reply"
   ```
   - Нельзя слать много команд сразу
   - Нужно ждать ответа перед следующей командой

3. **Node state:**
   - vtysh хранит node state **локально** (не в daemon)
   - Если пишешь напрямую в socket - сам управляй node state
   - Используй `end` или `exit` для перехода

---

## Дополнение: Логирование в FRR

### Уровни логирования

| Уровень | ID | Описание |
|---------|----|----------|
| **EMERG** | 0 | Система неработоспособна |
| **ALERT** | 1 | Требуется немедленное действие |
| **CRIT** | 2 | Критические условия |
| **ERR** | 3 | Ошибки |
| **WARNING** | 4 | Предупреждения |
| **NOTICE** | 5 | Информационные сообщения |
| **INFO** | 6 | Информация |
| **DEBUG** | 7 | Отладка |

### Команды логирования

```
# Базовые команды
log file /var/log/frr/frr.log [level]
log stdout [level]
log syslog [level]

# Для логирования команд
log commands
no log commands

# Для отладки BGP
debug bgp neighbor-events
debug bgp updates
```

### Где логируются команды

**`vty_command()`** в `lib/vty.c`:
```c
zlog_notice("%s%s", prompt_str, buf);
```

**Условие логирования:**
- `vty_log_commands == true` (устанавливается командой `log commands`)
- Команда не пустая
- Команда не "echo PING"

**Пример лога:**
```
2026/05/22 15:30:45 vty[123]@localhost: router bgp 65001
2026/05/22 15:30:46 vty[123]@localhost: neighbor 10.0.0.1 remote-as 65002
2026/05/22 15:30:47 vty[123]@localhost: exit
```

### Логи которые BGP пишет

**Постоянные логи (всегда):**
```
%%ADJCHANGE: neighbor 10.0.0.1 in vrf default Down Router ID changed
%%ADJCHANGE: neighbor 10.0.0.1 in vrf default Up
Resetting peer 10.0.0.1 due to change in addpath config
```

**Условные логи (с `log-neighbor-changes`):**
```
%%ADJCHANGE: neighbor 10.0.0.1 in vrf default Down Remote AS changed
%%ADJCHANGE: neighbor 10.0.0.1 in vrf default Down Admin. shutdown
```

### Логирование конфигурационных изменений

**BGP НЕ пишет логи** при изменении параметров:
| Команда | Лог? | Комментарий |
|---------|------|-------------|
| `neighbor X description` | ❌ | Нет логов |
| `neighbor X prefix-list` | ❌ | Нет логов |
| `neighbor X route-map` | ❌ | Нет логов |
| `neighbor X filter-list` | ❌ | Нет логов |
| `router bgp X` | ❌ |-builder сессии |
| `neighbor X remote-as Y` | ⚠️ | Только Down/Up через session reset |

**BGP пишет логи** при:
- Neighbor UP (если log-neighbor-changes)
- Neighbor DOWN (если log-neighbor-changes)
- Router ID change (всегда)
- Update source change (всегда)

### Резюме логирования для BGP

**Цель:** перехватить изменение конфигурации

**Методы:**
| Метод | Работает? | Комментарий |
|-------|-----------|-------------|
| `log-neighbor-changes` | ⚠️ | Только UP/DOWN состояний |
| `log commands` | ✅ | Логирует_ALL_ команды |
| Битое хранение | ❌ | BGP не пишет config change |

**Рекомендация:** Включить `log commands` и читать лог через tail/f

---

## Приложение: Проверка в коде
