# ╱╱╭╮╱╱╱╱╱╭╮
# ╱╱┃┃╱╱╱╱╱┃┃
# ╭━╯┣━━┳━━┫┃╭┳━━┳━╮
# ┃╭╮┃╭╮┃╭━┫╰╯┫┃━┫╭╯
# ┃╰╯┃╰╯┃╰━┫╭╮┫┃━┫┃
# ╰━━┻━━┻━━┻╯╰┻━━┻╯

docker-build:
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
	docker exec -it frr-playground curl localhost:8081/health

request-vtysh:
	docker exec -it frr-playground curl localhost:8081/frr/vtysh/running-config

request-socket:
	docker exec -it frr-playground curl localhost:8081/frr/socket/running-config

request-grpc:
	docker exec -it frr-playground curl localhost:8081/frr/grpc/check