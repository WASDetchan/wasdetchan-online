prod_up:
	docker compose -f compose.yaml -f compose.prod.yaml up --build -d
prod_down:
	docker compose -f compose.yaml -f compose.prod.yaml down
