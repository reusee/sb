package sb

import (
	"fmt"
	"io"
)

func dumpStream(
	stream Stream,
	prefix string,
	w io.Writer,
) {
	for {
		token, err := stream.Next()
		if err != nil { // NOCOVER
			panic(err)
		}
		if token == nil {
			break
		}
		fmt.Fprintf(w, "%s%+v\n", prefix, token)
	}
}
