# ╭━━━╮╭━━━╮╭━╮╭━╮╭━╮╭━╮╭━━━╮╭━╮╱╭╮
# ┃╭━╮┃┃╭━╮┃┃┃╰╯┃┃┃┃╰╯┃┃┃╭━╮┃┃┃╰╮┃┃
# ┃┃╱╰╯┃┃╱┃┃┃╭╮╭╮┃┃╭╮╭╮┃┃┃╱┃┃┃╭╮╰╯┃
# ┃┃╱╭╮┃┃╱┃┃┃┃┃┃┃┃┃┃┃┃┃┃┃┃╱┃┃┃┃╰╮┃┃
# ┃╰━╯┃┃╰━╯┃┃┃┃┃┃┃┃┃┃┃┃┃┃╰━╯┃┃┃╱┃┃┃
# ╰━━━╯╰━━━╯╰╯╰╯╰╯╰╯╰╯╰╯╰━━━╯╰╯╱╰━╯

SHELL := /bin/bash

RUN_WORKFLOW = $(MAKE) --no-print-directory welcome; mkdir -p logs/workflows; ts=$$(date +"%Y-%m-%d_%H-%M-%S"); name=$$(basename "$(1)" .yaml); log="logs/workflows/$${ts}_$${name}.log"; echo "log: $$log"; set -o pipefail; go run cmd/playground/main.go "$(1)" 2>&1 | tee >(perl -pe 's/\x1b\[[0-9;]*[A-Za-z]//g' > "$$log")

go-rebuild-frr-local-proto:
	rm -rf frrpb
	mkdir -p frrpb

	@# Проверяем, существует ли локальный файл перед копированием
	@if [ ! -f "third_party/frr/grpc/frr-northbound.proto" ]; then \
		echo "Ошибка: локальный файл frr-northbound.proto не найден!"; \
		echo "Сначала выполните: make frr-sync и make frr-checkout"; \
		exit 1; \
	fi

	@# Копируем файл из локального репозитория вместо curl
	cp third_party/frr/grpc/frr-northbound.proto frrpb/frr.proto

	protoc \
	-I=. \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--go_opt=Mfrrpb/frr.proto=github.com/shamil-developer/FRR-Playground/frrpb \
	--go-grpc_opt=Mfrrpb/frr.proto=github.com/shamil-developer/FRR-Playground/frrpb \
	frrpb/frr.proto

	go mod tidy


welcome:
	@printf "%s\n" "$$(tput setaf 2)"
	@printf "%s\n" "       ▒▒▒   ▒▒▒       ███████╗██╗░░░██╗░█████╗░░██████╗██████╗░███╗░░██╗"
	@printf "%s\n" "    ▒▒▒▒▒▒   ▒▒▒▒▒▒    ██╔════╝██║░░░██║██╔══██╗██╔════╝██╔══██╗████╗░██║"
	@printf "%s\n" "    ▒▒▒▒▒▒   ▒▒▒▒▒▒    █████╗░░╚██╗░██╔╝██║░░██║╚█████╗░██║░░██║██╔██╗██║"
	@printf "%s\n" "    ▒▒▒▒▒▒             ██╔══╝░░░╚████╔╝░██║░░██║░╚═══██╗██║░░██║██║╚████║"
	@printf "%s\n" "    ▒▒▒▒▒▒   ▒▒▒▒▒▒    ███████╗░░╚██╔╝░░╚█████╔╝██████╔╝██████╔╝██║░╚███║"
	@printf "%s\n" "    ▒▒▒▒▒▒   ▒▒▒▒▒▒    ╚══════╝░░░╚═╝░░░░╚════╝░╚═════╝░╚═════╝░╚═╝░░╚══╝"
	@printf "%s\n" "       ▒▒▒   ▒▒▒"
	@printf "%s\n" "$$(tput sgr0)"


# ╭━━━╮╭━━━╮╭━━━╮
# ┃╭━━╯┃╭━╮┃┃╭━╮┃
# ┃╰━━╮┃╰━╯┃┃╰━╯┃
# ┃╭━━╯┃╭╮╭╯┃╭╮╭╯
# ┃┃╱╱╱┃┃┃╰╮┃┃┃╰╮
# ╰╯╱╱╱╰╯╰━╯╰╯╰━╯

frr-pull: welcome
	@mkdir -p third_party
	@if [ ! -d "third_party/frr" ]; then \
		echo "Клонирование репозитория FRR..."; \
		git clone https://github.com/FRRouting/frr.git third_party/frr; \
	else \
		echo "Обновление репозитория FRR (git pull)..."; \
		cd third_party/frr && git pull; \
	fi

frr-checkout: welcome
	@if [ -z "$(BRANCH)" ]; then \
		echo "Ошибка: укажите ветку. Пример: make checkout BRANCH=master"; \
		exit 1; \
	fi
	@if [ ! -d "third_party/frr" ]; then \
		echo "Ошибка: репозиторий не найден. Сначала выполните 'make sync'"; \
		exit 1; \
	fi
	@echo "Переключение на ветку $(BRANCH)..."
	@cd third_party/frr && git checkout $(BRANCH)


# ╭━━┳━━╮
# ┃╭╮┃╭╮┃
# ┃╰╯┃╰╯┃
# ╰━╮┣━━╯
# ╭━╯┃
# ╰━━╯

go-deps: welcome
	go mod tidy

go-run: welcome
	go run cmd/playground/main.go

workflow:
	@$(call RUN_WORKFLOW,$(FILE))

workflow-01-check:
	@$(call RUN_WORKFLOW,configs/playground/workflows/01-check.yaml)

workflow-02-capabilities:
	@$(call RUN_WORKFLOW,configs/playground/workflows/02-capabilities.yaml)

workflow-03-interfaces:
	@$(call RUN_WORKFLOW,configs/playground/workflows/03-interfaces.yaml)

workflow-04-routes:
	@$(call RUN_WORKFLOW,configs/playground/workflows/04-routes.yaml)

workflow-05-vrf:
	@$(call RUN_WORKFLOW,configs/playground/workflows/05-vrf.yaml)

workflow-06-prefix-lists:
	@$(call RUN_WORKFLOW,configs/playground/workflows/06-prefix-lists.yaml)

workflow-07-route-maps:
	@$(call RUN_WORKFLOW,configs/playground/workflows/07-route-maps.yaml)

workflow-08-bgp-base:
	@$(call RUN_WORKFLOW,configs/playground/workflows/08-bgp-base.yaml)

workflow-13-isis:
	@$(call RUN_WORKFLOW,configs/playground/workflows/13-isis.yaml)

workflow-14-bfd:
	@$(call RUN_WORKFLOW,configs/playground/workflows/14-bfd.yaml)

workflow-15-sbfd:
	@$(call RUN_WORKFLOW,configs/playground/workflows/15-sbfd.yaml)

workflow-16-rip:
	@$(call RUN_WORKFLOW,configs/playground/workflows/16-rip.yaml)

workflow-17-ripng:
	@$(call RUN_WORKFLOW,configs/playground/workflows/17-ripng.yaml)

workflow-18-vrrp:
	@$(call RUN_WORKFLOW,configs/playground/workflows/18-vrrp.yaml)

workflow-19-diagnostics:
	@$(call RUN_WORKFLOW,configs/playground/workflows/19-diagnostics.yaml)

workflow-20-segment-routing:
	@$(call RUN_WORKFLOW,configs/playground/workflows/20-segment-routing.yaml)

workflow-21-pim:
	@$(call RUN_WORKFLOW,configs/playground/workflows/21-pim.yaml)

workflow-all:
	$(MAKE) workflow-01-check
	$(MAKE) workflow-02-capabilities
	$(MAKE) workflow-03-interfaces
	$(MAKE) workflow-04-routes
	$(MAKE) workflow-05-vrf
	$(MAKE) workflow-06-prefix-lists
	$(MAKE) workflow-07-route-maps
	$(MAKE) workflow-08-bgp-base

# ╱╱╭╮╱╱╱╱╱╭╮
# ╱╱┃┃╱╱╱╱╱┃┃
# ╭━╯┣━━┳━━┫┃╭┳━━┳━╮
# ┃╭╮┃╭╮┃╭━┫╰╯┫┃━┫╭╯
# ┃╰╯┃╰╯┃╰━┫╭╮┫┃━┫┃
# ╰━━┻━━┻━━┻╯╰┻━━┻╯

docker-build:
	docker build -t frr-playground .

docker-build-no-cache:
	docker build --no-cache -t frr-playground .

docker-run:
	docker run \
		--name frr-playground \
		--rm \
		-it \
		--privileged \
		-p 8080:8080 \
		-p 50051:50051 \
		-p 50052:50052 \
		-p 50053:50053 \
		-p 50054:50054 \
		-p 50055:50055 \
		-p 50056:50056 \
		-p 50057:50057 \
		-p 50058:50058 \
		-p 50059:50059 \
		-p 50060:50060 \
		-p 50061:50061 \
		-p 50062:50062 \
		-p 50063:50063 \
		-p 50064:50064 \
		-v $(PWD)/configs/frr/frr.conf:/etc/frr/frr.conf \
		-v $(PWD)/configs/frr/daemons:/etc/frr/daemons \
		-v $(PWD)/configs/frr/vtysh.conf:/etc/frr/vtysh.conf \
		frr-playground

docker-shell:
	docker exec -it frr-playground sh

docker-vtysh:
	docker exec -it frr-playground vtysh

docker-logs:
	docker exec -it frr-playground tail -f /var/log/frr/frr.log

docker-stop:
	docker stop frr-playground

docker-clean:
	docker rm -f frr-playground
	docker rmi -f frr-playground

# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╮
# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╯╰╮
# ╭━┳━━┳━━┳╮╭┳━━┳━┻╮╭╯
# ┃╭┫┃━┫╭╮┃┃┃┃┃━┫━━┫┃
# ┃┃┃┃━┫╰╯┃╰╯┃┃━╋━━┃╰╮
# ╰╯╰━━┻━╮┣━━┻━━┻━━┻━╯
# ╱╱╱╱╱╱╱┃┃
# ╱╱╱╱╱╱╱╰╯

request-health:
	curl http://localhost:8081/health

request-vtysh:
	curl http://localhost:8081/frr/vtysh/running-config

request-socket:
	curl http://localhost:8081/frr/socket/running-config

request-grpc-check:
	curl http://localhost:8081/frr/grpc/check

request-grpc-create-candidate:
	curl http://localhost:8081/frr/grpc/create-candidate

# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╮
# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃
# ╭━━┳━━┳━━┳━━┳━━┳━━┳━━┳━━┳━━╋━╋━━╋━━┳┓
# ┃  ╱╱╱╱╱╱╱╱ ╱╱╱╱╱╱╱╱  ╱╱╱╱╱╱╱╱┃┃
# ╰━━┻━━┻━━┻━━┻━━┻━━┻━━┻━━┻━━┻━┻━━┻━━┻┛

go-watch:
	@echo "=== FRR Notification Watcher ==="
	@echo "Cleaning up old watcher binary..."
	@docker exec frr-playground sh -c 'rm -f /usr/local/bin/frr-watcher 2>/dev/null || true'
	@echo "Copying main.go to container..."
	@docker cp go-watcher/main.go frr-playground:/tmp/go-watcher/main.go
	@echo "Building watcher in container..."
	@docker exec frr-playground sh -c 'mkdir -p /tmp/go-watcher && cd /tmp/go-watcher && go mod init frr-watcher 2>/dev/null || true && go build -o /usr/local/bin/frr-watcher .'
	@echo "=== Starting watcher... Press Ctrl+C to stop ==="
	@docker exec frr-playground /usr/local/bin/frr-watcher