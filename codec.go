package ffxiv

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

var decoderMap = map[uint16]Decoder{
	0x0:   decodeNone,
	0x1:   decodeNone,
	0x100: decodeZlib,
	0x101: decodeZlib,
}

// Decoder is a type to represent different encodings
type Decoder func([]byte) ([]byte, error)

// Decode returns a decoded byte slice for the provided data given a specified encoding.
func Decode(encoding uint16, data []byte) ([]byte, error) {
	e, ok := decoderMap[encoding]
	if !ok {
		return nil, errors.New("unsupported codec")
	}
	return e(data)
}

func decodeNone(b []byte) ([]byte, error) {
	return b, nil
}

func decodeZlib(data []byte) ([]byte, error) {
	b := make([]byte, 0)

	in := bytes.NewBuffer(data)
	out := bytes.NewBuffer(b)
	r, err := zlib.NewReader(in)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	io.Copy(out, r)
	return out.Bytes(), nil
}
