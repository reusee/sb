package sb

import "github.com/reusee/pr2"

// 8 * 1024 = 8K
var bytesPool8 = pr2.NewPool(1024, func(_ pr2.PoolPutFunc) *[]byte {
	bs := make([]byte, 8)
	return &bs
})

// 32K * 32 = 1M
var bytesPool32K = pr2.NewPool(32, func(_ pr2.PoolPutFunc) *[]byte {
	bs := make([]byte, 32*1024)
	return &bs
})
