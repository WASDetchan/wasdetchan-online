prod_up:
	docker compose -f compose.yaml -f compose.prod.yaml up -d

prod_down:
	docker compose -f compose.yaml -f compose.prod.yaml down

local_up:
	go generate
	docker compose up --build

push_image:
	docker pussh wasdetchan-online-app:latest wasdetchan@wasdetchan.ru
