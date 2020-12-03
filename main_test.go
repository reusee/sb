package sb

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {

	//bytesPool8.LogCallers = true
	//bytesPool32K.LogCallers = true

	defer func() {
		for _, stack := range bytesPool8.Callers {
			if len(stack) > 0 {
				panic(fmt.Errorf("pool leak: %s", stack))
			}
		}
		for _, stack := range bytesPool32K.Callers {
			if len(stack) > 0 {
				panic(fmt.Errorf("pool leak: %s", stack))
			}
		}
	}()

	m.Run()
}
