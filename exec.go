package goutil

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Exec struct {
}

func (this *Exec) Exec(cmd []string, timeout int) (string, int) {
	var out bytes.Buffer
	duration := time.Duration(timeout) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), duration)
	var command *exec.Cmd
	command = exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		log.Println(err, cmd)
		return err.Error(), -1
	}
	status := command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	return out.String(), status
}
