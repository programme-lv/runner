package isolate

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"golang.org/x/exp/slog"
)

type Isolate struct {
	idsInUse []int
	mutex    sync.Mutex
}

func NewIsolate() (*Isolate, error) {
	versionCmdStr := fmt.Sprintf("isolate --version")
	logger := slog.With(slog.String("cmd", versionCmdStr))

	versionCmd := exec.Command("/usr/bin/bash", "-c", versionCmdStr)
	out, err := versionCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	logger.Info("ran isolate version command", slog.String("output", string(out)))

	return &Isolate{}, nil
}

func (isolate *Isolate) isBoxIdInUse(boxId int) bool {
	for _, idInUse := range isolate.idsInUse {
		if idInUse == boxId {
			return true
		}
	}
	return false
}

func (isolate *Isolate) NewBox() (*IsolateBox, error) {
	isolate.mutex.Lock()
	defer isolate.mutex.Unlock()

	boxId := 0
	for isolate.isBoxIdInUse(boxId) {
		boxId++
	}

	logger := slog.With(slog.Int("box-id", boxId))

	cleanCmdStr := fmt.Sprintf("isolate --cleanup --box-id %d", boxId)
	logger = logger.With(slog.String("cmd", cleanCmdStr))

	cleanCmd := exec.Command("/usr/bin/bash", "-c", cleanCmdStr)
	cleanOut, err := cleanCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	logger.Info("ran isolate cleanup command", slog.String("output", string(cleanOut)))

	initCmdStr := fmt.Sprintf("isolate --init --box-id %d", boxId)
	logger = logger.With(slog.String("cmd", initCmdStr))

	initCmd := exec.Command("/usr/bin/bash", "-c", initCmdStr)
	initOut, err := initCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	initOutStr := string(initOut)
	logger.Info("ran isolate init command", slog.String("output", initOutStr))

	boxPath := initOutStr

	isolate.idsInUse = append(isolate.idsInUse, boxId)
	return NewIsolateBox(isolate, boxId, boxPath), nil
}

func (isolate *Isolate) EraseBox(boxId int) error {
	isolate.mutex.Lock()
	defer isolate.mutex.Unlock()

	cleanCmdStr := fmt.Sprintf("isolate --cleanup --box-id %d", boxId)
	cleanCmd := exec.Command("/usr/bin/bash", "-c", cleanCmdStr)
	cleanOut, err := cleanCmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Println(string(cleanOut))
	return nil
}

func (isolate *Isolate) RunCommand(
	boxId int,
	command string,
	constraints RuntimeConstraints,
	stdin io.ReadCloser,
	stdout io.WriteCloser,
	stderr io.WriteCloser) (cmd *exec.Cmd, err error) {

	runCmdStr := fmt.Sprintf("isolate --box-id %d --run %s",
		boxId,
		command)

	cmd = exec.Command("/usr/bin/bash", "-c", runCmdStr)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err = cmd.Start(); err != nil {
		return
	}

	return
}
