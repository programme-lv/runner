package isolate

import (
	"io"
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
	stdin io.ReadCloser,
	constraints *RuntimeConstraints) (*IsolateProcess, error) {

	if constraints == nil {
		c := DefaultRuntimeConstraints()
		constraints = &c
	}

	return box.isolate.StartCommand(box.id, command, stdin, *constraints)
}
