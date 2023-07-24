package isolate

import (
	"io"
	"os/exec"
)

type IsolateResult struct {
}

type IsolateProcess struct {
    cmd *exec.Cmd
    stdout io.ReadCloser
    stderr io.ReadCloser
}

func (process *IsolateProcess) Wait() (*IsolateResult, error) {
    err := process.cmd.Wait()
    if err != nil {
        return nil, err
    }
    return nil, nil
}

func (process *IsolateProcess) Stdout() io.ReadCloser {
    return process.stdout
}

func (process *IsolateProcess) Stderr() io.ReadCloser {
    return process.stderr
}

