// Package merge holds the Tailwind-aware class merger named by gsx.toml's
// class_merger. Components never call this directly — gsx invokes it wherever
// class values merge (duplicate class attrs, fallthrough attrs bags).
package merge

import twmerge "github.com/jackielii/tailwind-merge-go/pkg/twmerge"

var merger = twmerge.CreateTwMerge(twmerge.GetDefaultConfig())

// Merge is gsx's canonical func([]string) string merger seam.
func Merge(classes []string) string { return merger(classes) }
