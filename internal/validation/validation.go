package validation

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func Email(v string) bool { return emailPattern.MatchString(strings.ToLower(strings.TrimSpace(v))) }
func Required(v string, max int) bool {
	v = strings.TrimSpace(v)
	return v != "" && utf8.RuneCountInString(v) <= max
}
func OneOf(v string, values ...string) bool {
	for _, x := range values {
		if v == x {
			return true
		}
	}
	return false
}
func Score(v int) bool          { return v >= 0 && v <= 100 }
func Area(v float64) bool       { return v > 0 && v < 1e9 }
func Normalize(v string) string { return strings.Join(strings.Fields(strings.TrimSpace(v)), " ") }
