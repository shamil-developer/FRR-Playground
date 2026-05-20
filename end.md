# Доступные

Если `grpcurl` пример короткий, полный набор `update/delete` смотри в `params` соответствующего workflow step.

___

## Base / Capabilities

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
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

| yang module      | yang module path    | daemon  |
| ---------------- | ------------------- | ------- |
| `frr-zebra.yang` | `GetCapabilities()` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto localhost:50051 frr.Northbound/GetCapabilities
```

___

## Diagnostics

___

> [!TIP]
>
> ```vtysh
> show version
> show daemons
> show northbound
> show yang operational-data
> ```

| yang module      | yang module path    | daemon  |
| ---------------- | ------------------- | ------- |
| `frr-zebra.yang` | `GetCapabilities()` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto localhost:50051 frr.Northbound/GetCapabilities
```

> [!WARNING]
>
> Это не полный CLI output один-в-один. Через gRPC реально проверяем FRR version, supported encodings, список YANG modules и факт, что daemon endpoint живой.

___

> [!TIP]
>
> ```vtysh
> show logging
> ```

| yang module        | yang module path      | daemon  |
| ------------------ | --------------------- | ------- |
| `frr-logging.yang` | `/frr-logging:logging` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "CONFIG",
    "encoding": "JSON",
    "path": [
        "/frr-logging:logging"
    ]
}
JSON
```

> [!WARNING]
>
> Через gRPC читается logging configuration. Runtime log buffer/messages как в полном `show logging` через northbound gRPC не отдаются.

___

## Interfaces

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
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
> show interface
> show ip interface brief
> show running-config interface
> ```

| yang module          | yang module path                                   | daemon  |
| -------------------- | -------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib`                               | `zebra` |
| `frr-zebra.yang`     | `/frr-interface:lib/interface/.../frr-zebra:zebra` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "STATE",
    "encoding": "JSON",
    "withDefaults": true,
    "path": [
        "/frr-interface:lib"
    ]
}
JSON
```

> [!WARNING]
>
> `show ip interface brief` закрывается как данные, а не как готовая CLI-таблица: gRPC отдает interface config/state в JSON/YANG tree, а brief-view собирает сам CLI.

___

> [!TIP]
>
> ```vtysh
> configure terminal
> interface lo
>  description TEST
>  ip address 192.0.2.1/32
>  shutdown
>  no shutdown
>  no ip address 192.0.2.1/32
> ```

| yang module          | yang module path                                          | daemon  |
| -------------------- | --------------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']/description`     | `zebra` |
| `frr-zebra.yang`     | `/frr-interface:lib/interface[name='lo']/frr-zebra:zebra` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
        {
            "path": "/frr-interface:lib/interface[name='lo']/description",
            "value": "TEST"
        },
        {
            "path": "/frr-interface:lib/interface[name='lo']/frr-zebra:zebra/ipv4-addrs[ip='192.0.2.1'][prefix-length='32']"
        },
        {
            "path": "/frr-interface:lib/interface[name='lo']/frr-zebra:zebra/enabled",
            "value": "false"
        }
    ]
}
JSON
```

___

## Static Routes

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
"type": "STATE",
"encoding": "JSON",
"path": [
    "/frr-vrf:lib"
]
}
JSON
```

> [!WARNING]
>
> `show nexthop-group`, `show ip nht` и `show nexthop` закрываются только через route/RIB/nexthop данные внутри `/frr-vrf:lib`. Полного отдельного CLI view один-в-один через gRPC нет.

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50053 frr.Northbound/Get <<'JSON'
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

| yang module        | yang module path                               | daemon    |
| ------------------ | ---------------------------------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing`                         | `staticd` |
| `frr-staticd.yang` | `/frr-routing:routing/.../frr-staticd:staticd` | `staticd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
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

grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
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

grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
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

| yang module        | yang module path                               | daemon    |
| ------------------ | ---------------------------------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing`                         | `staticd` |
| `frr-staticd.yang` | `/frr-routing:routing/.../frr-staticd:staticd` | `staticd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
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
> configure terminal
> ip route 10.20.20.0/24 192.168.1.1 10
> ip route 10.30.30.0/24 192.168.1.1 tag 100
> ipv6 route 2001:db8:100::/64 2001:db8::1
> show ip route static
> show ip route 10.10.10.0/24 json
> show ipv6 route
> show ipv6 route json
> ```

| yang module        | yang module path                                          | daemon    |
| ------------------ | --------------------------------------------------------- | --------- |
| `frr-routing.yang` | `/frr-routing:routing`                                    | `staticd` |
| `frr-staticd.yang` | `/frr-routing:routing/.../frr-staticd:staticd/route-list` | `staticd` |
| `frr-vrf.yang`     | `/frr-vrf:lib/.../frr-zebra:zebra/ribs`                   | `zebra`   |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "STATE",
    "encoding": "JSON",
    "withDefaults": true,
    "path": [
        "/frr-vrf:lib/vrf[name='default']/frr-zebra:zebra/ribs"
    ]
}
JSON
```

___

## VRF

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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

## Prefix Lists / ACL

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
> configure terminal
> ipv6 prefix-list TEST6 seq 10 permit 2001:db8::/32 le 64
> ip access-list standard TEST
>  permit 10.0.0.0/8
> access-list 10 permit 10.0.0.0/8
> show ipv6 prefix-list
> show access-list
> show ip prefix-list detail
> ```

| yang module       | yang module path              | daemon  |
| ----------------- | ----------------------------- | ------- |
| `frr-filter.yang` | `/frr-filter:lib`             | `zebra` |
| `frr-filter.yang` | `/frr-filter:lib/prefix-list` | `zebra` |
| `frr-filter.yang` | `/frr-filter:lib/access-list` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "CONFIG",
    "encoding": "JSON",
    "withDefaults": true,
    "path": [
        "/frr-filter:lib"
    ]
}
JSON
```

___

## Route Maps

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/EditCandidate <<'JSON'
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
> configure terminal
> route-map TEST permit 10
>  match interface lo
>  match metric 100
>  match ip address TEST
>  match ipv6 address prefix-list TEST6
>  set ip next-hop 192.0.2.254
>  set tag 100
> show route-map json
> ```

| yang module          | yang module path     | daemon  |
| -------------------- | -------------------- | ------- |
| `frr-route-map.yang` | `/frr-route-map:lib` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "CONFIG",
    "encoding": "JSON",
    "withDefaults": true,
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
>  match community TEST
>  match as-path TEST
>  set local-preference 200
>  set community 65001:100
>  set weight 100
>  set origin igp
> ```

| yang module              | yang module path                             | daemon |
| ------------------------ | -------------------------------------------- | ------ |
| `frr-route-map.yang`     | `/frr-route-map:lib`                         | `bgpd` |
| `frr-bgp-route-map.yang` | `/frr-route-map:lib/.../frr-bgp-route-map:*` | `bgpd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50052 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
        {
            "path": "/frr-route-map:lib/route-map[name='TEST']/entry[sequence='10']/set-action[action='frr-bgp-route-map:set-local-preference']/rmap-set-action/frr-bgp-route-map:local-pref",
            "value": "200"
        }
    ]
}
JSON
```

___

## ISIS

___

> [!TIP]
>
> ```vtysh
> show isis
> ```

| yang module      | yang module path  | daemon  |
| ---------------- | ----------------- | ------- |
| `frr-isisd.yang` | `/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50057 frr.Northbound/Get <<'JSON'
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

| yang module      | yang module path  | daemon  |
| ---------------- | ----------------- | ------- |
| `frr-isisd.yang` | `/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
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

| yang module          | yang module path                                         | daemon  |
| -------------------- | -------------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']`                | `isisd` |
| `frr-isisd.yang`     | `/frr-interface:lib/interface[name='lo']/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
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

| yang module          | yang module path                                               | daemon  |
| -------------------- | -------------------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']/state`                | `isisd` |
| `frr-isisd.yang`     | `/frr-interface:lib/interface[name='lo']/state/frr-isisd:isis` | `isisd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50057 frr.Northbound/Get <<'JSON'
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

| yang module          | yang module path                                           | daemon  |
| -------------------- | ---------------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib/interface[name='lo']/frr-isisd:isis`   | `isisd` |
| `frr-isisd.yang`     | `/frr-isisd:isis/instance[area-tag='TEST'][vrf='default']` | `isisd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50057 frr.Northbound/EditCandidate <<'JSON'
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

## BFD

___

> [!TIP]
>
> ```vtysh
> show bfd peers
> show bfd peers brief
> show bfd peers counters
> ```

| yang module     | yang module path | daemon |
| --------------- | ---------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd` | `bfdd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
"type": "CONFIG",
"encoding": "JSON",
"path": [
    "/frr-bfdd:bfdd"
]
}
JSON
```

> [!WARNING]
>
> BFD peer config/state через gRPC работает. Для counters надежный путь - читать конкретный peer/session path (`.../stats`), а не ожидать готовую общую CLI-таблицу `brief/counters` один-в-один.

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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
> show bfd peer 192.0.2.2 counters
> ```

| yang module     | yang module path                  | daemon |
| --------------- | --------------------------------- | ------ |
| `frr-bfdd.yang` | `/frr-bfdd:bfdd/bfd/sessions/...` | `bfdd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
{
    "type": "STATE",
    "encoding": "JSON",
    "withDefaults": true,
    "path": [
        "/frr-bfdd:bfdd/bfd/sessions/single-hop[peer='192.0.2.2'][interface='lo'][vrf='default'][local-address='']/stats"
    ]
}
JSON
```

___

## SBFD

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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/Get <<'JSON'
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
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50054 frr.Northbound/EditCandidate <<'JSON'
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

___

## Key Chain

___

> [!TIP]
>
> ```vtysh
> key chain TEST
>  key 1
>   key-string SECRET
> ```

| yang module           | yang module path             | daemon         |
| --------------------- | ---------------------------- | -------------- |
| `ietf-key-chain.yang` | `/ietf-key-chain:key-chains` | `ospfd/ospf6d` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50055 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
        {
            "path": "/ietf-key-chain:key-chains/key-chain[name='TEST']/key[key-id='1']/key-string/keystring",
            "value": "SECRET"
        }
    ]
}
JSON
```

___

## Configuration Transactions

___

> [!TIP]
>
> ```vtysh
> commit check
> commit
> discard
> show configuration transaction
> configuration load
> show configuration running
> ```

| yang module    | yang module path                                              | daemon         |
| -------------- | ------------------------------------------------------------- | -------------- |
| Northbound RPC | `Commit/ListTransactions/LoadToCandidate/DeleteCandidate/Get` | любой endpoint |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d '{"candidateId":1,"phase":"VALIDATE"}' localhost:50051 frr.Northbound/Commit
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d '{"candidateId":1}' localhost:50051 frr.Northbound/DeleteCandidate
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d '{}' localhost:50051 frr.Northbound/ListTransactions
```

___

## Zebra / Debug

___

> [!TIP]
>
> ```vtysh
> show zebra
> show zebra dplane
> show running-config zebra
> debug zebra events
> debug zebra kernel
> ```

| yang module      | yang module path          | daemon  |
| ---------------- | ------------------------- | ------- |
| `frr-zebra.yang` | `/frr-zebra:zebra`        | `zebra` |
| `frr-zebra.yang` | `/frr-zebra:zebra/debugs` | `zebra` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50051 frr.Northbound/Get <<'JSON'
{
    "type": "ALL",
    "encoding": "JSON",
    "withDefaults": true,
    "path": [
        "/frr-zebra:zebra"
    ]
}
JSON
```

> [!WARNING]
>
> `show zebra` и `show running-config zebra` покрываются через `/frr-zebra:zebra` config/state. `show zebra dplane` покрывается только по config/state-частям вроде `dplane-queue-limit`; runtime dplane queue/statistics как полный CLI output через gRPC не отдается.

___

## Segment Routing

___

> [!TIP]
>
> ```vtysh
> segment-routing
> show segment-routing srv6
> ```

| yang module      | yang module path                  | daemon  |
| ---------------- | --------------------------------- | ------- |
| `frr-pathd.yang` | `/frr-pathd:pathd/srte`           | `pathd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50061 frr.Northbound/EditCandidate <<'JSON'
{
  "candidateId": 1,
  "update": [
    {
      "path": "/frr-pathd:pathd/srte/segment-list[name='TEST']/protocol-origin",
      "value": "local"
    },
    {
      "path": "/frr-pathd:pathd/srte/segment-list[name='TEST']/originator",
      "value": "config"
    },
    {
      "path": "/frr-pathd:pathd/srte/segment-list[name='TEST']/segment[index='10']/sid-value",
      "value": "16010"
    }
  ]
}
JSON
```

> [!WARNING]
>
> Реально проверен SR-TE model внутри `pathd`: segment-list создается, читается и удаляется через gRPC. Полный SRv6 CLI view `show segment-routing srv6` один-в-один этот workflow не повторяет.

___

## PIM

___

> [!TIP]
>
> ```vtysh
> router pim
>  packets 5
>  join-prune-interval 30
> ```

| yang module    | yang module path | daemon |
| -------------- | ---------------- | ------ |
| `frr-pim.yang` | `/frr-pim:pim`   | `pimd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50058 frr.Northbound/EditCandidate <<'JSON'
{
  "candidateId": 1,
  "update": [
    {
      "path": "/frr-pim:pim/address-family[address-family='frr-routing:ipv4']/packets",
      "value": "5"
    },
    {
      "path": "/frr-pim:pim/address-family[address-family='frr-routing:ipv4']/join-prune-interval",
      "value": "30"
    }
  ]
}
JSON
```

> [!WARNING]
>
> Реально проверена стабильная router-level PIM config через `pimd`. Interface-level PIM в текущей FRR 10.6.1 сборке падает внутри `pimd` при commit, поэтому не используется. Neighbor/mroute operational views как полный CLI output через gRPC пока не отдаются.

___

## RIP

> [!TIP]
>
> ```vtysh
> router rip
>  network 10.0.0.0/8
>  network lo
>  passive-interface lo
>  default-information originate
>  redistribute static metric 2
>  timers basic 5 30 60
> show ip rip
> no router rip
> ```

| yang module     | yang module path  | daemon |
| --------------- | ----------------- | ------ |
| `frr-ripd.yang` | `/frr-ripd:ripd`  | `ripd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50062 frr.Northbound/EditCandidate <<'JSON'
{
  "candidateId": 1,
  "update": [
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/network[.='10.0.0.0/8']"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/interface[.='lo']"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/passive-interface[.='lo']"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/default-information-originate",
      "value": "true"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/redistribute[protocol='static']/metric",
      "value": "2"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/timers/update-interval",
      "value": "5"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/timers/holddown-interval",
      "value": "30"
    },
    {
      "path": "/frr-ripd:ripd/instance[vrf='default']/timers/flush-interval",
      "value": "60"
    }
  ]
}
JSON
```

___

## RIPng

> [!TIP]
>
> ```vtysh
> router ripng
>  network 2001:db8::/32
>  network lo
>  passive-interface lo
>  default-information originate
>  redistribute static metric 2
>  timers basic 5 30 60
> show ipv6 ripng
> no router ripng
> ```

| yang module       | yang module path     | daemon   |
| ----------------- | -------------------- | -------- |
| `frr-ripngd.yang` | `/frr-ripngd:ripngd` | `ripngd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50063 frr.Northbound/EditCandidate <<'JSON'
{
  "candidateId": 1,
  "update": [
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/network[.='2001:db8::/32']"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/interface[.='lo']"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/passive-interface[.='lo']"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/default-information-originate",
      "value": "true"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/redistribute[protocol='static']/metric",
      "value": "2"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/timers/update-interval",
      "value": "5"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/timers/holddown-interval",
      "value": "30"
    },
    {
      "path": "/frr-ripngd:ripngd/instance[vrf='default']/timers/flush-interval",
      "value": "60"
    }
  ]
}
JSON
```

___

## VRRP

> [!TIP]
>
> ```vtysh
> interface lo
>  vrrp 10 priority 110
>  vrrp 10 advertisement-interval 1000
>  vrrp 10 shutdown
>  vrrp 10 ip 192.0.2.254
>  vrrp 10 ipv6 2001:db8::254
> show vrrp
> no vrrp 10
> ```

| yang module          | yang module path                                      | daemon  |
| -------------------- | ----------------------------------------------------- | ------- |
| `frr-interface.yang` | `/frr-interface:lib`                                  | `vrrpd` |
| `frr-vrrpd.yang`    | `/frr-interface:lib/interface/.../frr-vrrpd:vrrp`     | `vrrpd` |

```bash
grpcurl -plaintext -import-path frrpb -proto frr-northbound.proto -d @ localhost:50064 frr.Northbound/EditCandidate <<'JSON'
{
  "candidateId": 1,
  "update": [
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']"
    },
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']/priority",
      "value": "110"
    },
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']/advertisement-interval",
      "value": "1000"
    },
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']/shutdown",
      "value": "true"
    },
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']/v4/virtual-address[.='192.0.2.254']"
    },
    {
      "path": "/frr-interface:lib/interface[name='lo']/frr-vrrpd:vrrp/vrrp-group[virtual-router-id='10']/v6/virtual-address[.='2001:db8::254']"
    }
  ]
}
JSON
```

# Не доступные команды

## BGP

| Доказательства                                                                                                                                                           |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| При запросе `GetCapabilities` не возвращает модуль `frr-bgp` а только `frr-bgp-route-map`                                                                                |
| `grep -R "frr_bgp_info" /src/frr/bgpd` ничего не находит                                                                                                                 |
| `grep -n "frr_bgp_route_map_info" /src/frr/bgpd/bgp_main.c` выодит `123: &frr_bgp_route_map_info,` а `grep -n "frr_bgp_info" /src/frr/bgpd/bgp_main.c` ничего не находит |

> [!CAUTION]
>
> ```vtysh
> router bgp 65001
> ```
>
> ```vtysh
> bgp router-id
> ```
>
> ```vtysh
> show bgp summary
> ```
>
> ```vtysh
> show bgp neighbor
> ```
>
> ```vtysh
> show bgp ipv4 unicast
> ```
>
> ```vtysh
> show bgp ipv4 unicast summary
> ```
>
> ```vtysh
> show bgp ipv4 unicast neighbors
> ```
>
> ```vtysh
> show bgp ipv4 unicast json
> ```
>
> ```vtysh
> neighbor x.x.x.x remote-as
> ```
>
> ```vtysh
> neighbor description
> ```
>
> ```vtysh
> neighbor timers
> ```
>
> ```vtysh
> neighbor shutdown
> ```
>
> ```vtysh
> neighbor route-map
> ```
>
> ```vtysh
> neighbor x.x.x.x activate
> ```
>
> ```vtysh
> neighbor x.x.x.x soft-reconfiguration inbound
> ```
>
> ```vtysh
> neighbor x.x.x.x prefix-list ... in
> ```
>
> ```vtysh
> neighbor x.x.x.x route-map ... in
> ```
>
> ```vtysh
> address-family ipv4 unicast
> ```
>
> ```vtysh
> network x.x.x.x/24
> ```
>
> ```vtysh
> address-family ipv4 unicast
>  redistribute static
> ```
>
> ```vtysh
> aggregate-address ...
> ```
>
> ```vtysh
> default-information originate
> ```
>
> ```vtysh
> address-family l2vpn evpn
> ```
>
> ```vtysh
> advertise-all-vni
> ```
>
> ```vtysh
> show evpn
> ```
>
> ```vtysh
> show evpn vni
> ```
>
> ```vtysh
> show evpn mac vni
> ```
>
> ```vtysh
> show evpn arp-cache
> ```
>
> ```vtysh
> show bgp l2vpn evpn
> ```
>
> ```vtysh
> show running-config bgpd
> ```
___

## OSPF

| Доказательства                                                                                                                                                                   |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| При запросе `GetCapabilities` не возвращает модуль `frr-ospf` а только `frr-ospf-route-map`                                                                                      |
| `grep -R "frr_ospfd_info" /src/frr/ospfd` ничего не находит                                                                                                                      |
| `grep -n "frr_ospf_route_map_info" /src/frr/ospfd/ospf_main.c` выодит `141:	&frr_ospf_route_map_info,` а `grep -n "frr_ospfd_info" /src/frr/ospfd/ospf_main.c` ничего не находит |

> [!CAUTION]
>
> ```vtysh
> show bgp l2vpn evpn json
> ```
>
> ```vtysh
> router ospf
> ```
>
> ```vtysh
> network ... area 0
> ```
>
> ```vtysh
> passive-interface ...
> ```
>
> ```vtysh
> area ... range ...
> ```
>
> ```vtysh
> router ospf
>  redistribute static
> ```
>
> ```vtysh
> show ip ospf
> ```
>
> ```vtysh
> show ip ospf interface
> ```
>
> ```vtysh
> show ip ospf route
> ```
>
> ```vtysh
> show ip ospf neighbor
> ```
>
> ```vtysh
> show ip ospf database
> ```
>
> ```vtysh
> show running-config ospfd
> ```

## ISIS Runtime

| Доказательства                                                                                                                         |
| -------------------------------------------------------------------------------------------------------------------------------------- |
| `isisd GetCapabilities` публикует `frr-isisd`, и ISIS config/state частично работает через workflow `13-isis.yaml`                     |
| Полный LSDB view как `show isis database` через northbound gRPC в текущем FRR не реализован готовым YANG/JSON деревом                  |
| Для настоящей ISIS LSDB в lab нужен второй ISIS router; single-router проверка не даст полноценную базу                                |

> [!CAUTION]
>
> ```vtysh
> show isis database
> ```

## BGP Runtime / Debug Commands

| Доказательства                                                                                                                           |
| ---------------------------------------------------------------------------------------------------------------------------------------- |
| `grep -R "clear bgp" /src/frr/bgpd` находит реализации только в `bgp_vty.c` (`/* one clear bgp command to rule them all */`)             |
| `grep -R "debug bgp" /src/frr/bgpd` находит реализации только в `bgp_debug.c`, `bgp_vty.c`, `bgpd.xref`                                  |
| `grep -R "rpc clear" /src/frr/yang` не содержит BGP RPC, а содержит только `clear-rip-route`, `clear-ripng-route`, `clear-evpn-dup-addr` |
| `grep -R "action clear" /src/frr/yang` ничего не находит                                                                                 |
| При запросе `GetCapabilities` отсутствуют debug/runtime RPC или action для BGP                                                           |


> [!CAUTION]
>
> ```vtysh
> clear bgp *
> ```
>
> ```vtysh
> clear bgp neighbor soft
> ```
>
> ```vtysh
> clear ip bgp *
> ```
>
> ```vtysh
> debug bgp updates
> ```
>
> ```vtysh
> debug bgp neighbor-events
> ```
>
> ```vtysh
> debug bgp bestpath
> ```


## Generic CLI / System Commands

| Доказательства                                                                                                              |
| --------------------------------------------------------------------------------------------------------------------------- |
| `grep -R "rpc.*reload" /src/frr/yang` ничего не находит                                                                     |
| `grep -R "rpc.*write" /src/frr/yang` ничего не находит                                                                      |
| `grep -R "rpc.*tech" /src/frr/yang` ничего не находит                                                                       |
| `grep -R "rpc.*monitor" /src/frr/yang` ничего не находит                                                                    |
| `grep -R "action.*reload" /src/frr/yang` ничего не находит                                                                  |
| `grep -R "action.*write" /src/frr/yang` ничего не находит                                                                   |
| `grep -R "action.*tech" /src/frr/yang` ничего не находит                                                                    |
| `grep -R "action.*monitor" /src/frr/yang` ничего не находит                                                                 |
| `GetCapabilities` не возвращает runtime/system RPC или action для `reload`, `write memory`, `show tech`, `terminal monitor` |

> [!CAUTION]
>
> ```vtysh
> terminal monitor
> ```
>
> ```vtysh
> reload
> ```
>
> ```vtysh
> write memory
> ```
>
> ```vtysh
> show tech
> ```


## Interface Configuration

> [!CAUTION]
>
> ```vtysh
> interface lo
>  mtu 1500
> ```

| Доказательства                                                                                         |
| ------------------------------------------------------------------------------------------------------ |
| `GetCapabilities` возвращает модуль `frr-interface`                                                    |
| `grep -n "mtu" /src/frr/yang/frr-interface.yang` находит `leaf mtu`                                    |
| `Get(type=STATE)` для `/frr-interface:lib` возвращает `state.mtu`                                      |
| `Get(type=CONFIG)` для `/frr-interface:lib` не содержит `config.mtu`                                   |
| `EditCandidate` для `/frr-interface:lib/interface[name='lo']/config/mtu` возвращает `Failed to update` |

## Static Routes

| Доказательства                                                                                                                       |
| ------------------------------------------------------------------------------------------------------------------------------------ |
| `frr-staticd.yang` содержит `route-list`, `path-list`, `tag`, `metric`, `weight`, `bfd`, но не содержит `leaf name` для static route |
| `grep -n "leaf name" frr-lib/yang/frr-staticd.yang` не находит route name                                                            |
| `grpcurl Get` для `/frr-staticd:staticd/.../path-list.../name` возвращает `Data path not found` на `staticd / localhost:50053`       |

> [!CAUTION]
>
> ```vtysh
> ip route 10.40.40.0/24 192.168.1.1 name TEST
> ```


## Configuration Transactions

| Доказательства                                                                                                                    |
| --------------------------------------------------------------------------------------------------------------------------------- |
| `frr-northbound.proto` содержит `rollback_support` только в `GetCapabilitiesResponse`, отдельного `Rollback` RPC нет              |
| `GetCapabilities` в контейнере не отдаёт `rollbackSupport: true`                                                                  |
| `vtysh -c 'rollback configuration 1'` в контейнере возвращает `Unknown command`                                                   |
| В исходниках `rollback configuration` завязан на `HAVE_CONFIG_ROLLBACKS`, то есть это compile-time CLI feature, не gRPC operation |

> [!CAUTION]
>
> ```vtysh
> rollback configuration
> ```


## Filter / Community / Policy Objects

| Доказательства                                                                                                                        |
| ------------------------------------------------------------------------------------------------------------------------------------- |
| `frr-bgp-filter.yang` есть в дереве YANG, но `bgpd GetCapabilities` его не публикует                                                  |
| `bgpd_yang_modules[]` регистрирует `frr_bgp_route_map_info`, но не регистрирует `frr_bgp_filter_info`                                 |
| `grpcurl Get` для `/frr-bgp-filter:lib` на `bgpd / localhost:50052` возвращает `Data path not found`                                  |
| Community/extcommunity/as-path CLI остаются в `bgpd` VTY-коде, но полноценного опубликованного gRPC tree для этих объектов сейчас нет |

> [!CAUTION]
>
> ```vtysh
> bgp community-list standard TEST permit 65001:100
> bgp extcommunity-list standard TEST permit rt 65001:100
> bgp as-path access-list TEST permit .*
> show bgp community-list
> ```


## OSPF6

| Доказательства                                                                                                                   |
| -------------------------------------------------------------------------------------------------------------------------------- |
| `ospf6d GetCapabilities` публикует `frr-ospf-route-map` и `frr-ospf6-route-map`, но не публикует полноценный `frr-ospf6d` module |
| `ospf6d_yang_modules[]` регистрирует `frr_ospf6_route_map_info`, но не регистрирует `frr_ospf6d_info`                            |
| `grpcurl Get` для `/frr-ospf6d:ospf6d` на `ospf6d / localhost:50056` возвращает `Data path not found`                            |

> [!CAUTION]
>
> ```vtysh
> router ospf6
> show ipv6 ospf6 neighbor
> show ipv6 ospf6 database
> ```

## PIM Runtime

| Доказательства                                                                                                      |
| ------------------------------------------------------------------------------------------------------------------- |
| `pimd GetCapabilities` публикует `frr-pim`, и router-level PIM config работает через workflow `21-pim.yaml`         |
| Interface-level PIM config в текущей FRR 10.6.1 сборке падает внутри `pimd` при commit                              |
| Готовый operational view для PIM neighbors и multicast routes через northbound gRPC сейчас не реализован             |

> [!CAUTION]
>
> ```vtysh
> show ip pim neighbor
> show ip mroute
> ```


## Internal Runtime / CPU / Memory Debug

| Доказательства                                                                                                 |
| -------------------------------------------------------------------------------------------------------------- |
| `show memory` реализован как VTY command в `lib_vty.c` / `vtysh.c`, отдельного YANG module/RPC нет             |
| `show event cpu` реализован как VTY command в `event.c` / `vtysh.c`, отдельного YANG module/RPC нет            |
| `vtysh -c 'show thread cpu'` в текущем контейнере возвращает `Unknown command`                                 |
| `GetCapabilities` не публикует runtime/system module для memory/thread/event CPU                               |
| `grep -R "rpc.*memory\\|rpc.*thread\\|rpc.*event" frr-lib/yang` не находит подходящего RPC для этих CLI-команд |

> [!CAUTION]
>
> ```vtysh
> show memory
> show thread cpu
> show event cpu
> ```
