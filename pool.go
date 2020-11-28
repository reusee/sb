package sb

import "sync"

var eightBytesPool = sync.Pool{
	New: func() any {
		bs := make([]byte, 8)
		return &bs
	},
}

var copyBufferPool = sync.Pool{
	New: func() any {
		bs := make([]byte, 32*1024)
		return &bs
	},
}
