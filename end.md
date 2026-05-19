# Доступные/Не доступные команды

Если `grpcurl` пример короткий, полный набор `update/delete` смотри в `params` соответствующего workflow step.

___

> [!IMPORTANT]
> **Команда сетевика**
>
> ```vtysh
> show running-config
> ```

| yang module    | yang module path | daemon  |
| -------------- | ---------------- | ------- |
| `frr-vrf.yang` | `/frr-vrf:lib`   | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "CONFIG",
    "encoding": "JSON",
    "path": [
        "/frr-vrf:lib"
    ]
}
JSON
```

___

```bash
нет vtysh аналога
```

| yang module      | yang module path     | daemon  |
| ---------------- | -------------------- | ------- |
| `frr-zebra.yang` | `GetCapabilities()`  | `zebra` |

```bash
grpcurl -plaintext localhost:50051 frr.Northbound/GetCapabilities
```

___

```bash
show interface lo json
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "STATE",
    "encoding": "JSON",
    "path": [
        "/frr-interface:lib/interface[name='lo']/state"
    ]
}
JSON
```

___

```bash
show ip route json
```

| yang module    | yang module path | daemon  |
| -------------- | ---------------- | ------- |
| `frr-vrf.yang` | `/frr-vrf:lib`   | `zebra` |


```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-vrf:lib"
]
}
JSON
```

___

```bash
show running-config
```

| yang module        | yang module path       | daemon    |
| ------------------ | ---------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing` | `staticd` |

```bash
grpcurl -plaintext -d @ localhost:50053 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-routing:routing/control-plane-protocols"
]
}
JSON
```

___

```bash
configure terminal
ip route 10.10.10.0/24 192.168.1.1
```

| yang module        | yang module path       | daemon    |
| ------------------ | ---------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing` | `staticd` |
| `frr-staticd.yang` | `/frr-routing:routing/.../frr-staticd:staticd` | `staticd` |

```bash
grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
            {
            "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']",
            "value": "{\n\"prefix\": \"10.10.10.0/24\",\n\"src-prefix\": \"::/0\",\n\"afi-safi\": \"frr-routing:ipv4-unicast\"\n}\n"
            }
    ]
}
JSON

grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
        {
            "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']/path-list[table-id='0'][distance='1']",
            "value": "{\n\"table-id\": 0,\n\"distance\": 1,\n\"tag\": 0\n}\n"
        }
    ]
}
JSON

grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
        {
            "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']/path-list[table-id='0'][distance='1']/frr-nexthops/nexthop[nh-type='ip4'][gateway='192.168.1.1'][vrf='default'][interface='']",
            "value": "{\n\"nh-type\": \"ip4\",\n\"gateway\": \"192.168.1.1\",\n\"vrf\": \"default\",\n\"interface\": \"\"\n}\n"
        }
    ]
}
JSON
```

___

```bash
configure terminal
no ip route 10.10.10.0/24 192.168.1.1
```

| yang module        | yang module path       | daemon    |
| ------------------ | ---------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing` | `staticd` |
| `frr-staticd.yang` | `/frr-routing:routing/.../frr-staticd:staticd` | `staticd` |

```bash
grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']"
    }
]
}
JSON
```

___

```bash
show vrf
```

| yang module    | yang module path | daemon  |
| -------------- | ---------------- | ------- |
| `frr-vrf.yang` | `/frr-vrf:lib`   | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-vrf:lib"
]
}
JSON
```

___

```bash
vrf BLUE
```

| yang module    | yang module path | daemon  |
| -------------- | ---------------- | ------- |
| `frr-vrf.yang` | `/frr-vrf:lib`   | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-vrf:lib/vrf[name='BLUE']",
    "value": "{\n\"name\": \"BLUE\"\n}\n"
    }
]
}
JSON
```

___

```bash
no vrf BLUE
```

| yang module    | yang module path | daemon  |
| -------------- | ---------------- | ------- |
| `frr-vrf.yang` | `/frr-vrf:lib`   | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-vrf:lib/vrf[name='BLUE']"
    }
]
}
JSON
```

___

```bash
show ip prefix-list json
```

| yang module       | yang module path  | daemon  |
| ----------------- | ----------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-filter:lib"
]
}
JSON
```

___

```bash
configure terminal
ip prefix-list TEST7 seq 70 permit 70.70.0.0/16 le 32
```

| yang module       | yang module path  | daemon  |
| ----------------- | ----------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']"
    },
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='70']/action",
    "value": "permit"
    },
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='70']/ipv4-prefix",
    "value": "70.70.0.0/16"
    },
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='70']/ipv4-prefix-length-lesser-or-equal",
    "value": "32"
    }
]
}
JSON
```

___

```bash
configure terminal
ip prefix-list TEST7 seq 80 deny 80.80.0.0/16 le 32
```

| yang module       | yang module path  | daemon  |
| ----------------- | ----------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='80']/action",
    "value": "deny"
    },
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='80']/ipv4-prefix",
    "value": "80.80.0.0/16"
    },
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='80']/ipv4-prefix-length-lesser-or-equal",
    "value": "32"
    }
]
}
JSON
```

___

```bash
configure terminal
no ip prefix-list TEST7 seq 80 deny 80.80.0.0/16 le 32
```

| yang module       | yang module path  | daemon  |
| ----------------- | ----------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']/entry[sequence='80']"
    }
]
}
JSON
```

___

```bash
configure terminal
no ip prefix-list TEST7
```

| yang module       | yang module path  | daemon  |
| ----------------- | ----------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-filter:lib/prefix-list[type='ipv4'][name='TEST7']"
    }
]
}
JSON
```

___

```bash
show route-map
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-route-map:lib"
]
}
JSON
```

___

```bash
configure terminal
route-map TEST permit 10
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='10']/action",
    "value": "permit"
    }
]
}
JSON
```

___

```bash
route-map TEST permit 10
match ip address prefix-list TEST7
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='10']/match-condition[condition='frr-route-map:ipv4-prefix-list']/rmap-match-condition/list-name",
    "value": "TEST7"
    }
]
}
JSON
```

___

```bash
route-map TEST permit 10
set metric 100
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='10']/set-action[action='frr-route-map:set-metric']/rmap-set-action/value",
    "value": "100"
    }
]
}
JSON
```

___

```bash
route-map TEST deny 20
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |


```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='20']/action",
    "value": "deny"
    }
]
}
JSON
```

___

```bash
no route-map TEST deny 20
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='20']"
    }
]
}
JSON
```

___

```bash
no route-map TEST
```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-route-map:lib/route-map[name='TEST']"
    }
]
}
JSON
```

___

```bash
show isis
```

| yang module      | yang module path | daemon  |
| ---------------- | ---------------- | ------- |
| `frr-isisd.yang` | `/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-isisd:isis"
]
}
JSON
```

___

```bash
router isis TEST
 net 49.0001.0000.0000.0001.00
```

| yang module      | yang module path | daemon  |
| ---------------- | ---------------- | ------- |
| `frr-isisd.yang` | `/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']"
    },
    {
    "path": "/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']/area-address",
    "value": "49.0001.0000.0000.0001.00"
    }
]
}
JSON
```

___

```bash
interface lo
 ip router isis TEST
```

| yang module           | yang module path                                      | daemon  |
| --------------------- | ----------------------------------------------------- | ------- |
| `frr-interface.yang`  | `/frr-interface:lib/interface[name='lo']`             | `isisd` |
| `frr-isisd.yang`      | `/frr-interface:lib/interface[name='lo']/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-interface:lib/interface[name='lo']/frr-isisd:isis/area-tag",
    "value": "TEST"
    },
    {
    "path": "/frr-interface:lib/interface[name='lo']/frr-isisd:isis/ipv4-routing",
    "value": true
    }
]
}
JSON
```

___

```bash
show isis neighbor
```

| yang module          | yang module path                                            | daemon  |
| -------------------- | ----------------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']/state`             | `isisd` |
| `frr-isisd.yang`     | `/frr-interface:lib/interface[name='lo']/state/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-interface:lib/interface[name='lo']/state/frr-isisd:isis"
]
}
JSON
```

___

```bash
show isis database
```

| yang module      | yang module path                                | daemon  |
| ---------------- | ----------------------------------------------- | ------- |
| `frr-isisd.yang` | `/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']"
]
}
JSON
```

___

```bash
no ip router isis TEST
no router isis TEST
```

| yang module          | yang module path                                      | daemon  |
| -------------------- | ----------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']/frr-isisd:isis` | `isisd` |
| `frr-isisd.yang`     | `/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']` | `isisd` |

```bash
grpcurl -plaintext -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-interface:lib/interface[name='lo']/frr-isisd:isis"
    },
    {
    "path": "/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']"
    }
]
}
JSON
```

___

```bash
show bfd peers
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd"
]
}
JSON
```

___

```bash
bfd
 profile TEST
  detect-multiplier 3
  transmit-interval 300
  receive-interval 300
  echo-mode
  echo-interval 50
  passive-mode
  shutdown
  log-session-changes
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-bfdd:bfdd/bfd"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/profile[name='TEST']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/profile[name='TEST']/detection-multiplier",
    "value": 3
    }
]
}
JSON
```

___

```bash
bfd
 peer 192.0.2.2 local-address 127.0.0.1 interface lo
  profile TEST
  detect-multiplier 3
  transmit-interval 300
  receive-interval 300
  echo-mode
  passive-mode
  shutdown
  log-session-changes
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/single-hop[dest-addr='192.0.2.2'][interface='lo'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/single-hop[dest-addr='192.0.2.2'][interface='lo'][vrf='default']/source-addr",
    "value": "127.0.0.1"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/single-hop[dest-addr='192.0.2.2'][interface='lo'][vrf='default']/profile",
    "value": "TEST"
    }
]
}
JSON
```

___

```bash
bfd
 peer 192.0.2.3 multihop local-address 127.0.0.1
  profile TEST
  detect-multiplier 3
  transmit-interval 300
  receive-interval 300
  minimum-ttl 254
  passive-mode
  shutdown
  log-session-changes
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/multi-hop[source-addr='127.0.0.1'][dest-addr='192.0.2.3'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/multi-hop[source-addr='127.0.0.1'][dest-addr='192.0.2.3'][vrf='default']/profile",
    "value": "TEST"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/multi-hop[source-addr='127.0.0.1'][dest-addr='192.0.2.3'][vrf='default']/minimum-ttl",
    "value": 254
    }
]
}
JSON
```

___

```bash
show bfd peer 192.0.2.2 interface lo counters
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd/bfd/sessions/single-hop[dest-addr='192.0.2.2'][interface='lo'][vrf='default']/stats"
]
}
JSON
```

___

```bash
show bfd peer 192.0.2.3 multihop local-address 127.0.0.1 counters
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd/bfd/sessions/multi-hop[source-addr='127.0.0.1'][dest-addr='192.0.2.3'][vrf='default']/stats"
]
}
JSON
```

___

```bash
no peer 192.0.2.3 multihop local-address 127.0.0.1
no peer 192.0.2.2 interface lo
no profile TEST
no bfd
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/multi-hop[source-addr='127.0.0.1'][dest-addr='192.0.2.3'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/single-hop[dest-addr='192.0.2.2'][interface='lo'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/profile[name='TEST']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd"
    }
]
}
JSON
```

___

```bash
bfd
 peer 2001:db8::1 bfd-mode sbfd-echo bfd-name TEST-ECHO multihop local-address 2001:db8::1 srv6-source-ipv6 2001:db8::1 srv6-encap-data 2001:db8::2
  detect-multiplier 3
  transmit-interval 300
  receive-interval 300
  echo-mode
  echo-interval 50
  passive-mode
  shutdown
  log-session-changes
 peer 192.0.2.4 bfd-mode sbfd-init bfd-name TEST-INIT multihop local-address 127.0.0.1 remote-discr 1001
  detect-multiplier 3
  transmit-interval 300
  receive-interval 300
  passive-mode
  shutdown
  log-session-changes
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"update": [
    {
    "path": "/frr-bfdd:bfdd/bfd"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/sbfd-echo[source-addr='2001:db8::1'][bfd-name='TEST-ECHO'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/sbfd-init[source-addr='127.0.0.1'][dest-addr='192.0.2.4'][bfd-name='TEST-INIT'][vrf='default']"
    }
]
}
JSON
```

___

```bash
show bfd bfd-name TEST-ECHO
show bfd bfd-name TEST-ECHO counters
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd/bfd/sessions/sbfd-echo"
]
}
JSON
```

___

```bash
show bfd bfd-name TEST-INIT
show bfd bfd-name TEST-INIT counters
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd/bfd/sessions/sbfd-init[source-addr='127.0.0.1'][dest-addr='192.0.2.4'][bfd-name='TEST-INIT'][vrf='default']/stats"
]
}
JSON
```

___

```bash
no peer 2001:db8::1 bfd-mode sbfd-echo bfd-name TEST-ECHO multihop local-address 2001:db8::1 srv6-source-ipv6 2001:db8::1 srv6-encap-data 2001:db8::2
no peer 192.0.2.4 bfd-mode sbfd-init bfd-name TEST-INIT multihop local-address 127.0.0.1 remote-discr 1001
no bfd
```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/sbfd-init[source-addr='127.0.0.1'][dest-addr='192.0.2.4'][bfd-name='TEST-INIT'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd/sessions/sbfd-echo[source-addr='2001:db8::1'][bfd-name='TEST-ECHO'][vrf='default']"
    },
    {
    "path": "/frr-bfdd:bfdd/bfd"
    }
]
}
JSON
```
