go-rebuild-frr-proto:
	rm -rf frrgrpc
	
	mkdir -p frrgrpc

	curl -L \
	https://raw.githubusercontent.com/FRRouting/frr/master/grpc/frr-northbound.proto \
	-o frrgrpc/frr.proto

	protoc \
	-I=. \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--go_opt=Mfrrgrpc/frr.proto=github.com/shamil-developer/FRR-Playground/frrgrpc \
	--go-grpc_opt=Mfrrgrpc/frr.proto=github.com/shamil-developer/FRR-Playground/frrgrpc \
	frrgrpc/frr.proto

	rm -f frrgrpc/frr.proto

# ╭━━┳━━╮
# ┃╭╮┃╭╮┃
# ┃╰╯┃╰╯┃
# ╰━╮┣━━╯
# ╭━╯┃
# ╰━━╯

go-deps:
	go mod tidy

go-run:
	go run main.go

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
		-v $(PWD)/configs/frr/frr.conf:/etc/frr/frr.conf \
		-v $(PWD)/configs/frr/daemons:/etc/frr/daemons \
		-v $(PWD)/configs/frr/vtysh.conf:/etc/frr/vtysh.conf \
		frr-playground

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
