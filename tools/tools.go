package tools

import (
	"io"
	"os/exec"
	"runtime"
)

func ExecClear(stdout io.Writer) {
	if stdout == nil {
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = stdout
	cmd.Run()
}
