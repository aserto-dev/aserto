package iostream

import (
	"io"
	"os"
)

// type IO interface {
// 	StdIn() io.Reader
// 	StdOut() io.Writer
// 	StdErr() io.Writer
// }

type StdIO struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (i *StdIO) StdIn() io.Reader {
	return i.in
}

func (i *StdIO) StdOut() io.Writer {
	return i.out
}

func (i *StdIO) StdErr() io.Writer {
	return i.err
}

// type BufferIO struct {
// 	In  *bytes.Buffer
// 	Out *bytes.Buffer
// 	Err *bytes.Buffer
// }

// func (i *BufferIO) Input() io.Reader {
// 	return i.In
// }

// func (i *BufferIO) Output() io.Writer {
// 	return i.Out
// }

// func (i *BufferIO) Error() io.Writer {
// 	return i.Err
// }

func DefaultIO() *StdIO {
	return &StdIO{os.Stdin, os.Stdout, os.Stderr}
}

// func BytesIO() *BufferIO {
// 	return &BufferIO{&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}}
// }

// func NewUI(ios IO) *StdIO {
// 	return DefaultIO()
// }
