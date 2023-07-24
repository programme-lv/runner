package isolate

import (
	"fmt"
	"io"
	"os/exec"
)

type IsolateBox struct {
	id      int
	path    string
	isolate *Isolate
}

func NewIsolateBox(isolate *Isolate, id int, path string) *IsolateBox {
	return &IsolateBox{
		id:      id,
		path:    path,
		isolate: isolate,
	}
}

func (box *IsolateBox) Id() int {
	return box.id
}

func (box *IsolateBox) Close() error {
	// run isolate --cleanup --box-id box.id
	return box.isolate.EraseBox(box.id)
}

func (box *IsolateBox) Run(command string,
    constraints *RuntimeConstraints,
	stdin io.ReadCloser,
	stdout io.WriteCloser,
	stderr io.WriteCloser) (cmd *exec.Cmd, err error) {

	runCmdStr := fmt.Sprintf("isolate --box-id %d --run %s", box.id, command)
	cmd = exec.Command("/usr/bin/bash", "-c", runCmdStr)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err = cmd.Start(); err != nil {
		return 
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	return
}
