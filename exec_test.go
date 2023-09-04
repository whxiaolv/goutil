package goutil

import (
	"fmt"
	"testing"
)

func Test_exec(t *testing.T) {
	t.Run("testExec", testExec)
	//t.Run("testUrlEncode", testUrlEncode)
}

var (
	execCommand Exec
)

func testExec(t *testing.T) {

	execString, status := execCommand.Exec([]string{"ls"}, 1)
	fmt.Println(execString, status)

}
