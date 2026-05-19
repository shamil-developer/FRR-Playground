go-rebuild-frr-proto:
	rm -rf frrpb
	
	mkdir -p frrpb

	curl -L \
	https://raw.githubusercontent.com/FRRouting/frr/master/grpc/frr-northbound.proto \
	-o frrpb/frr.proto

	protoc \
	-I=. \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--go_opt=Mfrrpb/frr.proto=github.com/shamil-developer/FRR-Playground/frrpb \
	--go-grpc_opt=Mfrrpb/frr.proto=github.com/shamil-developer/FRR-Playground/frrpb \
	frrpb/frr.proto

	rm -f frrpb/frr.proto

# ╭━━┳━━╮
# ┃╭╮┃╭╮┃
# ┃╰╯┃╰╯┃
# ╰━╮┣━━╯
# ╭━╯┃
# ╰━━╯

go-deps:
	go mod tidy

go-run:
	go run cmd/playground/main.go

workflow:
	go run cmd/playground/main.go $(FILE)

workflow-01-check:
	go run cmd/playground/main.go configs/playground/workflows/01-check.yaml

workflow-02-capabilities:
	go run cmd/playground/main.go configs/playground/workflows/02-capabilities.yaml

workflow-03-interfaces:
	go run cmd/playground/main.go configs/playground/workflows/03-interfaces.yaml

workflow-04-routes:
	go run cmd/playground/main.go configs/playground/workflows/04-routes.yaml

workflow-05-vrf:
	go run cmd/playground/main.go configs/playground/workflows/05-vrf.yaml

workflow-06-prefix-lists:
	go run cmd/playground/main.go configs/playground/workflows/06-prefix-lists.yaml

workflow-07-route-maps:
	go run cmd/playground/main.go configs/playground/workflows/07-route-maps.yaml

workflow-08-bgp-base:
	go run cmd/playground/main.go configs/playground/workflows/08-bgp-base.yaml

workflow-13-isis:
	go run cmd/playground/main.go configs/playground/workflows/13-isis.yaml

workflow-14-bfd:
	go run cmd/playground/main.go configs/playground/workflows/14-bfd.yaml

workflow-15-sbfd:
	go run cmd/playground/main.go configs/playground/workflows/15-sbfd.yaml

workflow-16-rip:
	go run cmd/playground/main.go configs/playground/workflows/16-rip.yaml

workflow-17-ripng:
	go run cmd/playground/main.go configs/playground/workflows/17-ripng.yaml

workflow-18-vrrp:
	go run cmd/playground/main.go configs/playground/workflows/18-vrrp.yaml

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
