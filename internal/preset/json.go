package preset

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"strconv"
	"unicode/utf8"
)

func ParseJSON(src []byte) (Preset, error) {
	if !utf8.Valid(src) {
		return Preset{}, errors.New("preset JSON: invalid UTF-8")
	}

	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	preset, err := decodePreset(decoder)
	if err != nil {
		return Preset{}, fmt.Errorf("preset JSON: %w", err)
	}

	if _, err := decoder.Token(); err == nil {
		return Preset{}, errors.New("preset JSON: trailing JSON value")
	} else if !errors.Is(err, io.EOF) {
		return Preset{}, fmt.Errorf("preset JSON: trailing data: %w", err)
	}
	if err := Validate(preset); err != nil {
		return Preset{}, fmt.Errorf("preset JSON: %w", err)
	}
	return clonePreset(preset), nil
}

func CanonicalJSON(preset Preset) ([]byte, error) {
	if err := Validate(preset); err != nil {
		return nil, fmt.Errorf("canonical preset JSON: %w", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	writeJSONField(&buffer, 1, "$schema", preset.Schema, true)
	writeJSONIntField(&buffer, 1, "schemaVersion", preset.SchemaVersion, true)
	writeJSONField(&buffer, 1, "style", string(preset.Style), true)
	writeJSONField(&buffer, 1, "radius", preset.Radius, true)
	writeJSONField(&buffer, 1, "fontSans", preset.FontSans, true)
	writeJSONField(&buffer, 1, "fontHeading", preset.FontHeading, true)
	buffer.WriteString("  \"theme\": {\n")
	writeThemeValues(&buffer, "light", preset.Theme.Light, true)
	writeThemeValues(&buffer, "dark", preset.Theme.Dark, false)
	buffer.WriteString("  }\n")
	buffer.WriteString("}\n")
	return buffer.Bytes(), nil
}

func decodePreset(decoder *json.Decoder) (Preset, error) {
	if err := expectObjectStart(decoder, "preset"); err != nil {
		return Preset{}, err
	}

	var preset Preset
	seen := make(map[string]bool, 5)
	for decoder.More() {
		key, err := objectKey(decoder, "preset")
		if err != nil {
			return Preset{}, err
		}
		if seen[key] {
			return Preset{}, fmt.Errorf("%s: duplicate key", key)
		}
		seen[key] = true

		switch key {
		case "$schema":
			preset.Schema, err = decodeString(decoder, "$schema")
		case "schemaVersion":
			preset.SchemaVersion, err = decodeSchemaVersion(decoder)
		case "style":
			var style string
			style, err = decodeString(decoder, "style")
			preset.Style = Style(style)
		case "radius":
			preset.Radius, err = decodeString(decoder, "radius")
		case "fontSans":
			preset.FontSans, err = decodeString(decoder, "fontSans")
		case "fontHeading":
			preset.FontHeading, err = decodeString(decoder, "fontHeading")
		case "theme":
			preset.Theme, err = decodeTheme(decoder)
		default:
			return Preset{}, fmt.Errorf("%s: unknown key", key)
		}
		if err != nil {
			return Preset{}, err
		}
	}
	if err := expectObjectEnd(decoder, "preset"); err != nil {
		return Preset{}, err
	}

	for _, required := range []string{"$schema", "schemaVersion", "style", "radius", "fontSans", "fontHeading", "theme"} {
		if !seen[required] {
			return Preset{}, fmt.Errorf("%s: missing required field", required)
		}
	}
	return preset, nil
}

func decodeTheme(decoder *json.Decoder) (Theme, error) {
	if err := expectObjectStart(decoder, "theme"); err != nil {
		return Theme{}, err
	}

	var theme Theme
	seen := make(map[string]bool, 2)
	for decoder.More() {
		key, err := objectKey(decoder, "theme")
		if err != nil {
			return Theme{}, err
		}
		path := "theme." + key
		if seen[key] {
			return Theme{}, fmt.Errorf("%s: duplicate key", path)
		}
		seen[key] = true

		switch key {
		case "light":
			theme.Light, err = decodeThemeValues(decoder, path)
		case "dark":
			theme.Dark, err = decodeThemeValues(decoder, path)
		default:
			return Theme{}, fmt.Errorf("%s: unknown key", path)
		}
		if err != nil {
			return Theme{}, err
		}
	}
	if err := expectObjectEnd(decoder, "theme"); err != nil {
		return Theme{}, err
	}
	for _, required := range []string{"light", "dark"} {
		if !seen[required] {
			return Theme{}, fmt.Errorf("theme.%s: missing required field", required)
		}
	}
	return theme, nil
}

func decodeThemeValues(decoder *json.Decoder, path string) (ThemeValues, error) {
	if err := expectObjectStart(decoder, path); err != nil {
		return nil, err
	}

	values := make(ThemeValues, len(tokenDefinitions))
	seen := make(map[string]bool, len(tokenDefinitions))
	for decoder.More() {
		key, err := objectKey(decoder, path)
		if err != nil {
			return nil, err
		}
		fieldPath := path + "." + key
		if seen[key] {
			return nil, fmt.Errorf("%s: duplicate key", fieldPath)
		}
		seen[key] = true
		if !isTokenName(key) {
			return nil, fmt.Errorf("%s: unknown key", fieldPath)
		}
		value, err := decodeString(decoder, fieldPath)
		if err != nil {
			return nil, err
		}
		values[key] = value
	}
	if err := expectObjectEnd(decoder, path); err != nil {
		return nil, err
	}
	return values, nil
}

func decodeSchemaVersion(decoder *json.Decoder) (int, error) {
	token, err := decoder.Token()
	if err != nil {
		return 0, fmt.Errorf("schemaVersion: %w", err)
	}
	number, ok := token.(json.Number)
	if !ok {
		return 0, fmt.Errorf("schemaVersion: expected integer, got %T", token)
	}
	version, err := strconv.Atoi(number.String())
	if err != nil {
		return 0, fmt.Errorf("schemaVersion: expected integer")
	}
	if version != SchemaVersion {
		return 0, fmt.Errorf("schemaVersion: unsupported version %d", version)
	}
	return version, nil
}

func decodeString(decoder *json.Decoder, path string) (string, error) {
	token, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	value, ok := token.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", path, token)
	}
	return value, nil
}

func expectObjectStart(decoder *json.Decoder, path string) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '{' {
		return fmt.Errorf("%s: expected object", path)
	}
	return nil
}

func expectObjectEnd(decoder *json.Decoder, path string) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '}' {
		return fmt.Errorf("%s: expected object end", path)
	}
	return nil
}

func objectKey(decoder *json.Decoder, path string) (string, error) {
	token, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	key, ok := token.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected object key", path)
	}
	return key, nil
}

func writeJSONField(buffer *bytes.Buffer, indent int, name, value string, comma bool) {
	writeIndent(buffer, indent)
	writeJSONString(buffer, name)
	buffer.WriteString(": ")
	writeJSONString(buffer, value)
	writeLineEnd(buffer, comma)
}

func writeJSONIntField(buffer *bytes.Buffer, indent int, name string, value int, comma bool) {
	writeIndent(buffer, indent)
	writeJSONString(buffer, name)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.Itoa(value))
	writeLineEnd(buffer, comma)
}

func writeThemeValues(buffer *bytes.Buffer, mode string, values ThemeValues, comma bool) {
	writeIndent(buffer, 2)
	writeJSONString(buffer, mode)
	buffer.WriteString(": {\n")
	for i, name := range TokenNames() {
		writeJSONField(buffer, 3, name, values[name], i != len(tokenDefinitions)-1)
	}
	writeIndent(buffer, 2)
	buffer.WriteByte('}')
	writeLineEnd(buffer, comma)
}

func writeJSONString(buffer *bytes.Buffer, value string) {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	buffer.Write(encoded)
}

func writeIndent(buffer *bytes.Buffer, indent int) {
	for range indent {
		buffer.WriteString("  ")
	}
}

func writeLineEnd(buffer *bytes.Buffer, comma bool) {
	if comma {
		buffer.WriteByte(',')
	}
	buffer.WriteByte('\n')
}

func clonePreset(preset Preset) Preset {
	preset.Theme.Light = cloneThemeValues(preset.Theme.Light)
	preset.Theme.Dark = cloneThemeValues(preset.Theme.Dark)
	return preset
}

func cloneThemeValues(values ThemeValues) ThemeValues {
	cloned := make(ThemeValues, len(values))
	maps.Copy(cloned, values)
	return cloned
}
