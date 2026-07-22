// Package gsxui embeds the component registry the gsxui CLI vendors from.
// The embedded tree is the single source of truth: the CLI derives the
// component list, dependencies, and JS presence from it at runtime.
package gsxui

import "embed"

//go:embed ui assets merge
var Files embed.FS
