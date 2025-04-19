package goroutines

import (
	"io"
	"os"
)

type repeatReader struct {
	content string
	count   int
}

func (r *repeatReader) Read(p []byte) (n int, err error) {
	if r.count <= 0 {
		return 0, io.EOF
	}
	r.count--
	n = copy(p, []byte(r.content))
	return n, nil
}

func CustomReaderSample() {
	r := &repeatReader{
		content: "hello\n",
		count:   5,
	}
	io.Copy(os.Stdout, r)
}
