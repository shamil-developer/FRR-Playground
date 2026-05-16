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
