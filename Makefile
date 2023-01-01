APP_HOST := plum.lab
APP_DIR := /opt/glory

.PHONY: all
all:
	docker-compose down
	rm -f ./docker/app/config.json
	cp ./config.json ./docker/app/config.json
	docker-compose up -d

.PHONY: clean
clean:
	docker-compose down
	rm -f ./docker/app/config.json

.PHONY: deploy_cmd
deploy_cmd:
	cp ./config.json ./cmd/config.json
	scp -r ./cmd $(APP_HOST):$(APP_DIR)

.PHONY: deploy_app
deploy_app:
	scp -r ./docker/app/glory $(APP_HOST):$(APP_DIR)/docker/app/glory

.PHONY: deploy_full
deploy_full:
	rm -f ./docker/app/config.json ./bin/config.json
	cp ./config.json ./docker/app/config.json
	cp ./config.json ./cmd/config.json
	rm -rf ./docker/db/data .docker/db/log
	scp -r ./docker ./cmd ./.env ./docker-compose.yml $(APP_HOST):$(APP_DIR)
