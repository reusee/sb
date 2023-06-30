package sb

import "github.com/reusee/pr3"

// 8 * 1024 = 8K
var bytesPool8 = pr3.NewPool(1024, func() []byte {
	return make([]byte, 8)
})

// 32K * 32 = 1M
var bytesPool32K = pr3.NewPool(32, func() []byte {
	return make([]byte, 32*1024)
})
