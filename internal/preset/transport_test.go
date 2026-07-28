package preset

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestShareCodeRoundTripsCanonicalPreset(t *testing.T) {
	t.Parallel()

	preset := Default(StyleMaia)
	code, err := EncodeShare(preset)
	if err != nil {
		t.Fatalf("EncodeShare: %v", err)
	}
	if !strings.HasPrefix(code, "gsxui:v1:") {
		t.Fatalf("share code = %q, want gsxui:v1 prefix", code)
	}
	if strings.Contains(strings.TrimPrefix(code, "gsxui:v1:"), "=") {
		t.Fatalf("share code contains base64 padding: %q", code)
	}

	got, err := DecodeShare(code)
	if err != nil {
		t.Fatalf("DecodeShare: %v", err)
	}
	if !presetsEqual(got, preset) {
		t.Fatalf("DecodeShare = %#v, want %#v", got, preset)
	}
	got.Theme.Light["primary"] = "red"
	if preset.Theme.Light["primary"] == "red" {
		t.Fatal("DecodeShare returned shared theme storage")
	}
}

func TestDecodeShareRejectsInvalidTransport(t *testing.T) {
	t.Parallel()

	canonical := readTestFile(t, "testdata/default-nova.json")
	encoded := base64.RawURLEncoding.EncodeToString(canonical)
	minified := bytes.ReplaceAll(canonical, []byte("\n"), nil)
	invalidUTF8 := base64.RawURLEncoding.EncodeToString([]byte{0xff, 0xfe})
	trailingBitsPreset := Default(StyleNova)
	trailingBitsPreset.Theme.Light["primary"] = "red"
	trailingBitsPayload, err := CanonicalJSON(trailingBitsPreset)
	if err != nil {
		t.Fatal(err)
	}
	trailingBitsEncoded := base64.RawURLEncoding.EncodeToString(trailingBitsPayload)
	trailingBitsAlternate := nonCanonicalTrailingBits(t, trailingBitsEncoded, trailingBitsPayload)

	tests := []struct {
		name string
		code string
		want string
	}{
		{name: "unknown version", code: "gsxui:v2:" + encoded, want: "version"},
		{name: "missing prefix", code: "v1:" + encoded, want: "prefix"},
		{name: "padded base64", code: "gsxui:v1:" + encoded + "=", want: "base64"},
		{name: "invalid base64", code: "gsxui:v1:not+url", want: "base64"},
		{
			name: "embedded carriage return",
			code: "gsxui:v1:" + encoded[:12] + "\r" + encoded[12:],
			want: "canonical",
		},
		{
			name: "embedded line feed",
			code: "gsxui:v1:" + encoded[:12] + "\n" + encoded[12:],
			want: "canonical",
		},
		{
			name: "non-zero unused trailing bits",
			code: "gsxui:v1:" + trailingBitsAlternate,
			want: "canonical",
		},
		{name: "invalid UTF-8", code: "gsxui:v1:" + invalidUTF8, want: "UTF-8"},
		{
			name: "non-canonical valid JSON",
			code: "gsxui:v1:" + base64.RawURLEncoding.EncodeToString(minified),
			want: "canonical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := DecodeShare(tt.code); err == nil ||
				!strings.Contains(err.Error(), tt.want) {
				t.Fatalf("DecodeShare error = %v, want %q", err, tt.want)
			}
		})
	}
}

func nonCanonicalTrailingBits(t *testing.T, encoded string, payload []byte) string {
	t.Helper()
	if remainder := len(payload) % 3; remainder != 1 && remainder != 2 {
		t.Fatalf("payload length %d has no unused trailing base64 bits", len(payload))
	}
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	last := strings.IndexByte(alphabet, encoded[len(encoded)-1])
	if last == -1 {
		t.Fatalf("last encoded byte %q is outside base64url alphabet", encoded[len(encoded)-1])
	}
	alternate := encoded[:len(encoded)-1] + string(alphabet[last+1])
	decoded, err := base64.RawURLEncoding.DecodeString(alternate)
	if err != nil {
		t.Fatalf("alternate base64 spelling did not decode: %v", err)
	}
	if !bytes.Equal(decoded, payload) {
		t.Fatal("alternate base64 spelling changed decoded payload")
	}
	return alternate
}

func TestInputResolverResolvesFileStdinRawJSONAndShare(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	code, err := EncodeShare(Default(StyleNova))
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "preset.json")
	if err := os.WriteFile(path, golden, 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		resolver InputResolver
		argument string
	}{
		{name: "file", resolver: InputResolver{}, argument: path},
		{
			name:     "stdin JSON",
			resolver: InputResolver{Stdin: bytes.NewReader(golden)},
			argument: "-",
		},
		{
			name:     "stdin share",
			resolver: InputResolver{Stdin: strings.NewReader(code + "\n")},
			argument: "-",
		},
		{name: "raw JSON", resolver: InputResolver{}, argument: string(golden)},
		{name: "share", resolver: InputResolver{}, argument: code},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.resolver.Resolve(context.Background(), tt.argument)
			if err != nil {
				t.Fatalf("Resolve: %v", err)
			}
			if !presetsEqual(got, Default(StyleNova)) {
				t.Fatalf("Resolve = %#v, want default Nova", got)
			}
		})
	}
}

func TestInputResolverLoadsHTTPSJSONAndShareWithBoundedContext(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	code, err := EncodeShare(Default(StyleNova))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		body string
	}{
		{name: "JSON", body: string(golden)},
		{name: "share", body: code},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				deadline, ok := request.Context().Deadline()
				if !ok {
					return nil, errors.New("request context has no deadline")
				}
				remaining := time.Until(deadline)
				if remaining <= 9*time.Second || remaining > 10*time.Second {
					return nil, errors.New("request context deadline is not 10 seconds")
				}
				return response(request, http.StatusOK, tt.body), nil
			})}
			got, err := (InputResolver{HTTPClient: client}).Resolve(
				context.Background(),
				"https://themes.example.test/preset",
			)
			if err != nil {
				t.Fatalf("Resolve: %v", err)
			}
			if !presetsEqual(got, Default(StyleNova)) {
				t.Fatalf("Resolve = %#v, want default Nova", got)
			}
		})
	}
}

func TestInputResolverRejectsUnsafeURLsAndResponses(t *testing.T) {
	t.Parallel()

	golden := readTestFile(t, "testdata/default-nova.json")
	var transportCalls atomic.Int32
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		transportCalls.Add(1)
		switch request.URL.Host {
		case "redirect.example.test":
			resp := response(request, http.StatusFound, "")
			resp.Header.Set("Location", "http://downgrade.example.test/preset")
			return resp, nil
		case "status.example.test":
			return response(request, http.StatusTeapot, string(golden)), nil
		case "invalid.example.test":
			return response(request, http.StatusOK, `{"not":"a preset"}`), nil
		default:
			return response(request, http.StatusOK, string(golden)), nil
		}
	})}

	tests := []struct {
		name     string
		argument string
		want     string
		noCall   bool
	}{
		{name: "HTTP", argument: "http://themes.example.test/preset", want: "https", noCall: true},
		{
			name:     "credentials",
			argument: "https://user:secret@themes.example.test/preset",
			want:     "credentials",
			noCall:   true,
		},
		{
			name:     "fragment",
			argument: "https://themes.example.test/preset#part",
			want:     "fragment",
			noCall:   true,
		},
		{
			name:     "redirect downgrade",
			argument: "https://redirect.example.test/preset",
			want:     "redirect",
		},
		{
			name:     "non-2xx",
			argument: "https://status.example.test/preset",
			want:     "418",
		},
		{
			name:     "invalid body",
			argument: "https://invalid.example.test/preset",
			want:     "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := transportCalls.Load()
			_, err := (InputResolver{HTTPClient: client}).Resolve(context.Background(), tt.argument)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Resolve error = %v, want %q", err, tt.want)
			}
			if tt.noCall && transportCalls.Load() != before {
				t.Fatal("unsafe URL reached the HTTP transport")
			}
		})
	}
}

func TestInputResolverRejectsOversizedInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resolver InputResolver
		argument string
	}{
		{
			name:     "raw JSON",
			resolver: InputResolver{MaxBytes: 5},
			argument: `{"a":1}`,
		},
		{
			name:     "direct share",
			resolver: InputResolver{MaxBytes: 5},
			argument: "gsxui:v1:abcdef",
		},
		{
			name:     "stdin",
			resolver: InputResolver{Stdin: strings.NewReader("123456"), MaxBytes: 5},
			argument: "-",
		},
		{
			name: "HTTPS body",
			resolver: InputResolver{
				MaxBytes: 5,
				HTTPClient: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
					return response(request, http.StatusOK, "123456"), nil
				})},
			},
			argument: "https://themes.example.test/preset",
		},
	}

	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "large.json")
	if err := os.WriteFile(path, []byte("123456"), 0o600); err != nil {
		t.Fatal(err)
	}
	tests = append(tests, struct {
		name     string
		resolver InputResolver
		argument string
	}{name: "file", resolver: InputResolver{MaxBytes: 5}, argument: path})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := tt.resolver.Resolve(context.Background(), tt.argument); err == nil ||
				!strings.Contains(err.Error(), "5 bytes") {
				t.Fatalf("Resolve error = %v, want 5 byte limit", err)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func response(request *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    request,
	}
}
