package preset

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	sharePrefix          = "gsxui:v1:"
	defaultMaxInputBytes = int64(1 << 20)
	requestTimeout       = 10 * time.Second
)

type InputResolver struct {
	Stdin      io.Reader
	HTTPClient *http.Client
	MaxBytes   int64
}

func EncodeShare(preset Preset) (string, error) {
	canonical, err := CanonicalJSON(preset)
	if err != nil {
		return "", fmt.Errorf("encode share: %w", err)
	}
	return sharePrefix + base64.RawURLEncoding.EncodeToString(canonical), nil
}

func DecodeShare(code string) (Preset, error) {
	if !strings.HasPrefix(code, "gsxui:") {
		return Preset{}, errors.New("decode share: invalid prefix")
	}
	parts := strings.SplitN(code, ":", 3)
	if len(parts) != 3 || parts[1] != "v1" {
		version := ""
		if len(parts) > 1 {
			version = parts[1]
		}
		return Preset{}, fmt.Errorf("decode share: unsupported transport version %q", version)
	}
	if strings.Contains(parts[2], "=") {
		return Preset{}, errors.New("decode share: padded base64 is not canonical")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Preset{}, fmt.Errorf("decode share: invalid base64: %w", err)
	}
	if parts[2] != base64.RawURLEncoding.EncodeToString(payload) {
		return Preset{}, errors.New("decode share: base64 spelling is not canonical")
	}
	if !utf8.Valid(payload) {
		return Preset{}, errors.New("decode share: payload is not valid UTF-8")
	}
	preset, err := ParseJSON(payload)
	if err != nil {
		return Preset{}, fmt.Errorf("decode share: %w", err)
	}
	canonical, err := CanonicalJSON(preset)
	if err != nil {
		return Preset{}, fmt.Errorf("decode share: normalize payload: %w", err)
	}
	if !bytes.Equal(payload, canonical) {
		return Preset{}, errors.New("decode share: payload is valid but not canonical JSON")
	}
	return clonePreset(preset), nil
}

func (resolver InputResolver) Resolve(ctx context.Context, argument string) (Preset, error) {
	if argument == "-" {
		if resolver.Stdin == nil {
			return Preset{}, errors.New("resolve preset from stdin: no stdin reader configured")
		}
		src, err := readBounded(resolver.Stdin, resolver.maxBytes())
		if err != nil {
			return Preset{}, fmt.Errorf("resolve preset from stdin: %w", err)
		}
		return resolvePresetContent(src, "stdin")
	}
	if strings.HasPrefix(argument, "gsxui:") {
		if err := checkInputSize([]byte(argument), resolver.maxBytes()); err != nil {
			return Preset{}, fmt.Errorf("resolve share preset: %w", err)
		}
		return DecodeShare(argument)
	}
	if startsJSONDocument(argument) {
		if err := checkInputSize([]byte(argument), resolver.maxBytes()); err != nil {
			return Preset{}, fmt.Errorf("resolve raw preset JSON: %w", err)
		}
		preset, err := ParseJSON([]byte(argument))
		if err != nil {
			return Preset{}, fmt.Errorf("resolve raw preset JSON: %w", err)
		}
		return clonePreset(preset), nil
	}

	parsedURL, urlErr := url.Parse(argument)
	if urlErr == nil && parsedURL.Scheme != "" {
		return resolver.resolveHTTPS(ctx, parsedURL)
	}

	file, err := os.Open(argument)
	if err != nil {
		if urlErr != nil {
			return Preset{}, fmt.Errorf("resolve preset path %q (URL parse also failed: %v): %w", argument, urlErr, err)
		}
		return Preset{}, fmt.Errorf("resolve preset path %q: %w", argument, err)
	}
	defer file.Close()
	src, err := readBounded(file, resolver.maxBytes())
	if err != nil {
		return Preset{}, fmt.Errorf("resolve preset path %q: %w", argument, err)
	}
	return resolvePresetContent(src, argument)
}

func (resolver InputResolver) resolveHTTPS(ctx context.Context, parsedURL *url.URL) (Preset, error) {
	if !strings.EqualFold(parsedURL.Scheme, "https") {
		return Preset{}, fmt.Errorf("resolve preset URL: only https URLs are allowed, got %q", parsedURL.Scheme)
	}
	if parsedURL.Host == "" {
		return Preset{}, errors.New("resolve preset URL: https URL requires a host")
	}
	if parsedURL.User != nil {
		return Preset{}, errors.New("resolve preset URL: credentials are not allowed")
	}
	if parsedURL.Fragment != "" {
		return Preset{}, errors.New("resolve preset URL: fragments are not allowed")
	}

	client := resolver.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	clientCopy := *client
	originalRedirect := client.CheckRedirect
	clientCopy.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if !strings.EqualFold(request.URL.Scheme, "https") {
			return fmt.Errorf("redirect to non-HTTPS scheme %q is not allowed", request.URL.Scheme)
		}
		if request.URL.User != nil {
			return errors.New("redirect credentials are not allowed")
		}
		if request.URL.Fragment != "" {
			return errors.New("redirect fragments are not allowed")
		}
		if originalRedirect != nil {
			return originalRedirect(request, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	requestContext, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return Preset{}, fmt.Errorf("resolve preset URL: %w", err)
	}
	response, err := clientCopy.Do(request)
	if err != nil {
		return Preset{}, fmt.Errorf("resolve preset URL %q: %w", parsedURL.String(), err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return Preset{}, fmt.Errorf("resolve preset URL %q: unexpected HTTP status %d", parsedURL.String(), response.StatusCode)
	}

	maxBytes := resolver.maxBytes()
	if response.ContentLength > maxBytes {
		return Preset{}, fmt.Errorf("resolve preset URL %q: input exceeds %d bytes", parsedURL.String(), maxBytes)
	}
	src, err := readBounded(response.Body, maxBytes)
	if err != nil {
		return Preset{}, fmt.Errorf("resolve preset URL %q: %w", parsedURL.String(), err)
	}
	return resolvePresetContent(src, parsedURL.String())
}

func resolvePresetContent(src []byte, source string) (Preset, error) {
	trimmed := bytes.TrimSpace(src)
	if bytes.HasPrefix(trimmed, []byte("gsxui:")) {
		preset, err := DecodeShare(string(trimmed))
		if err != nil {
			return Preset{}, fmt.Errorf("resolve preset from %s: %w", source, err)
		}
		return clonePreset(preset), nil
	}
	preset, err := ParseJSON(trimmed)
	if err != nil {
		return Preset{}, fmt.Errorf("resolve preset from %s: %w", source, err)
	}
	return clonePreset(preset), nil
}

func startsJSONDocument(argument string) bool {
	decoder := json.NewDecoder(strings.NewReader(argument))
	token, err := decoder.Token()
	if err != nil {
		return false
	}
	delimiter, ok := token.(json.Delim)
	return ok && (delimiter == '{' || delimiter == '[')
}

func readBounded(reader io.Reader, maxBytes int64) ([]byte, error) {
	limit := maxBytes + 1
	if maxBytes == math.MaxInt64 {
		limit = maxBytes
	}
	src, err := io.ReadAll(io.LimitReader(reader, limit))
	if err != nil {
		return nil, err
	}
	if int64(len(src)) > maxBytes {
		return nil, checkInputSize(src, maxBytes)
	}
	return src, nil
}

func checkInputSize(src []byte, maxBytes int64) error {
	if int64(len(src)) > maxBytes {
		return fmt.Errorf("input exceeds %d bytes", maxBytes)
	}
	return nil
}

func (resolver InputResolver) maxBytes() int64 {
	if resolver.MaxBytes <= 0 {
		return defaultMaxInputBytes
	}
	return resolver.MaxBytes
}
