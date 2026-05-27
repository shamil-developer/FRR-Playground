# Руководство по запуску и эксплуатации

Этот проект проверяет работу FRR gRPC Northbound через YAML-сценарии из
`configs/playground/workflows`. Перед запуском сценариев нужно подготовить
локальные исходники FRR, собрать Docker-образ и поднять контейнер.

## Что нужно установить локально

- Docker.
- Go версии из `go.mod`.
- `make`.
- `git`.
- `protoc`, `protoc-gen-go` и `protoc-gen-go-grpc` нужны только если меняется
  `frrpb/frr.proto` или нужно пересобрать Go-код из proto.

## Быстрый старт с нуля

```sh
make frr-pull
make frr-checkout BRANCH=frr-10.6.1
make go-deps
make docker-build
make docker-run
```

После `make docker-run` дождитесь строки:

```text
FRR STARTED SUCCESSFULLY
```

Затем во втором терминале запустите базовую проверку:

```sh
make workflow-01-check
```

Если проверка прошла, можно запускать остальные workflow.

## Почему нужен `third_party/frr`

Официальный Docker-образ FRR не включает gRPC из коробки. Поэтому `Dockerfile`
собирает FRR из исходников с флагами `--enable-grpc` и `--enable-mgmtd`.

`Dockerfile` делает:

```dockerfile
COPY third_party/frr/ .
```

Значит перед `make docker-build` каталог `third_party/frr` уже должен
существовать. Его готовят команды:

```sh
make frr-pull
make frr-checkout BRANCH=frr-10.6.1
```

`make frr-pull` клонирует или обновляет FRR:

```sh
git clone https://github.com/FRRouting/frr.git third_party/frr
```

`make frr-checkout BRANCH=...` переключает локальный FRR на нужную ветку или тег.
В этом документе используется версия `frr-10.6.1`, потому что она указана как
рабочая версия для проекта.

## Подготовка Go-кода

Обычная подготовка зависимостей:

```sh
make go-deps
```

Это запускает:

```sh
go mod tidy
```

Если нужно пересобрать Go-код из локального FRR proto, выполните:

```sh
make go-rebuild-frr-local-proto
```

Эта команда:

- очищает `frrpb`;
- копирует `third_party/frr/grpc/frr-northbound.proto` в `frrpb/frr.proto`;
- запускает `protoc`;
- обновляет зависимости через `go mod tidy`.

Перед ней обязательно должны быть выполнены:

```sh
make frr-pull
make frr-checkout BRANCH=frr-10.6.1
```

## Сборка Docker-образа

Обычная сборка:

```sh
make docker-build
```

Сборка без Docker cache:

```sh
make docker-build-no-cache
```

Образ называется:

```text
frr-playground
```

## Запуск контейнера FRR

```sh
make docker-run
```

Контейнер запускается с именем `frr-playground`, в privileged-режиме и
публикует gRPC-порты FRR:

| Daemon | Port  |
| ------ | ----- |
| zebra  | 50051 |
| bgpd   | 50052 |
| staticd | 50053 |
| bfdd   | 50054 |
| ospfd  | 50055 |
| ospf6d | 50056 |
| isisd  | 50057 |
| pimd   | 50058 |
| pbrd   | 50059 |
| fabricd | 50060 |
| pathd  | 50061 |
| ripd   | 50062 |
| ripngd | 50063 |
| vrrpd  | 50064 |

Эти же адреса описаны в `configs/playground/options.yaml`.

Контейнер монтирует локальные файлы:

- `configs/frr/frr.conf` -> `/etc/frr/frr.conf`
- `configs/frr/daemons` -> `/etc/frr/daemons`
- `configs/frr/vtysh.conf` -> `/etc/frr/vtysh.conf`

## Эксплуатационные команды Docker

Зайти в shell контейнера:

```sh
make docker-shell
```

Открыть `vtysh`:

```sh
make docker-vtysh
```

Смотреть лог FRR:

```sh
make docker-logs
```

Остановить контейнер:

```sh
make docker-stop
```

Удалить контейнер и образ:

```sh
make docker-clean
```

## Запуск workflow

Workflow запускаются локальным Go-приложением. Контейнер FRR должен быть уже
запущен, потому что приложение подключается к gRPC-портам `localhost:50051` -
`localhost:50064`.

Запуск workflow по произвольному YAML-файлу:

```sh
make workflow FILE=configs/playground/workflows/01-check.yaml
```

Запуск напрямую через Go:

```sh
go run cmd/playground/main.go configs/playground/workflows/01-check.yaml
```

или через флаг:

```sh
go run cmd/playground/main.go -workflow configs/playground/workflows/01-check.yaml
```

`make workflow-*` дополнительно сохраняет очищенный от ANSI-кодов лог в:

```text
logs/workflows/<timestamp>_<workflow-name>.log
```

## Доступные workflow

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

Запуск базового набора:

```sh
make workflow-all
```

Сейчас `workflow-all` запускает только workflow `01` - `08`. Workflow `13` -
`21` нужно запускать отдельно.

## Watcher уведомлений MGMTd

Для проверки уведомлений через `/run/frr/mgmtd_fe.sock` есть отдельная команда:

```sh
make go-watch
```

Она:

- собирает `go-watcher` под архитектуру контейнера;
- копирует бинарник в `frr-playground:/usr/local/bin/frr-watcher`;
- запускает watcher внутри контейнера.

Перед запуском `make go-watch` контейнер `frr-playground` должен быть поднят.

## Служебные request-команды

В `Makefile` есть цели `request-health`, `request-vtysh`, `request-socket`,
`request-grpc-check` и `request-grpc-create-candidate`. Они обращаются к HTTP API
на `localhost:8081`.

Текущий основной сценарий работы использует прямые gRPC workflow-команды и
Docker-контейнер FRR. HTTP API на `8081` в этом гайде не поднимается, поэтому
эти цели не входят в обязательную последовательность запуска.

## Типовой порядок работы

Первый запуск после клонирования проекта:

```sh
make frr-pull
make frr-checkout BRANCH=frr-10.6.1
make go-deps
make docker-build
make docker-run
```

Проверка во втором терминале:

```sh
make workflow-01-check
```

Повторная работа, если образ уже собран:

```sh
make docker-run
make workflow-01-check
```

Если менялись конфиги FRR:

```sh
make docker-stop
make docker-run
```

Если менялись исходники FRR или версия FRR:

```sh
make frr-pull
make frr-checkout BRANCH=<branch-or-tag>
make docker-build-no-cache
make docker-run
```

Если менялся `frr-northbound.proto`:

```sh
make go-rebuild-frr-local-proto
make docker-build-no-cache
make docker-run
```

## Ссылки

- [Исходники FRR на GitHub](https://github.com/FRRouting/frr)
- [FRR 10.6.1](https://github.com/FRRouting/frr/releases/tag/frr-10.6.1)
- [Installing Dependencies](https://docs.frrouting.org/projects/dev-guide/en/latest/building-frr-for-ubuntu2x04.html)
- [Northbound gRPC Documentation](https://docs.frrouting.org/en/latest/grpc.html)
- [Northbound API Architecture](https://docs.frrouting.org/projects/dev-guide/en/latest/northbound/northbound.html)
- [FRR gRPC Proto File](https://github.com/FRRouting/frr/blob/master/grpc/frr-northbound.proto)
- [Northbound gRPC Developer Guide](https://docs.frrouting.org/projects/dev-guide/en/latest/grpc.html)
