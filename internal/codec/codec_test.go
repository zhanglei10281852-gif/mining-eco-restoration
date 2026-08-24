package codec_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/codec"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	in := map[string]any{"name": "矿山", "count": 3}
	s, e := codec.Encode(in)
	if e != nil {
		t.Fatal(e)
	}
	var out map[string]any
	if e = codec.Decode(s, &out); e != nil || out["name"] != "矿山" {
		t.Fatal(out, e)
	}
}
func TestDecodeInvalid(t *testing.T) {
	var out any
	if e := codec.Decode("not-base64", &out); e == nil {
		t.Fatal()
	}
	if codec.Must(map[string]int{"x": 1}) == "" {
		t.Fatal()
	}
}
