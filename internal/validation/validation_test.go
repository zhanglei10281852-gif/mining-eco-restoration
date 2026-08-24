package validation_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/validation"
	"testing"
)

func TestEmailCases(t *testing.T) {
	good := []string{"a@b.co", "ops.team@example.cn"}
	bad := []string{"", "a@", "a b@c.com", "a@b"}
	for _, v := range good {
		if !validation.Email(v) {
			t.Fatal(v)
		}
	}
	for _, v := range bad {
		if validation.Email(v) {
			t.Fatal(v)
		}
	}
}
func TestRequiredAndNormalize(t *testing.T) {
	if !validation.Required(" name ", 10) {
		t.Fatal("required")
	}
	if validation.Required("", 10) {
		t.Fatal("empty")
	}
	if validation.Required("12345", 4) {
		t.Fatal("max")
	}
	if got := validation.Normalize("  a   b "); got != "a b" {
		t.Fatal(got)
	}
}
func TestEnumsAndBounds(t *testing.T) {
	if !validation.OneOf("x", "a", "b", "x") {
		t.Fatal()
	}
	if validation.OneOf("z", "a") {
		t.Fatal()
	}
	for _, v := range []int{0, 100} {
		if !validation.Score(v) {
			t.Fatal(v)
		}
	}
	if validation.Score(-1) || validation.Score(101) {
		t.Fatal()
	}
	if validation.Area(0) || !validation.Area(10) {
		t.Fatal()
	}
	if validation.Area(-2) {
		t.Fatal()
	}
}
