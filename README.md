# glory

Google Driveにある学位論文の原稿(docx)の進捗状況を表示する。

## バージョン
go1.15.5

依存パッケージは`go.mod`を参照

## ディレクトリ構成
```
glory/          # 設定や、定数などの情報
  ├ countChars/ # Docxの文字数をカウントする
  ├ docker/     # dockerの設定
  ├ fetchFile/ (old) # docxファイルを取得し、文字数を計算し、DBに登録する処理(Google APIが上手く取得できなかったため、使用しない予定)
  └ web/        # HTTPサーバ
```

## 使い方
**各ディレクトリにあるREADMEも参考にしてください**

`config.json`で設定する。設定方法は`confgi_sample.json`を参照。`dir_id`はGoogle DriveにおけるディレクトリのID。

`.env`を作り、`MYSQL_ROOT_PASSWORD`, `MYSQL_USER`, `MYSQL_PASSWORD`を設定する。

```
$ docker-compose up -d

// mysqlに入りユーザにテーブルの操作権限を与えるなどしてください

// fetchFile/ web/において
$ make build

```

## 画面例
![進捗グラフ表示例](assets/screenshot.png)

