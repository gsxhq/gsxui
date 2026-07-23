package main

import "testing"

// The dev loop's port agreement: with NO env set, the server must default to
// 7777 — the same default vite.config.ts and `gsx dev` assume — or a fresh
// checkout's `make site-dev` shows the "backend unavailable" interstitial
// forever (server on one port, proxy on another). PORT stays the container
// path (Dockerfile sets ENV PORT=8080; Cloud Run injects it), GO_PORT wins
// in dev.
func TestListenPort(t *testing.T) {
	env := func(m map[string]string) func(string) string {
		return func(k string) string { return m[k] }
	}
	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"no env → dev default 7777 (vite/gsx-dev agreement)", nil, "7777"},
		{"PORT only (container/Cloud Run)", map[string]string{"PORT": "8080"}, "8080"},
		{"GO_PORT wins over PORT (dev loop)", map[string]string{"GO_PORT": "7878", "PORT": "8080"}, "7878"},
		{"GO_PORT only", map[string]string{"GO_PORT": "7777"}, "7777"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := listenPort(env(tc.env)); got != tc.want {
				t.Fatalf("listenPort = %q, want %q", got, tc.want)
			}
		})
	}
}
