package config

import "testing"

func TestParseRedisURL(t *testing.T) {
	rc, err := ParseRedisURL("redis://:s3cret@cache.internal:6379/2")
	if err != nil {
		t.Fatal(err)
	}
	if rc.Addr != "cache.internal:6379" || rc.Password != "s3cret" || rc.DB != 2 {
		t.Fatalf("unexpected config: %+v", rc)
	}
}
