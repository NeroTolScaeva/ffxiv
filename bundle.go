package ffxiv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
)

// magicIdentifier is the integer equivalent of the magic bytes used to identify the header of FFXIV bundles
const magicIdentifier uint32 = 0x41a05252

// BundleHeaderSize is the known length of a bundle header
const BundleHeaderSize = 40

// magic will store the []byte version of magicIdentifier here to avoid recomputing it
var magic []byte
var magicOnce sync.Once

// BundleHeader is the header of an FFXIV message
type BundleHeader struct {
	Magic      [4]uint32 `json:"magic"`
	Timestamp  uint64    `json:"timestamp"`
	Length     uint32    `json:"bundle_length"`
	Connection uint16    `json:"connection_type"`
	Count      uint16    `json:"message_count"`
	Encoding   uint16    `json:"encoding"`
	Unknown    [4]byte   `json:"-"`
}

// NewBundleHeader returns a new header from a given byte array.  The byte
// array's size is expected to be exactly BundleHeaderSize bytes.
func NewBundleHeader(b []byte) (*BundleHeader, error) {
	if len(b) != BundleHeaderSize {
		return nil, errors.New("invalid bundle header size")
	}
	r := bytes.NewReader(b)
	h := &BundleHeader{}
	if err := binary.Read(r, binary.LittleEndian, h); err != nil {
		return nil, err
	}
	return h, nil
}

// MagicIdentifier is a function to lazily return the unique identifier of an FFXIV bundle
// headers based on a known constant.  This function will panic if something goes wrong with
// the computation
func MagicIdentifier() []byte {
	magicOnce.Do(func() {
		b := bytes.NewBuffer(make([]byte, 0, 4))
		err := binary.Write(b, binary.LittleEndian, magicIdentifier)
		if err != nil {
			// An error here is unrecoverable
			panic(err)
		}
		magic = b.Bytes()
	})
	return magic
}
