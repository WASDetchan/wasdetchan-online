prod_up:
	docker compose -f compose.yaml -f compose.prod.yaml up -d
prod_down:
	docker compose -f compose.yaml -f compose.prod.yaml down
