package isolate

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type IsolateResult struct {
}

type IsolateProcess struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func (process *IsolateProcess) Wait() error {
	err := process.cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (process *IsolateProcess) LogOutput() {
    stdoutScanner := bufio.NewScanner(process.Stdout())
    stderrScanner := bufio.NewScanner(process.Stderr())
	go func() {
		for stdoutScanner.Scan() {
			fmt.Printf(stdoutScanner.Text())
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			fmt.Printf(stderrScanner.Text())
		}
	}()
}

func (process *IsolateProcess) Stdout() io.ReadCloser {
	return process.stdout
}

func (process *IsolateProcess) Stderr() io.ReadCloser {
	return process.stderr
}
