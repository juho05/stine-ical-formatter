package formatter

import (
	"io"

	"golang.org/x/text/encoding/unicode"
)

func newReader(r io.Reader) io.Reader {
	// STiNE exports have no BOM
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	return decoder.Reader(r)
}
