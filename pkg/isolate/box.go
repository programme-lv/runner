package isolate

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
    return box.isolate.ReleaseBox(box.id)
}
