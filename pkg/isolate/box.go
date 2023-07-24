package isolate

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

type IsolateBox struct {
	id      int
	path    string
	isolate *Isolate
	logger  *slog.Logger
}

func NewIsolateBox(isolate *Isolate, id int, path string) *IsolateBox {
	return &IsolateBox{
		id:      id,
		path:    path,
		isolate: isolate,
        logger:  slog.With(slog.Int("box-id", id), slog.String("box-path", path)),
	}
}

func (box *IsolateBox) Id() int {
	return box.id
}

func (box *IsolateBox) Path() string {
	return box.path
}

func (box *IsolateBox) Close() error {
	return box.isolate.EraseBox(box.id)
}

func (box *IsolateBox) Run(
	command string,
	stdin io.ReadCloser,
	constraints *RuntimeConstraints) (*IsolateProcess, error) {

	if constraints == nil {
		c := DefaultRuntimeConstraints()
		constraints = &c
	}

	return box.isolate.StartCommand(box.id, command, stdin, *constraints)
}

func (box *IsolateBox) AddFile(path string, content []byte) error {
    box.logger.Info("adding file to box", slog.String("file-path", path))
	path = filepath.Join(box.path, "box", path)
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
    box.logger.Info("added file to box", slog.String("file-path", path))
	return nil
}
