package glory

import (
	"fmt"
	"log"
	"os/exec"
)

// SetupTest はテスト開始時に行うコマンドを実行する
func SetupTest(cmds []string) {
	for _, c := range cmds {
		cmd := exec.Command("sh", "-c", c)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("set up test command %s, error: %v", c, err)
		}
	}
}

// TearDown はテスト終了時に行うコマンドを実行する
func TearDown(cmds []string) {
	for _, c := range cmds {
		cmd := exec.Command("sh", "-c", c)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("tear down command %s, error: %v", c, err)
		}
	}
}

func ExecSQL(file string) string {
	return fmt.Sprintf("mysql -u user -puserpass -h 127.0.0.1 < %s", file)
}
