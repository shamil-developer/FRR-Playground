# FRR → gRPC Controller Flow

---

# 1. Проверка подключения

| vtysh                 | gRPC       | Что делать      |
| --------------------- | ---------- | --------------- |
| `show version`        | ❌          | Проверка FRR    |
| `show daemons`        | ❌          | Проверка daemon |
| `show running-config` | ⚠️ Частично | Получить config |

---

# 2. Capabilities

| vtysh | gRPC                | Что делать             |
| ----- | ------------------- | ---------------------- |
| ❌     | `GetCapabilities()` | Узнать возможности API |

---

# 3. Interfaces

| vtysh                 | gRPC    | Что делать        |
| --------------------- | ------- | ----------------- |
| `show interface`      | `Get()` | Читать interfaces |
| `show interface json` | `Get()` | Operational state |

---

# 4. Routes

| vtysh                | gRPC              | Что делать           |
| -------------------- | ----------------- | -------------------- |
| `show ip route`      | `Get()`           | Читать routes        |
| `show ip route json` | `Get()`           | JSON routes          |
| `ip route ...`       | `EditCandidate()` | Создать static route |

---

# 5. VRF

| vtysh      | gRPC              | Что делать  |
| ---------- | ----------------- | ----------- |
| `vrf BLUE` | `EditCandidate()` | Создать VRF |
| `show vrf` | `Get()`           | Читать VRF  |

---

# 6. Prefix Lists

| vtysh                 | gRPC              | Что делать          |
| --------------------- | ----------------- | ------------------- |
| `ip prefix-list ...`  | `EditCandidate()` | Создать prefix-list |
| `show ip prefix-list` | `Get()`           | Читать prefix-list  |

---

# 7. Route Maps

| vtysh            | gRPC              | Что делать        |
| ---------------- | ----------------- | ----------------- |
| `route-map ...`  | `EditCandidate()` | Создать route-map |
| `match ...`      | `EditCandidate()` | Match rules       |
| `set ...`        | `EditCandidate()` | Policy actions    |
| `show route-map` | `Get()`           | Читать route-map  |

---

# 8. BGP Base

| vtysh               | gRPC              | Что делать  |
| ------------------- | ----------------- | ----------- |
| `router bgp 65001`  | `EditCandidate()` | Создать BGP |
| `bgp router-id`     | `EditCandidate()` | Router ID   |
| `show bgp summary`  | `Get()`           | BGP summary |
| `show bgp neighbor` | `Get()`           | Neighbors   |

---

# 9. BGP Neighbor

| vtysh                        | gRPC              | Что делать   |
| ---------------------------- | ----------------- | ------------ |
| `neighbor x.x.x.x remote-as` | `EditCandidate()` | Создать peer |
| `neighbor description`       | `EditCandidate()` | Description  |
| `neighbor timers`            | `EditCandidate()` | Timers       |
| `neighbor shutdown`          | `EditCandidate()` | Shutdown     |
| `neighbor route-map`         | `EditCandidate()` | Policy       |
| `show bgp neighbor`          | `Get()`           | State        |

---

# 10. Address Family

| vtysh                         | gRPC              | Что делать   |
| ----------------------------- | ----------------- | ------------ |
| `address-family ipv4 unicast` | `EditCandidate()` | AF           |
| `network x.x.x.x/24`          | `EditCandidate()` | Advertise    |
| `redistribute static`         | `EditCandidate()` | Redistribute |

---

# 11. EVPN / VXLAN

| vtysh                       | gRPC              | Что делать     |
| --------------------------- | ----------------- | -------------- |
| `address-family l2vpn evpn` | `EditCandidate()` | EVPN           |
| `advertise-all-vni`         | `EditCandidate()` | EVPN advertise |
| `show bgp l2vpn evpn`       | `Get()`           | EVPN state     |
| `show evpn vni`             | `Get()`           | VNI            |
| `show evpn mac vni`         | ⚠️ Частично        | MAC table      |

---

# 12. OSPF

| vtysh                   | gRPC              | Что делать  |
| ----------------------- | ----------------- | ----------- |
| `router ospf`           | `EditCandidate()` | OSPF        |
| `network ... area 0`    | `EditCandidate()` | Add network |
| `show ip ospf neighbor` | `Get()`           | Neighbors   |
| `show ip ospf database` | `Get()`           | LSDB        |

---

# 13. ISIS

| vtysh                | gRPC              | Что делать |
| -------------------- | ----------------- | ---------- |
| `router isis`        | `EditCandidate()` | ISIS       |
| `net ...`            | `EditCandidate()` | NET        |
| `show isis neighbor` | `Get()`           | Neighbor   |
| `show isis database` | `Get()`           | LSDB       |

---

# 14. BFD

| vtysh            | gRPC              | Что делать |
| ---------------- | ----------------- | ---------- |
| `bfd`            | `EditCandidate()` | BFD        |
| `peer ...`       | `EditCandidate()` | Peer       |
| `show bfd peers` | `Get()`           | State      |

---

# 15. Candidate Workflow

| vtysh                | gRPC                | Что делать       |
| -------------------- | ------------------- | ---------------- |
| `configure terminal` | `CreateCandidate()` | Create candidate |
| `commit check`       | `Commit(VALIDATE)`  | Validate         |
| `commit`             | `Commit()`          | Apply            |
| `discard`            | `DeleteCandidate()` | Remove candidate |
| `update`             | `UpdateCandidate()` | Rebase           |

---

# 16. Transactions

| vtysh                            | gRPC                                   | Что делать   |
| -------------------------------- | -------------------------------------- | ------------ |
| `show configuration transaction` | `ListTransactions()`                   | Transactions |
| `rollback configuration`         | `GetTransaction() + LoadToCandidate()` | Rollback     |

---

# 17. Config Load

| vtysh                        | gRPC                | Что делать  |
| ---------------------------- | ------------------- | ----------- |
| `configuration load`         | `LoadToCandidate()` | Load config |
| `show configuration running` | `Get()`             | Read config |

---

# 18. Locking

| vtysh | gRPC             | Что делать |
| ----- | ---------------- | ---------- |
| ❌     | `LockConfig()`   | Lock       |
| ❌     | `UnlockConfig()` | Unlock     |

---

# 19. YANG RPC

| vtysh | gRPC        | Что делать  |
| ----- | ----------- | ----------- |
| ❌     | `Execute()` | Execute RPC |

---

# 20. НЕ ДЕЛАЕТСЯ через gRPC

| vtysh               | gRPC |
| ------------------- | ---- |
| `clear bgp *`       | ❌    |
| `debug bgp updates` | ❌    |
| `show logging`      | ❌    |
| `terminal monitor`  | ❌    |
| `reload`            | ❌    |
| `write memory`      | ❌    |
| `show tech`         | ❌    |

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