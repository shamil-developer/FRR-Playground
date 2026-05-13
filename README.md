# FRR-Playground

Playground для изучения и тестирования способов взаимодействия с FRR через Go.

## Как собран Dockerfile

Так как официальный Docker-образ FRR не поддерживает gRPC из коробки, Dockerfile в проекте собирает FRR вручную с включенной поддержкой gRPC.

[Исходники FRR на GitHub](https://github.com/FRRouting/frr)
[Версия 10.6.1](https://github.com/FRRouting/frr/releases/tag/frr-10.6.1)

Официальная документация FRR по сборке:

- [Installing Dependencies](https://docs.frrouting.org/projects/dev-guide/en/latest/building-frr-for-ubuntu2x04.html)

## Быстрый старт

### Сборка Docker-образа

```sh
make docker-build
```

### Запуск контейнера

```sh
make docker-run
```

После запуска вывод будет примерно таким:

```sh
shmail@Shamils-MacBook-Pro FRR-Playground % make docker-run

docker run \
	--name frr-playground \
	--rm \
	-it \
	--privileged \
	-p 8080:8080 \
	-p 50051:50051 \
	-v /Users/shmail/GitHub/shamil-developer/FRR-Playground/configs/frr/frr.conf:/etc/frr/frr.conf \
	-v /Users/shmail/GitHub/shamil-developer/FRR-Playground/configs/frr/daemons:/etc/frr/daemons \
	-v /Users/shmail/GitHub/shamil-developer/FRR-Playground/configs/frr/vtysh.conf:/etc/frr/vtysh.conf \
	frr-playground

2026/05/13 08:18:03 ZEBRA: [YDG3W-JND95] FD Limit set: 1048576 is stupidly large...
2026/05/13 08:18:03 STATIC: [YDG3W-JND95] FD Limit set: 1048576 is stupidly large...
2026/05/13 08:18:03 BGP: [YDG3W-JND95] FD Limit set: 1048576 is stupidly large...

[41|zebra] sending configuration
[47|bgpd] sending configuration
[56|staticd] sending configuration

[41|zebra] done
[56|staticd] done
[47|bgpd] done

Waiting for children to finish applying config...

2026/05/13 08:18:09 server started on :8081
```

После появления строки:

```sh
server started on :8081
```

приложение готово к работе.

## Архитектура

Схема работы:

```text
You -> Go Server -> FRR
```

Все запросы проходят через Go-приложение, которое напрямую взаимодействует с FRR.

Поэтому Go-приложение запускается внутри одного контейнера вместе с FRR.

## Способы получения данных из FRR

В проекте реализованы самые популярные и практичные способы взаимодействия с FRR:

1. gRPC
2. vtysh
3. Unix Socket

## gRPC

### Документация

- [Northbound gRPC Documentation](https://docs.frrouting.org/en/latest/grpc.html)
- [Northbound API Architecture](https://docs.frrouting.org/projects/dev-guide/en/latest/northbound/northbound.html)
- [FRR gRPC Proto File](https://github.com/FRRouting/frr/blob/master/grpc/frr-northbound.proto)
- [Northbound gRPC Developer Guide](https://docs.frrouting.org/projects/dev-guide/en/latest/grpc.html)

### Команда

```sh
make request-grpc
```

## vtysh

### Документация

- [VTY Shell Documentation](https://docs.frrouting.org/en/latest/vtysh.html)
- [VTYSH Developer Documentation](https://docs.frrouting.org/projects/dev-guide/en/latest/vtysh.html)
- [vtysh Debian Man Page](https://manpages.debian.org/testing/frr/vtysh.1.en.html)

### Команда

```sh
make request-vtysh
```

## Unix Socket

### Документация

- [VTY Shell UNIX Socket Architecture](https://docs.frrouting.org/en/stable-7.3/overview.html)
- [VTYSH Permissions and Socket Access](https://docs.frrouting.org/en/stable-5.0/vtysh.html)
- [FRR Setup and VTY Socket Information](https://github.com/FRRouting/frr/blob/master/doc/user/setup.rst)

### Команда

```sh
make request-socket
```