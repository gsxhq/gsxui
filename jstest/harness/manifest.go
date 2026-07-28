package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/gsxhq/gsxui/site/examples"
)

// entry is one addressable example page. Component is the REGISTERED
// component name ("navigation-menu"), not the example directory name
// ("navigationmenu") — Go package names can't contain hyphens, but
// examples.For is keyed by the registered name and that is what tests read.
type entry struct {
	Component string `json:"component"`
	Example   string `json:"example"`
	URL       string `json:"url"`
}

// buildManifest enumerates every registered example in registration order.
// Named previews replace their example's base node in the browser corpus:
// each independently addressable node gets a unique test identity, while the
// base route remains available to the site and focused tests.
func buildManifest() []entry {
	var out []entry
	for _, component := range examples.Components() {
		for _, ex := range examples.For(component) {
			if len(ex.Previews) != 0 {
				for _, preview := range ex.Previews {
					query := url.Values{
						examples.PreviewQueryKey: []string{preview.Name},
					}
					out = append(out, entry{
						Component: component,
						Example:   ex.Name + "/" + preview.Name,
						URL:       "/x/" + component + "/" + ex.Name + "?" + query.Encode(),
					})
				}
				continue
			}
			out = append(out, entry{
				Component: component,
				Example:   ex.Name,
				URL:       "/x/" + component + "/" + ex.Name,
			})
		}
	}
	return out
}

// writeManifest serialises buildManifest to path. Playwright's globalSetup
// runs this before workers import spec files, so the specs can read the
// example list synchronously and generate one test per example.
func writeManifest(path string) error {
	b, err := json.MarshalIndent(buildManifest(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
