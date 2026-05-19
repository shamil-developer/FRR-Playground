# FRR → gRPC Controller Flow

---

# 1. Проверка подключения

| Статус | vtysh                 | gRPC       | Что делать      |
| ------ | --------------------- | ---------- | --------------- |
| ❌      | `show version`        | ❌          | Проверка FRR    |
| ❌      | `show daemons`        | ❌          | Проверка daemon |
| ⚠️      | `show running-config` | ⚠️ Частично | Получить config |

---

# 2. Capabilities

| Статус | vtysh | gRPC                | Что делать             |
| ------ | ----- | ------------------- | ---------------------- |
| ✅      | ❌     | `GetCapabilities()` | Узнать возможности API |

---

# 3. Interfaces

| Статус | vtysh                 | gRPC    | Что делать        |
| ------ | --------------------- | ------- | ----------------- |
| ✅      | `show interface`      | `Get()` | Читать interfaces |
| ✅      | `show interface json` | `Get()` | Operational state |
| ✅      | `interface ...`       | `EditCandidate()` | Interface config |
| ✅      | `description ...`     | `EditCandidate()` | Description |
| ✅      | `ip address ...`      | `EditCandidate()` | IPv4 address |
| ✅      | `shutdown` / `no shutdown` | `EditCandidate()` | Admin state |
| ✅      | `mtu ...`             | `EditCandidate()` | MTU |
| ✅      | `no ip address ...`   | `EditCandidate()` | Cleanup address |

Важно:
- Для переносимого lab лучше использовать `lo` или заранее созданный test interface.
- Operational counters/status могут отличаться по окружению, поэтому asserts должны проверять стабильные поля.

---

# 4. Routes

| Статус | vtysh                | gRPC              | Что делать           |
| ------ | -------------------- | ----------------- | -------------------- |
| ✅      | `show ip route`      | `Get()`           | Читать routes        |
| ✅      | `show ip route json` | `Get()`           | JSON routes          |
| ✅      | `ip route ...`       | `EditCandidate()` | Создать static route |
| ✅      | `ip route ... name`  | `EditCandidate()` | Route tag/name |
| ✅      | `ip route ... distance` | `EditCandidate()` | Admin distance |
| ✅      | `no ip route ...`    | `EditCandidate()` | Cleanup static route |
| ✅      | `show running-config` | `Get()`          | Verify route config |

Важно:
- После каждого `ip route ...` нужно делать `Get()` running config и `Get()` route state, затем JSON assert.
- Cleanup должен возвращать routing table/config к исходному состоянию.

---

# 5. VRF

| Статус | vtysh      | gRPC              | Что делать  |
| ------ | ---------- | ----------------- | ----------- |
| ✅      | `vrf BLUE` | `EditCandidate()` | Создать VRF |
| ✅      | `show vrf` | `Get()`           | Читать VRF  |
| ✅      | `no vrf BLUE` | `EditCandidate()` | Cleanup VRF |
| ✅      | `show running-config` | `Get()` | Verify VRF config |

Важно:
- Если внутри VRF будут routes/interfaces, cleanup надо делать в правильном порядке: сначала зависимые объекты, потом VRF.

---

# 6. Prefix Lists

| Статус | vtysh                 | gRPC              | Что делать          |
| ------ | --------------------- | ----------------- | ------------------- |
| ✅      | `ip prefix-list ...`  | `EditCandidate()` | Создать prefix-list |
| ✅      | `seq ... permit ...`  | `EditCandidate()` | Permit entry |
| ✅      | `seq ... deny ...`    | `EditCandidate()` | Deny entry |
| ✅      | `ge ... le ...`       | `EditCandidate()` | Prefix range |
| ✅      | `show ip prefix-list` | `Get()`           | Читать prefix-list  |
| ✅      | `no ip prefix-list ... seq ...` | `EditCandidate()` | Cleanup entry |
| ✅      | `no ip prefix-list ...` | `EditCandidate()` | Cleanup list |

Важно:
- Сетевик обычно хочет видеть sequence, permit/deny, ge/le и удаление как отдельные проверки.

---

# 7. Route Maps

| Статус | vtysh            | gRPC              | Что делать        |
| ------ | ---------------- | ----------------- | ----------------- |
| ✅      | `route-map ...`  | `EditCandidate()` | Создать route-map |
| ✅      | `match ...`      | `EditCandidate()` | Match rules       |
| ✅      | `set ...`        | `EditCandidate()` | Policy actions    |
| ✅      | `route-map ... permit ...` | `EditCandidate()` | Permit sequence |
| ✅      | `route-map ... deny ...` | `EditCandidate()` | Deny sequence |
| ✅      | `match ip address prefix-list ...` | `EditCandidate()` | Match prefix-list |
| ✅      | `set metric ...` | `EditCandidate()` | Set metric |
| ✅      | `show route-map` | `Get()`           | Читать route-map  |
| ✅      | `no route-map ...` | `EditCandidate()` | Cleanup route-map |

Важно:
- Route-map workflow должен показывать связку prefix-list -> match -> set и потом cleanup.

---

# 8. BGP Base

| Статус | vtysh               | gRPC              | Что делать  |
| ------ | ------------------- | ----------------- | ----------- |
| ❌      | `router bgp 65001`  | `EditCandidate()` | Создать BGP |
| ❌      | `bgp router-id`     | `EditCandidate()` | Router ID   |
| ❌      | `show bgp summary`  | `Get()`           | BGP summary |
| ❌      | `show bgp neighbor` | `Get()`           | Neighbors   |

Статус: сейчас через pure gRPC не реализуется в FRR 10.6.1.

Доказательства:
- `grpcurl GetCapabilities` на `bgpd:50052` показывает `frr-bgp-route-map`, но не показывает `frr-bgp`.
- В `frr-lib/bgpd/bgp_main.c` в `bgpd_yang_modules[]` подключен `frr_bgp_route_map_info`, но нет `frr_bgp_info`.
- В `frr-lib/bgpd` есть `bgp_routemap_nb.c`, но нет `bgp_nb.c` или northbound callbacks для BGP Base.
- Прямой `EditCandidate()` на путь `frr-bgp:bgp` возвращает `InvalidArgument: Failed to update`.
- `vtysh router bgp 65001` работает, но это CLI path из `bgp_vty.c`, не gRPC northbound.

Важно:
- Через gRPC можно задать общий `zebra` router-id:
  `/frr-vrf:lib/vrf[name='default']/frr-zebra:zebra/router-id`.
- Это не то же самое, что явный `bgp router-id` в BGP config.

---

# 9. BGP Neighbor

| Статус | vtysh                        | gRPC              | Что делать   |
| ------ | ---------------------------- | ----------------- | ------------ |
| ❌      | `neighbor x.x.x.x remote-as` | `EditCandidate()` | Создать peer |
| ❌      | `neighbor description`       | `EditCandidate()` | Description  |
| ❌      | `neighbor timers`            | `EditCandidate()` | Timers       |
| ❌      | `neighbor shutdown`          | `EditCandidate()` | Shutdown     |
| ❌      | `neighbor route-map`         | `EditCandidate()` | Policy       |
| ❌      | `show bgp neighbor`          | `Get()`           | State        |

Статус: сейчас через pure gRPC не реализуется.

Доказательства:
- Neighbor config находится в `frr-bgp.yang`, но модуль `frr-bgp` не объявлен в `bgpd` capabilities.
- В `frr-lib/bgpd` нет northbound callbacks для `neighbor remote-as`, `description`, `timers`, `shutdown`.
- `neighbor route-map` тоже относится к BGP neighbor config; сам `route-map` через gRPC есть, но привязка route-map к BGP neighbor через gRPC недоступна.

---

# 10. Address Family

| Статус | vtysh                         | gRPC              | Что делать   |
| ------ | ----------------------------- | ----------------- | ------------ |
| ❌      | `address-family ipv4 unicast` | `EditCandidate()` | AF           |
| ❌      | `network x.x.x.x/24`          | `EditCandidate()` | Advertise    |
| ❌      | `redistribute static`         | `EditCandidate()` | Redistribute |

Статус: сейчас через pure gRPC не реализуется.

Доказательства:
- Address-family, `network` и BGP redistribute находятся внутри `frr-bgp.yang`.
- `frr-bgp` не подключен к `bgpd` gRPC northbound.
- Static routes через `staticd` работают, но BGP redistribution static через BGP address-family нет.

---

# 11. EVPN / VXLAN ???

| Статус | vtysh                       | gRPC              | Что делать     |
| ------ | --------------------------- | ----------------- | -------------- |
| ❌      | `address-family l2vpn evpn` | `EditCandidate()` | EVPN           |
| ❌      | `advertise-all-vni`         | `EditCandidate()` | EVPN advertise |
| ❌      | `show bgp l2vpn evpn`       | `Get()`           | EVPN state     |
| ⚠️      | `show evpn`                 | `Execute()`       | Zebra EVPN summary |
| ⚠️      | `show evpn vni`             | `Get()`           | VNI            |
| ⚠️      | `show evpn mac vni`         | ⚠️ Частично        | MAC table      |
| ⚠️      | `show evpn arp-cache`       | `Execute()`       | ARP/ND cache   |

Статус: частично.

Что нельзя:
- BGP EVPN config через `address-family l2vpn evpn` и `advertise-all-vni` сейчас недоступен, потому что это часть `frr-bgp`.
- `show bgp l2vpn evpn` через `Get()` недоступен как BGP operational tree.

Что можно проверить отдельно:
- `zebra:50051` объявляет `frr-zebra`.
- В `frr-zebra.yang` есть EVPN/VNI YANG RPC: `get-evpn-info`, `get-vni-info`, `get-evpn-macs`, `get-evpn-arp-cache` и другие.
- Это не заменяет BGP EVPN config, но может быть полезно для чтения zebra EVPN/VNI state через `Execute()`.
- Чтобы это показать красиво, нужен handler `Execute()` и отдельный read-only workflow без попытки создать BGP EVPN.

---

# 12. OSPF

| Статус | vtysh                   | gRPC              | Что делать  |
| ------ | ----------------------- | ----------------- | ----------- |
| ❌      | `router ospf`           | `EditCandidate()` | OSPF        |
| ❌      | `network ... area 0`    | `EditCandidate()` | Add network |
| ❌      | `show ip ospf neighbor` | `Get()`           | Neighbors   |
| ❌      | `show ip ospf database` | `Get()`           | LSDB        |

Статус: сейчас OSPF Base через pure gRPC не реализуется.

Доказательства:
- `grpcurl GetCapabilities` на `ospfd:50055` показывает `frr-ospf-route-map`, но не показывает `frr-ospfd`.
- Прямой `EditCandidate()` на OSPF Base path возвращает `InvalidArgument: Failed to update`.
- Значит `router ospf`, `network ... area 0`, OSPF neighbors и LSDB через `Get()` сейчас недоступны как OSPF northbound tree.

---

# 13. ISIS

| Статус | vtysh                | gRPC              | Что делать |
| ------ | -------------------- | ----------------- | ---------- |
| ✅      | `router isis`        | `EditCandidate()` | ISIS       |
| ✅      | `net ...`            | `EditCandidate()` | NET        |
| ✅      | `ip router isis ...` | `EditCandidate()` | Interface ISIS |
| ⚠️      | `show isis neighbor` | `Get()`           | Neighbor state есть, но для реального соседа нужен второй роутер |
| ❌      | `show isis database` | нет полноценного `Get()` | LSDB не опубликован как gRPC/YANG data tree |

Статус: ISIS Base реализуем через gRPC частично-полноценно: конфигурация ISIS, NET и interface binding работают; полноценный LSDB через `Get()` в текущем FRR northbound недоступен.

Доказательства:
- `grpcurl GetCapabilities` на `isisd:50057` показывает `frr-isisd`.
- В `frr-lib/isisd/isis_nb.c` есть `frr_isisd_info` и callbacks для `/frr-isisd:isis/instance`.
- В `frr-isisd.yang` есть config для `/frr-isisd:isis/instance` и interface augment `/frr-interface:lib/interface/frr-isisd:isis`.
- Smoke-test через gRPC прошел: создан ISIS instance `TEST`, NET `49.0001.0000.0000.0001.00`, interface binding `lo -> ip router isis TEST`, затем все прочитано через `Get()` и удалено.
- В `frr-isisd.yang` нет RPC для `show isis database`; LSDB есть во `vtysh show isis database json`, но не опубликован как полноценное northbound `Get()` дерево.

Важно:
- Для настоящего `show isis neighbor` нужен второй роутер/контейнер, иначе через `Get()` видны только ISIS operational counters на interface.

---

# 14. BFD

| Статус | vtysh                              | gRPC              | Что делать |
| ------ | ---------------------------------- | ----------------- | ---------- |
| ✅      | `bfd`                              | `EditCandidate()` | Enable BFD container |
| ✅      | `profile TEST`                     | `EditCandidate()` | BFD profile |
| ✅      | `detect-multiplier 3`              | `EditCandidate()` | Detection multiplier |
| ✅      | `transmit-interval 300`            | `EditCandidate()` | TX interval |
| ✅      | `receive-interval 300`             | `EditCandidate()` | RX interval |
| ✅      | `echo-mode`                        | `EditCandidate()` | Echo mode |
| ✅      | `echo-interval ...`                | `EditCandidate()` | Echo timers |
| ✅      | `passive-mode`                     | `EditCandidate()` | Passive mode |
| ✅      | `shutdown`                         | `EditCandidate()` | Administrative down |
| ✅      | `log-session-changes`              | `EditCandidate()` | Session logs |
| ✅      | `peer 192.0.2.2 interface lo`      | `EditCandidate()` | Single-hop peer |
| ✅      | `peer ... local-address ...`       | `EditCandidate()` | Local source address |
| ✅      | `peer ... multihop`                | `EditCandidate()` | Multi-hop peer |
| ✅      | `show bfd peers`                   | `Get()`           | Peer state |
| ✅      | `show bfd peers brief`             | `Get()`           | Brief state |
| ✅      | `show bfd peers counters`          | `Get()`           | Counters |
| ✅      | `show bfd peer ... counters`       | `Get()`           | Peer counters |
| ✅      | `no peer ...`                      | `EditCandidate()` | Cleanup peer |
| ✅      | `no profile TEST`                  | `EditCandidate()` | Cleanup profile |
| ✅      | `no bfd`                           | `EditCandidate()` | Cleanup BFD container |

Статус: BFD Base реализуем через gRPC. Надо делать отдельный workflow.

Доказательства:
- `grpcurl GetCapabilities` на `bfdd:50054` показывает `frr-bfdd`.
- В `frr-lib/bfdd/bfdd_nb.c` есть `frr_bfdd_info` и callbacks для `/frr-bfdd:bfdd/bfd`.
- В `frr-bfdd.yang` есть `/frr-bfdd:bfdd/bfd`, `profile`, `sessions/single-hop`, `sessions/multi-hop`, `stats`.
- Smoke-test через gRPC прошел: создан `/frr-bfdd:bfdd/bfd`, прочитан через `Get()` и удален.
- `bfdd` включен с gRPC endpoint `50054`.

Важно:
- Для переносимого lab лучше начать с single-hop peer `192.0.2.2 interface lo` и profile `TEST`.
- Реальный BFD session `Up` без второй BFD стороны не появится; через gRPC проверяем config, state и counters, а состояние ожидаемо будет `down`.
- Широкий operational `Get()` по `/frr-bfdd:bfdd/bfd/sessions` в текущей сборке FRR может уронить `bfdd`; workflow читает state безопасно точечными paths по каждому peer.
- SBFD не смешиваем с BFD Base, он вынесен в отдельный пункт.

---

# 15. SBFD

| Статус | vtysh              | gRPC              | Что делать |
| ------ | ------------------ | ----------------- | ---------- |
| ✅      | `sbfd echo ...`    | `EditCandidate()` | SBFD echo session |
| ✅      | `sbfd init ...`    | `EditCandidate()` | SBFD init session |
| ✅      | `bfd-name ...`     | `EditCandidate()` | Session name |
| ✅      | `remote-discr ...` | `EditCandidate()` | Remote discriminator |
| ✅      | `show bfd peers`   | `Get()`           | SBFD state/stats |

Статус: реализуемо через `frr-bfdd`, но лучше отдельным workflow после BFD Base.

Доказательства:
- В `frr-bfdd.yang` есть `sessions/sbfd-echo` и `sessions/sbfd-init`.
- Для SBFD есть отдельные поля: `bfd-name`, `remote-discr`, `srv6-encap-data`, `srv6-source-ipv6`, `bfd-mode`, `multi-hop`.

Важно:
- Это не обычный BFD peer, поэтому не смешиваем с базовым BFD.
- Реальный Up state потребует второй endpoint/совместимую сторону.

---

# 16. Candidate Workflow

| Статус | vtysh                | gRPC                | Что делать       |
| ------ | -------------------- | ------------------- | ---------------- |
| ✅      | `configure terminal` | `CreateCandidate()` | Create candidate |
| ✅      | `commit check`       | `Commit(VALIDATE)`  | Validate         |
| ✅      | `commit`             | `Commit()`          | Apply            |
| ✅      | `discard`            | `DeleteCandidate()` | Remove candidate |
| ✅      | `update`             | `UpdateCandidate()` | Rebase           |

Статус: реализуемо.

Что уже есть:
- `CreateCandidate()`
- `EditCandidate()`
- `Commit()`

Что надо дописать в контроллере:
- handler для `DeleteCandidate()`.
- handler для `UpdateCandidate()`.
- отдельный сценарий `Commit(VALIDATE)` перед apply, если хотим показать `commit check`.

---

# 17. Transactions

| Статус | vtysh                            | gRPC                                   | Что делать   |
| ------ | -------------------------------- | -------------------------------------- | ------------ |
| ✅      | `show configuration transaction` | `ListTransactions()`                   | Transactions |
| ⚠️      | `rollback configuration`         | `GetTransaction() + LoadToCandidate()` | Rollback     |

Статус: частично реализуемо.

Что уже есть:
- `ListTransactions()`.
- `GetTransaction()`.

Что надо дописать:
- handler для `LoadToCandidate()`.
- workflow rollback: `GetTransaction()` -> `CreateCandidate()` -> `LoadToCandidate()` -> `Commit()`.

---

# 18. Config Load

| Статус | vtysh                        | gRPC                | Что делать  |
| ------ | ---------------------------- | ------------------- | ----------- |
| ✅      | `configuration load`         | `LoadToCandidate()` | Load config |
| ✅      | `show configuration running` | `Get()`             | Read config |

Статус: реализуемо для поддержанных YANG modules.

Что надо дописать:
- handler для `LoadToCandidate()`.

Ограничение:
- Load не сможет загрузить неподдержанный BGP Base config через gRPC, потому что `frr-bgp` отсутствует в capabilities.

---

# 19. Locking

| Статус | vtysh | gRPC             | Что делать |
| ------ | ----- | ---------------- | ---------- |
| ✅      | ❌     | `LockConfig()`   | Lock       |
| ✅      | ❌     | `UnlockConfig()` | Unlock     |

Статус: реализуемо.

Доказательства:
- Проверено через `grpcurl`: `LockConfig()` и `UnlockConfig()` на `bgpd` отвечают успешно.

Что надо дописать:
- handlers `lockConfig` и `unlockConfig`.

---

# 20. YANG RPC

| Статус | vtysh | gRPC        | Что делать  |
| ------ | ----- | ----------- | ----------- |
| ⚠️      | ❌     | `Execute()` | Execute RPC |

Статус: реализуемо только для реально объявленных YANG RPC.

Доказательства:
- В `frr-zebra.yang` есть RPC: `get-route-information`, `get-vrf-info`, `get-evpn-info`, `get-vni-info`, `get-evpn-macs`, `get-evpn-arp-cache` и другие.
- В `frr-bgp.yang` нет RPC для `clear bgp` или `show bgp`.

Что надо дописать:
- handler для `Execute()`.

Ограничение:
- `Execute()` не запускает произвольные `vtysh` команды. Он работает только с YANG RPC, которые объявлены в capabilities/schema.

---

# 21. НЕ ДЕЛАЕТСЯ через gRPC

| Статус | vtysh               | gRPC |
| ------ | ------------------- | ---- |
| ❌      | `clear bgp *`       | ❌    |
| ❌      | `debug bgp updates` | ❌    |
| ❌      | `show logging`      | ❌    |
| ❌      | `terminal monitor`  | ❌    |
| ❌      | `reload`            | ❌    |
| ❌      | `write memory`      | ❌    |
| ❌      | `show tech`         | ❌    |

---

# Что писать в Go сейчас

Правильная цепочка для demo/controller:

```text
GetCapabilities
↓
Get current config
↓
CreateCandidate
↓
EditCandidate
↓
Commit VALIDATE
↓
Commit APPLY
↓
Get operational state
↓
ListTransactions
```

---


| Статус | Пункт                      | Итог                                                                                  |
| ------ | -------------------------- | ------------------------------------------------------------------------------------- |
| ❌      | **8. BGP Base**            | **Нет**                                                                               |
| ❌      | **9. BGP Neighbor**        | **Нет**                                                                               |
| ❌      | **10. BGP Address Family** | **Нет**                                                                               |
| ⚠️      | **11. EVPN/VXLAN**         | **Частично**: BGP EVPN нет, zebra EVPN/VNI state/RPC можно пробовать                  |
| ❌      | **12. OSPF**               | **Нет** для `router ospf/network`; есть только OSPF route-map модули                  |
| ✅      | **13. ISIS**               | **Да**: `frr-isisd` есть, smoke-test через gRPC прошёл                                |
| ✅      | **14. BFD**                | **Да**: `frr-bfdd` есть, smoke-test через gRPC прошёл                                 |
| ✅      | **15. SBFD**               | **Да**: через `frr-bfdd`, но отдельным workflow после BFD Base                        |
| ✅      | **16. Candidate Workflow** | **Да**, но надо дописать handlers для `DeleteCandidate`, `UpdateCandidate`            |
| ⚠️      | **17. Transactions**       | **Да/частично**: `List/Get` есть, rollback через `LoadToCandidate` надо дописать      |
| ✅      | **18. Config Load**        | **Да**, но нужен handler `LoadToCandidate`                                            |
| ✅      | **19. Locking**            | **Да**, нужен handler `LockConfig/UnlockConfig`                                       |
| ⚠️      | **20. YANG RPC**           | **Да**, нужен handler `Execute`, но только для реальных YANG RPC, не для vtysh-команд |
