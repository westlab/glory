# glory web server

Google Drive File Stream likeに使えるファイルシステムをマウントし、
そこから取得する。

使用ツールは[google-drive-ocamlfuse](https://github.com/astrada/google-drive-ocamlfuse) 。
設定方法は適当におググりください。

## 利用方法
`docker-compose up -d`をして、次に接続用のMYSQLユーザを作成し、適切に権限を付与する。
そして次のようにファイルを生成する。

```shell
$ make build

$ vi batch.sh
#!/bin/bash
set -eu

export DATA_SOURCE_NAME='{user}:{passwd}@tcp({address}:{port})/{DB_name}'
/usr/bin/cd /home/kenta/go/src/github.com/westlab/glory/countchars
./bin/countchars >> ./log/log_$(date +\%Y\%m\%d).log 2>&1

$ chmod +x glory_test.sh

$ sudo ./glory_test.sh > /dev/null 2>&1
```

最後にcrontabで1日1回適当な時刻に処理をするように設定する。
