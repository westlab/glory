package glory

import "os"

var DataSourceName string

func init() {
	DataSourceName = os.Getenv("DATA_SOURCE_NAME")
	if DataSourceName == "" {
		DataSourceName = "user:userpass@tcp(127.0.0.1:3306)/glory_test"
	}
}
