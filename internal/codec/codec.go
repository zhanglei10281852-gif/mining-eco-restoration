package codec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(v any) (string, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func Decode(s string, v any) error {
	b, e := base64.RawURLEncoding.DecodeString(s)
	if e != nil {
		return fmt.Errorf("decode: %w", e)
	}
	if e = json.Unmarshal(b, v); e != nil {
		return fmt.Errorf("json: %w", e)
	}
	return nil
}
func Must(v any) string {
	s, e := Encode(v)
	if e != nil {
		panic(e)
	}
	return s
}
