# Доступные

Если `grpcurl` пример короткий, полный набор `update/delete` смотри в `params` соответствующего workflow step.

___

> [!TIP]
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

> [!TIP]
>
> ```vtysh
> нет vtysh аналога
> ```

| yang module      | yang module path     | daemon  |
| ---------------- | -------------------- | ------- |
| `frr-zebra.yang` | `GetCapabilities()`  | `zebra` |

```bash
grpcurl -plaintext localhost:50051 frr.Northbound/GetCapabilities
```

___

> [!TIP]
>
> ```vtysh
> show interface lo json
> ```

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

> [!TIP]
>
> ```vtysh
> show ip route json
> ```

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

> [!TIP]
>
> ```vtysh
> show running-config
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> ip route 10.10.10.0/24 192.168.1.1
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> no ip route 10.10.10.0/24 192.168.1.1
> ```

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

> [!TIP]
>
> ```vtysh
> show vrf
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

> [!TIP]
>
> ```vtysh
> vrf BLUE
> ```

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

> [!TIP]
>
> ```vtysh
> no vrf BLUE
> ```

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

> [!TIP]
>
> ```vtysh
> show ip prefix-list json
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> ip prefix-list TEST7 seq 70 permit 70.70.0.0/16 le 32
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> ip prefix-list TEST7 seq 80 deny 80.80.0.0/16 le 32
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> no ip prefix-list TEST7 seq 80 deny 80.80.0.0/16 le 32
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> no ip prefix-list TEST7
> ```

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

> [!TIP]
>
> ```vtysh
> show route-map
> ```

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

> [!TIP]
>
> ```vtysh
> configure terminal
> route-map TEST permit 10
> ```

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

> [!TIP]
>
> ```vtysh
> route-map TEST permit 10
> match ip address prefix-list TEST7
> ```

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

> [!TIP]
>
> ```vtysh
> route-map TEST permit 10
> set metric 100
> ```

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

> [!TIP]
>
> ```vtysh
> route-map TEST deny 20
> ```

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

> [!TIP]
>
> ```vtysh
> no route-map TEST deny 20
> ```

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

> [!TIP]
>
> ```vtysh
> no route-map TEST
> ```

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

> [!TIP]
>
> ```vtysh
> show isis
> ```

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

> [!TIP]
>
> ```vtysh
> router isis TEST
>  net 49.0001.0000.0000.0001.00
> ```

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

> [!TIP]
>
> ```vtysh
> interface lo
>  ip router isis TEST
> ```

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

> [!TIP]
>
> ```vtysh
> show isis neighbor
> ```

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

> [!TIP]
>
> ```vtysh
> no ip router isis TEST
> no router isis TEST
> ```

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

> [!TIP]
>
> ```vtysh
> show bfd peers
> ```

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

> [!TIP]
>
> ```vtysh
> bfd
>  profile TEST
>   detect-multiplier 3
>   transmit-interval 300
>   receive-interval 300
>   echo-mode
>   echo-interval 50
>   passive-mode
>   shutdown
>   log-session-changes
> ```

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

> [!TIP]
>
> ```vtysh
> bfd
>  peer 192.0.2.2 local-address 127.0.0.1 interface lo
>   profile TEST
>   detect-multiplier 3
>   transmit-interval 300
>   receive-interval 300
>   echo-mode
>   passive-mode
>   shutdown
>   log-session-changes
> ```

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

> [!TIP]
>
> ```vtysh
> bfd
>  peer 192.0.2.3 multihop local-address 127.0.0.1
>   profile TEST
>   detect-multiplier 3
>   transmit-interval 300
>   receive-interval 300
>   minimum-ttl 254
>   passive-mode
>   shutdown
>   log-session-changes
> ```

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

> [!TIP]
>
> ```vtysh
> show bfd peer 192.0.2.2 interface lo counters
> ```

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

> [!TIP]
>
> ```vtysh
> show bfd peer 192.0.2.3 multihop local-address 127.0.0.1 counters
> ```

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

> [!TIP]
>
> ```vtysh
> no peer 192.0.2.3 multihop local-address 127.0.0.1
> no peer 192.0.2.2 interface lo
> no profile TEST
> no bfd
> ```

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

> [!TIP]
>
> ```vtysh
> bfd
>  peer 2001:db8::1 bfd-mode sbfd-echo bfd-name TEST-ECHO multihop local-address 2001:db8::1 srv6-source-ipv6 2001:db8::1 srv6-encap-data 2001:db8::2
>   detect-multiplier 3
>   transmit-interval 300
>   receive-interval 300
>   echo-mode
>   echo-interval 50
>   passive-mode
>   shutdown
>   log-session-changes
>  peer 192.0.2.4 bfd-mode sbfd-init bfd-name TEST-INIT multihop local-address 127.0.0.1 remote-discr 1001
>   detect-multiplier 3
>   transmit-interval 300
>   receive-interval 300
>   passive-mode
>   shutdown
>   log-session-changes
> ```

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

> [!TIP]
>
> ```vtysh
> show bfd bfd-name TEST-ECHO
> show bfd bfd-name TEST-ECHO counters
> ```

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

> [!TIP]
>
> ```vtysh
> show bfd bfd-name TEST-INIT
> show bfd bfd-name TEST-INIT counters
> ```

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

> [!TIP]
>
> ```vtysh
> no peer 2001:db8::1 bfd-mode sbfd-echo bfd-name TEST-ECHO multihop local-address 2001:db8::1 srv6-source-ipv6 2001:db8::1 srv6-encap-data 2001:db8::2
> no peer 192.0.2.4 bfd-mode sbfd-init bfd-name TEST-INIT multihop local-address 127.0.0.1 remote-discr 1001
> no bfd
> ```

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

# Частично доступные команды

> [!WARNING]
>
> ```vtysh
> show version
> ```

> [!WARNING]
>
> ```vtysh
> show daemons
> ```

> [!WARNING]
>
> ```vtysh
> show logging
> ```

> [!WARNING]
>
> ```vtysh
> show isis database
> ```

> [!WARNING]
>
> ```vtysh
> show evpn
> ```

> [!WARNING]
>
> ```vtysh
> show evpn vni
> ```

> [!WARNING]
>
> ```vtysh
> show evpn mac vni
> ```

> [!WARNING]
>
> ```vtysh
> show evpn arp-cache
> ```

# Не доступные команды

> [!CAUTION]
>
> ```vtysh
> router bgp 65001
> ```

> [!CAUTION]
>
> ```vtysh
> bgp router-id
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp summary
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp neighbor
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp ipv4 unicast
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp ipv4 unicast summary
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp ipv4 unicast neighbors
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp ipv4 unicast json
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor x.x.x.x remote-as
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor description
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor timers
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor shutdown
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor route-map
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor x.x.x.x activate
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor x.x.x.x soft-reconfiguration inbound
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor x.x.x.x prefix-list ... in
> ```

> [!CAUTION]
>
> ```vtysh
> neighbor x.x.x.x route-map ... in
> ```

> [!CAUTION]
>
> ```vtysh
> address-family ipv4 unicast
> ```

> [!CAUTION]
>
> ```vtysh
> network x.x.x.x/24
> ```

> [!CAUTION]
>
> ```vtysh
> address-family ipv4 unicast
>  redistribute static
> ```

> [!CAUTION]
>
> ```vtysh
> aggregate-address ...
> ```

> [!CAUTION]
>
> ```vtysh
> default-information originate
> ```

> [!CAUTION]
>
> ```vtysh
> address-family l2vpn evpn
> ```

> [!CAUTION]
>
> ```vtysh
> advertise-all-vni
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp l2vpn evpn
> ```

> [!CAUTION]
>
> ```vtysh
> show bgp l2vpn evpn json
> ```

> [!CAUTION]
>
> ```vtysh
> router ospf
> ```

> [!CAUTION]
>
> ```vtysh
> network ... area 0
> ```

> [!CAUTION]
>
> ```vtysh
> passive-interface ...
> ```

> [!CAUTION]
>
> ```vtysh
> area ... range ...
> ```

> [!CAUTION]
>
> ```vtysh
> router ospf
>  redistribute static
> ```

> [!CAUTION]
>
> ```vtysh
> show ip ospf
> ```

> [!CAUTION]
>
> ```vtysh
> show ip ospf interface
> ```

> [!CAUTION]
>
> ```vtysh
> show ip ospf route
> ```

> [!CAUTION]
>
> ```vtysh
> show ip ospf neighbor
> ```

> [!CAUTION]
>
> ```vtysh
> show ip ospf database
> ```

> [!CAUTION]
>
> ```vtysh
> clear bgp *
> ```

> [!CAUTION]
>
> ```vtysh
> clear bgp neighbor soft
> ```

> [!CAUTION]
>
> ```vtysh
> clear ip bgp *
> ```

> [!CAUTION]
>
> ```vtysh
> debug bgp updates
> ```

> [!CAUTION]
>
> ```vtysh
> debug bgp neighbor-events
> ```

> [!CAUTION]
>
> ```vtysh
> debug bgp bestpath
> ```

> [!CAUTION]
>
> ```vtysh
> terminal monitor
> ```

> [!CAUTION]
>
> ```vtysh
> reload
> ```

> [!CAUTION]
>
> ```vtysh
> write memory
> ```

> [!CAUTION]
>
> ```vtysh
> show tech
> ```

# Надо проверить

> [!NOTE]
>
> ```vtysh
> show interface
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  description TEST
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  ip address 192.0.2.1/32
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  mtu 1500
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  shutdown
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  no shutdown
> ```

> [!NOTE]
>
> ```vtysh
> interface lo
>  no ip address 192.0.2.1/32
> ```

> [!NOTE]
>
> ```vtysh
> ip route 10.20.20.0/24 192.168.1.1 10
> ```

> [!NOTE]
>
> ```vtysh
> ip route 10.30.30.0/24 192.168.1.1 tag 100
> ```

> [!NOTE]
>
> ```vtysh
> ip route 10.40.40.0/24 192.168.1.1 name TEST
> ```

> [!NOTE]
>
> ```vtysh
> show ip route static
> ```

> [!NOTE]
>
> ```vtysh
> show ip route 10.10.10.0/24 json
> ```

> [!NOTE]
>
> ```vtysh
> ip prefix-list TEST7 seq 90 permit 90.90.0.0/16 ge 24 le 32
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match interface lo
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match metric 100
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set ip next-hop 192.0.2.254
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set local-preference 200
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set community 65001:100
> ```

> [!NOTE]
>
> ```vtysh
> show bfd peers brief
> ```

> [!NOTE]
>
> ```vtysh
> show bfd peers counters
> ```

> [!NOTE]
>
> ```vtysh
> commit check
> ```

> [!NOTE]
>
> ```vtysh
> commit
> ```

> [!NOTE]
>
> ```vtysh
> discard
> ```

> [!NOTE]
>
> ```vtysh
> show configuration transaction
> ```

> [!NOTE]
>
> ```vtysh
> rollback configuration
> ```

> [!NOTE]
>
> ```vtysh
> configuration load
> ```

> [!NOTE]
>
> ```vtysh
> show configuration running
> ```

> [!NOTE]
>
> ```vtysh
> show ip interface brief
> ```

> [!NOTE]
>
> ```vtysh
> show running-config interface
> ```

> [!NOTE]
>
> ```vtysh
> ipv6 route 2001:db8:100::/64 2001:db8::1
> show ipv6 route
> show ipv6 route json
> ```

> [!NOTE]
>
> ```vtysh
> ipv6 prefix-list TEST6 seq 10 permit 2001:db8::/32 le 64
> show ipv6 prefix-list
> ```

> [!NOTE]
>
> ```vtysh
> ip access-list standard TEST
>  permit 10.0.0.0/8
> show access-list
> ```

> [!NOTE]
>
> ```vtysh
> access-list 10 permit 10.0.0.0/8
> ```

> [!NOTE]
>
> ```vtysh
> bgp community-list standard TEST permit 65001:100
> bgp extcommunity-list standard TEST permit rt 65001:100
> bgp as-path access-list TEST permit .*
> show bgp community-list
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match ip address TEST
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match ipv6 address prefix-list TEST6
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match community TEST
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  match as-path TEST
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set tag 100
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set weight 100
> ```

> [!NOTE]
>
> ```vtysh
> route-map TEST permit 10
>  set origin igp
> ```

> [!NOTE]
>
> ```vtysh
> key chain TEST
>  key 1
>   key-string SECRET
> ```

> [!NOTE]
>
> ```vtysh
> show nexthop-group
> ```

> [!NOTE]
>
> ```vtysh
> show zebra
> show zebra dplane
> ```

> [!NOTE]
>
> ```vtysh
> show northbound
> show yang operational-data
> ```

> [!NOTE]
>
> ```vtysh
> router ospf6
> show ipv6 ospf6 neighbor
> show ipv6 ospf6 database
> ```

> [!NOTE]
>
> ```vtysh
> router rip
> show ip rip
> ```

> [!NOTE]
>
> ```vtysh
> router ripng
> show ipv6 ripng
> ```

> [!NOTE]
>
> ```vtysh
> router pim
> show ip pim neighbor
> show ip mroute
> ```

> [!NOTE]
>
> ```vtysh
> vrrp ...
> show vrrp
> ```

> [!NOTE]
>
> ```vtysh
> segment-routing
> show segment-routing srv6
> ```

> [!NOTE]
>
> ```vtysh
> show ip nht
> show nexthop
> ```

> [!NOTE]
>
> ```vtysh
> show route-map json
> show ip prefix-list detail
> ```

> [!NOTE]
>
> ```vtysh
> show running-config bgpd
> show running-config ospfd
> show running-config zebra
> ```

> [!NOTE]
>
> ```vtysh
> show memory
> show thread cpu
> show event cpu
> ```

> [!NOTE]
>
> ```vtysh
> debug zebra events
> debug zebra kernel
> ```
