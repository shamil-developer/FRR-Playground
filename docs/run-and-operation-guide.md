# Руководство по запуску и эксплуатации

## О Dockerfile

Так как официальный Docker-образ FRR не поддерживает gRPC из коробки, Dockerfile в проекте собирает FRR вручную с включенной поддержкой gRPC.

[Исходники FRR на GitHub](https://github.com/FRRouting/frr)
[Версия 10.6.1](https://github.com/FRRouting/frr/releases/tag/frr-10.6.1)

Официальная документация FRR по сборке:

- [Installing Dependencies](https://docs.frrouting.org/projects/dev-guide/en/latest/building-frr-for-ubuntu2x04.html)


### Сборка Docker-образа

```sh
make docker-build
```

### Запуск контейнера

```sh
make docker-run
```

> [!warning]
>
> Ожидайте вывода что все успешном запущено, перед тем как приступать к эксплуатации

### Запуск сценарив

- `make workflow FILE=configs/playground/workflows/01-check.yaml` - запуск любого YAML-сценария вручную по пути.
- `make workflow-01-check` - базовая проверка подключения к gRPC и чтения CONFIG tree.
- `make workflow-02-capabilities` - capabilities, transactions, lock/unlock и базовые northbound проверки.
- `make workflow-03-interfaces` - чтение interface state/config и изменение interface параметров.
- `make workflow-04-routes` - чтение RIB/static routes, создание static route и cleanup.
- `make workflow-05-vrf` - создание VRF, чтение VRF tree и cleanup.
- `make workflow-06-prefix-lists` - IPv4/IPv6 prefix-list, чтение, проверка и удаление.
- `make workflow-07-route-maps` - route-map, match/set actions, чтение, проверка и cleanup.
- `make workflow-08-bgp-base` - BGP base investigation workflow.
- `make workflow-13-isis` - ISIS config/state: instance, NET, interface, state checks и cleanup.
- `make workflow-14-bfd` - BFD profile, peers, counters/state и cleanup.
- `make workflow-15-sbfd` - SBFD initiator/reflector checks и cleanup.
- `make workflow-16-rip` - RIP config/state, timers, redistribute и cleanup.
- `make workflow-17-ripng` - RIPng config/state, timers, redistribute и cleanup.
- `make workflow-18-vrrp` - VRRP group config/state, IPv4/IPv6 virtual address и cleanup.
- `make workflow-19-diagnostics` - diagnostics через gRPC: version/capabilities/logging/zebra state.
- `make workflow-20-segment-routing` - pathd SR-TE segment-list через gRPC и cleanup.
- `make workflow-21-pim` - PIM router-level config через gRPC и cleanup.

### gRPC Документация

- [Northbound gRPC Documentation](https://docs.frrouting.org/en/latest/grpc.html)
- [Northbound API Architecture](https://docs.frrouting.org/projects/dev-guide/en/latest/northbound/northbound.html)
- [FRR gRPC Proto File](https://github.com/FRRouting/frr/blob/master/grpc/frr-northbound.proto)
- [Northbound gRPC Developer Guide](https://docs.frrouting.org/projects/dev-guide/en/latest/grpc.html)
