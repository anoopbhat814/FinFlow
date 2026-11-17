package main

// Populated at link time: -ldflags "-X main.version=v0.4.0 -X main.commit=abc1234".
var (
	version = "dev"
	commit  = "unknown"
)

func buildInfo() map[string]string {
	return map[string]string{"version": version, "commit": commit}
}
