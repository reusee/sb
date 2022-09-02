//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/reusee/e5"
	"github.com/reusee/sb"
	"os"
)

func main() {

	// marshal stream
	marshaler := sb.Marshal(42)

	// unmarshal sink
	var num int
	unmarshaler := sb.Unmarshal(&num)

	// copy stream to sink
	check(sb.Copy(
		marshaler,
		unmarshaler,
	))
	must(num == 42)

	// encode
	buf := new(bytes.Buffer)
	check(sb.Copy(
		sb.Marshal(80),
		sb.Encode(buf),
	))

	// decode
	check(sb.Copy(
		sb.Decode(buf),
		sb.Unmarshal(&num),
	))
	must(num == 80)

	// hash
	var sum []byte
	check(sb.Copy(
		sb.Marshal(map[int]string{
			42: "42",
			80: "80",
		}),
		sb.Hash(sha256.New, &sum, nil),
	))
	must(fmt.Sprintf("%x", sum) ==
		"48fc94ae5e9f6961bbf6e85288deab361e9d1bac2d13be8fa20ee4103295d033")

	var tokens sb.Tokens
	check(sb.Copy(
		// stream combinator
		sb.Tee(
			sb.Marshal(42),
			sb.Encode(buf),
		),
		// multiple sinks
		sb.Unmarshal(&num),
		sb.Hash(sha256.New, &sum, nil),
		sb.CollectTokens(&tokens),
	))
	must(num == 42)
	must(fmt.Sprintf("%x", sum) ==
		"151a3a0b4c88483512fc484d0badfedf80013ebb18df498bbee89ac5b69d7222")

	// stream comparison
	res, err := sb.Compare(
		sb.Marshal(42),
		sb.Marshal(map[int]int{42: 42}),
	)
	check(err)
	must(res == -1)

}

var (
	check = e5.Check
	pt    = fmt.Printf
)

func must(b bool) {
	if !b {
		pt(e5.NewStacktrace()(fmt.Errorf("should be true")).Error())
		os.Exit(-1)
	}
}
