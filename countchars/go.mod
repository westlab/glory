module glory/countchars

go 1.15

require (
	cloud.google.com/go v0.72.0 // indirect
	github.com/Songmu/flextime v0.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.3
	github.com/westlab/glory v0.0.0
	github.com/westlab/glory/fetchFile v0.0.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	google.golang.org/api v0.35.0
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201112120144-2985b7af83de // indirect
)

replace github.com/westlab/glory v0.0.0 => ../

replace github.com/westlab/glory/fetchFile v0.0.0 => ./
