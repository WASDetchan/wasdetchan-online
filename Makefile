docker_dev:
	docker build -t  wasdetchan-online .
	docker run  -p 8082:80  \
		--rm wasdetchan-online
