package sb

import "fmt"

type Offset int64

var _ error = Offset(0)

func (o Offset) Error() string {
	return fmt.Sprintf("offset: %d", o)
}
