.PHONY: go postgres dockerupp1 api dockerupp2

go:
	sudo docker exec -it pong_client_container /bin/bash

postgres:
	docker exec -it conergie-postgres /bin/bash

api:
	docker exec -it conergie-api-1 /bin/bash

dockerupp1:
	sudo PLAYER_ID=1 docker compose -p pong_player1 --profile player1 up

dockerupp2:
	sudo PLAYER_ID=2 docker compose -p pong_player2 --profile player2 up