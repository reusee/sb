package sb

import (
	"bytes"
	"testing"
)

func TestChunkedReader(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(
			ChunkedReader{
				R: bytes.NewReader([]byte("foobarbaz")),
				N: 3,
			},
		),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}

	var data [][]byte
	if err := Copy(
		Decode(buf), Unmarshal(

			&data)); err != nil {
		t.Fatal(err)
	}
	if len(data) != 3 {
		t.Fatal()
	}
	if !bytes.Equal(data[0], []byte("foo")) {
		t.Fatal()
	}
	if !bytes.Equal(data[1], []byte("bar")) {
		t.Fatal()
	}
	if !bytes.Equal(data[2], []byte("baz")) {
		t.Fatal()
	}
}

func TestMarshalReaderChunked(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := Copy(
		Marshal(
			MarshalReaderChunked(
				DefaultCtx,
				bytes.NewReader([]byte("foobarbaz")),
				3,
				nil,
			),
		),
		Encode(buf),
	); err != nil {
		t.Fatal(err)
	}

	var data [][]byte
	if err := Copy(
		Decode(buf), Unmarshal(

			&data)); err != nil {
		t.Fatal(err)
	}
	if len(data) != 3 {
		t.Fatal()
	}
	if !bytes.Equal(data[0], []byte("foo")) {
		t.Fatal()
	}
	if !bytes.Equal(data[1], []byte("bar")) {
		t.Fatal()
	}
	if !bytes.Equal(data[2], []byte("baz")) {
		t.Fatal()
	}
}

func BenchmarkChunkedReader(b *testing.B) {
	bs := bytes.Repeat([]byte("foo"), 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(bs)
		if err := Copy(
			Marshal(ChunkedReader{r, 128}),
			Discard,
		); err != nil {
			b.Fatal()
		}
	}
}

func BenchmarkMarshalReaderChunked(b *testing.B) {
	bs := bytes.Repeat([]byte("foo"), 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(bs)
		proc := MarshalReaderChunked(DefaultCtx, r, 128, nil)
		if err := Copy(
			&proc,
			Discard,
		); err != nil {
			b.Fatal()
		}
	}
}
