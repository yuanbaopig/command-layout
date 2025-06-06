package bufferpool

import (
	"bytes"
	"sync"
)

var (
	maxSize = 1 << 16
	bf      bytesBuffer
)

func init() {
	bf = bytesBuffer{
		Pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

type bytesBuffer struct {
	sync.Pool
}

func (b *bytesBuffer) getBuffer() *bytes.Buffer {
	return b.Get().(*bytes.Buffer)
}

func (b *bytesBuffer) putBuffer(buf *bytes.Buffer) {
	buf.Reset()

	if buf.Cap() > maxSize {
		return
	}

	b.Put(buf)
}

func SetBytesBuffSize(size int) {
	maxSize = size
}

func GetBytesBuffer() *bytes.Buffer {
	return bf.getBuffer()
}

func PutBytesBuffer(buf *bytes.Buffer) {
	bf.putBuffer(buf)
}
