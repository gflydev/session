package session

import (
	"testing"
)

func newDict(kv map[string]interface{}) Dict {
	d := newDictValue()
	for k, v := range kv {
		d.KV[k] = v
	}
	return d
}

func TestMSGPEncodeDecodeRoundTrip(t *testing.T) {
	src := newDict(map[string]interface{}{
		"user":  "alice",
		"count": int64(42),
		"admin": true,
	})

	encoded, err := MSGPEncode(src)
	if err != nil {
		t.Fatalf("MSGPEncode: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatal("MSGPEncode returned empty payload for non-empty dict")
	}

	var dst Dict
	if err := MSGPDecode(&dst, encoded); err != nil {
		t.Fatalf("MSGPDecode: %v", err)
	}

	if got := dst.KV["user"]; got != "alice" {
		t.Errorf(`dst["user"] = %v, want "alice"`, got)
	}
	if got := dst.KV["count"]; got != int64(42) {
		t.Errorf(`dst["count"] = %v (%T), want int64(42)`, got, got)
	}
	if got := dst.KV["admin"]; got != true {
		t.Errorf(`dst["admin"] = %v, want true`, got)
	}
}

func TestMSGPEncodeEmptyReturnsNil(t *testing.T) {
	encoded, err := MSGPEncode(newDictValue())
	if err != nil {
		t.Fatalf("MSGPEncode: %v", err)
	}
	if encoded != nil {
		t.Errorf("MSGPEncode of empty dict = %v, want nil", encoded)
	}
}

func TestMSGPDecodeEmptyIsNoop(t *testing.T) {
	dst := newDict(map[string]interface{}{"stale": "value"})
	if err := MSGPDecode(&dst, nil); err != nil {
		t.Fatalf("MSGPDecode(nil): %v", err)
	}
	if len(dst.KV) != 0 {
		t.Errorf("MSGPDecode(nil) left %d entries, want 0 (decode should clear)", len(dst.KV))
	}
}

func TestBase64EncodeDecodeRoundTrip(t *testing.T) {
	src := newDict(map[string]interface{}{
		"token": "xyz",
		"n":     int64(7),
	})

	encoded, err := Base64Encode(src)
	if err != nil {
		t.Fatalf("Base64Encode: %v", err)
	}

	var dst Dict
	if err := Base64Decode(&dst, encoded); err != nil {
		t.Fatalf("Base64Decode: %v", err)
	}

	if got := dst.KV["token"]; got != "xyz" {
		t.Errorf(`dst["token"] = %v, want "xyz"`, got)
	}
	if got := dst.KV["n"]; got != int64(7) {
		t.Errorf(`dst["n"] = %v, want int64(7)`, got)
	}
}

func TestBase64DecodeInvalidReturnsError(t *testing.T) {
	var dst Dict
	if err := Base64Decode(&dst, []byte("not-valid-base64!!!")); err == nil {
		t.Fatal("Base64Decode of invalid input: want error, got nil")
	}
}
