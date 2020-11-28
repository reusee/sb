package sb

import (
	"github.com/reusee/pr"
)

// 8K
var getEightBytes = pr.NewBytesPool(8, 1024)

// 1M
var get32KBytes = pr.NewBytesPool(32*1024, 32)
