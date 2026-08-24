package pagination_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/pagination"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	tests := []struct {
		l, o         string
		wantL, wantO int
	}{{"", "", 20, 0}, {"-1", "-2", 20, 0}, {"101", "4", 20, 4}, {"5", "7", 5, 7}, {"abc", "x", 20, 0}}
	for _, tt := range tests {
		p := pagination.Parse(tt.l, tt.o)
		if p.Limit != tt.wantL || p.Offset != tt.wantO {
			t.Fatalf("%+v got %+v", tt, p)
		}
	}
}
func TestParseManyLimits(t *testing.T) {
	for i := 0; i < 150; i++ {
		p := pagination.Parse("100", "0")
		if p.Limit != 100 {
			t.Fatal(p)
		}
	}
}
