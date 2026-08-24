package ids_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/ids"
	"strings"
	"testing"
)

func TestNewPrefixAndUniqueness(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		v := ids.New("task")
		if !strings.HasPrefix(v, "task_") {
			t.Fatal(v)
		}
		if seen[v] {
			t.Fatalf("duplicate %s", v)
		}
		seen[v] = true
	}
}
func TestPrefixes(t *testing.T) {
	for _, p := range []string{"usr", "ses", "aud", "evt", "smp"} {
		if !strings.HasPrefix(ids.New(p), p+"_") {
			t.Fatal(p)
		}
	}
}
