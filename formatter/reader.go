package formatter

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/charmap"
)

type reader struct {
	r io.Reader
}

func newReader(r io.Reader) io.Reader {
	decoder := charmap.ISO8859_1.NewDecoder()
	return decoder.Reader(reader{
		r: r,
	})
}

func (r reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	data := bytes.ReplaceAll(p[:n], []byte{0}, []byte{})
	copy(p, data)
	return len(data), err
}
