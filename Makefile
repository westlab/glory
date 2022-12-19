
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
