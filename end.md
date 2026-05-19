# Доступные/Не доступные команды

___

```bash
show running-config
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
| `frr-staticd.yang` | `/frr-routing:staticd` | `staticd` |

```bash
grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
    "candidateId": 1,
    "update": [
            {
            "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']",
            "value": "{\n\"prefix\": \"10.10.10.0/24\",\n\"src-prefix\": \"::/0\",\n\"afi-safi\": \"frr-routing:ipv4-unicast\"\n}\n"
            }
    ],
    "note": "пример укорочен; полный update/delete смотри в params шага"
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
    ],
    "note": "пример укорочен; полный update/delete смотри в params шага"
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
    ],
    "note": "пример укорочен; полный update/delete смотри в params шага"
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
| `frr-staticd.yang` | `/frr-routing:staticd` | `staticd` |

```bash
grpcurl -plaintext -d @ localhost:50053 frr.Northbound/EditCandidate <<'JSON'
{
"candidateId": 1,
"delete": [
    {
    "path": "/frr-routing:routing/control-plane-protocols/control-plane-protocol[type='frr-staticd:staticd'][name='staticd'][vrf='default']/frr-staticd:staticd/route-list[prefix='10.10.10.0/24'][src-prefix='::/0'][afi-safi='frr-routing:ipv4-unicast']"
    }
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
    }
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
    }
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
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
],
"note": "пример укорочен; полный update/delete смотри в params шага"
}
JSON
```

