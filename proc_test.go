package sb

import (
	"testing"
)

func TestProcAsMarshaler(t *testing.T) {
	var i int
	if err := Copy(
		Marshal(
			Marshal(
				Marshal(
					Marshal(
						Marshal(
							Marshal(
								Marshal(
									42,
								),
							),
						),
					),
				),
			),
		),
		Unmarshal(&i),
	); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal()
	}

	var strs []string
	if err := Copy(
		Marshal(
			Marshal(
				Marshal(
					Marshal(
						Marshal(
							Marshal(
								Marshal(
									[]any{
										Marshal(
											Marshal(
												Marshal("foo"),
											),
										),
										"bar",
										Marshal("baz"),
									},
								),
							),
						),
					),
				),
			),
		),
		Unmarshal(&strs),
	); err != nil {
		t.Fatal(err)
	}
	if len(strs) != 3 {
		t.Fatal()
	}
	if strs[0] != "foo" {
		t.Fatal()
	}
	if strs[1] != "bar" {
		t.Fatal()
	}
	if strs[2] != "baz" {
		t.Fatal()
	}
}
