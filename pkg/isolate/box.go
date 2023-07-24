package isolate

import (
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

func (box *IsolateBox) Path() string {
    return box.path
}

func (box *IsolateBox) Close() error {
	return box.isolate.EraseBox(box.id)
}

func (box *IsolateBox) Run(
	command string,
	stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser,
	constraints *RuntimeConstraints) (cmd *exec.Cmd, err error) {

	if constraints == nil {
		c := DefaultRuntimeConstraints()
		constraints = &c
	}

	return box.isolate.RunCommand(box.id, command, stdin, stdout, stderr, *constraints)
}
