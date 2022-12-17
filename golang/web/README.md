# glory web server

## 利用方法
`docker-compose up -d`をして、次に接続用のMYSQLユーザを作成し、適切に権限を付与する。
そして次のように実行すれば良い。

```shell
$ make build

$ vi glory_test.sh
#!/bin/bash

set -eu

export DATA_SOURCE_NAME='{user}:{passwd}@tcp({address}:{port})/{DB_name}'
export GIN_MODE='release'

cd /path/to/this/directory
./bin/glory_server &

$ chmod +x glory_test.sh

$ sudo ./glory_test.sh > /dev/null 2>&1
```
