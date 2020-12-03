package sb

import (
	"github.com/reusee/pr"
)

// 8 * 1024 = 8K
var bytesPool8 = pr.NewPool(1024, func() any {
	bs := make([]byte, 8)
	return &bs
})

// 32K * 32 = 1M
var bytesPool32K = pr.NewPool(32, func() any {
	bs := make([]byte, 32*1024)
	return &bs
})
