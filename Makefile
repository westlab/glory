APP_HOST := plum.lab
APP_DIR := /opt/glory

.PHONY: all
all:
	docker-compose down
	rm -f golang/app/config.json
	cp config.json golang/app/config.json
	docker-compose up -d

.PHONY: clean
clean:
	docker-compose down

.PHONY: exec_mysqld
exec_mysqld:
	docker exec -it mysqld /bin/bash

.PHONY: exec_gloryd
exec_gloryd:
	docker exec -it gloryd /bin/bash

.PHONY: deploy_bin
deploy_bin:
	scp -r ./golang/cmd/bin $(APP_HOST):$(APP_DIR)

.PHONY: deploy_full
deploy_full:
	rm -f golang/app/config.json
	cp config.json golang/app/config.json
	scp -r ./docker $(APP_HOST):$(APP_DIR)
	scp -r ./golang/app $(APP_HOST):$(APP_DIR)
	scp -r ./golang/cmd/bin $(APP_HOST):$(APP_DIR)
	scp ./.env $(APP_HOST):$(APP_DIR)
	scp ./config.json $(APP_HOST):$(APP_DIR)
	scp ./docker-compose.yml $(APP_HOST):$(APP_DIR)
