package isolate

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
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

	cleanCmdStr := fmt.Sprintf("isolate --cg --cleanup --box-id %d", boxId)
	logger = logger.With(slog.String("cmd", cleanCmdStr))

	cleanCmd := exec.Command("/usr/bin/bash", "-c", cleanCmdStr)
	cleanOut, err := cleanCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	logger.Info("ran isolate cleanup command", slog.String("output", string(cleanOut)))

	initCmdStr := fmt.Sprintf("isolate --cg --init --box-id %d", boxId)
	logger = logger.With(slog.String("cmd", initCmdStr))

	initCmd := exec.Command("/usr/bin/bash", "-c", initCmdStr)
	initOut, err := initCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	initOutStr := string(initOut)
	logger.Info("ran isolate init command", slog.String("output", initOutStr))

	boxPath := initOutStr
    boxPath = strings.TrimSuffix(boxPath, "\n")

	isolate.idsInUse = append(isolate.idsInUse, boxId)
	return NewIsolateBox(isolate, boxId, boxPath), nil
}

func (isolate *Isolate) EraseBox(boxId int) error {
	isolate.mutex.Lock()
	defer isolate.mutex.Unlock()
    logger := slog.With(slog.Int("box-id", boxId))

	cleanCmdStr := fmt.Sprintf("isolate --cg --cleanup --box-id %d", boxId)
	logger = logger.With(slog.String("cmd", cleanCmdStr))
	
    cleanCmd := exec.Command("/usr/bin/bash", "-c", cleanCmdStr)
	cleanOut, err := cleanCmd.CombinedOutput()
    
    logger = logger.With(slog.String("output", string(cleanOut)))
    logger.Info("erased isolate box")
	if err != nil {
		return err
	}
	return nil
}

func (isolate *Isolate) StartCommand(
	boxId int, command string, stdin io.ReadCloser,
	constraints RuntimeConstraints) (*IsolateProcess,error) {

    var process *IsolateProcess = &IsolateProcess{}
    var err error

	runCmdStr := fmt.Sprintf("isolate --cg --box-id %d %s %s --run /usr/bin/env %s",
		boxId,
        strings.Join(constraints.ToArgs(), " "),
        "--env=HOME=/box",
		command)

    logger := slog.With(slog.Int("box-id", boxId),
                        slog.String("cmd", runCmdStr))

    cmd := exec.Command("/usr/bin/bash", "-c", runCmdStr)
    cmd.Stdin = stdin
    process.stdout, err = cmd.StdoutPipe()
    if err != nil {
        return process, err
    }
    process.stderr, err = cmd.StderrPipe()
    if err != nil {
        return process, err
    }
    process.cmd = cmd

	if err = cmd.Start(); err != nil {
		return process, err
	}

    logger.Info("started isolate command")

	return process, err
}
