// Package canonical holds the structural component sources and the shapes they
// declare. It is compiled so that it type-checks and so that style-independent
// structure and behavior tests can run against the authoritative source, but it
// is never shipped: consumers receive the generated (canonical x recipe) output
// in package ui. Nothing may import it — internal/stylegen reads the shapes it
// needs from registry/canonical/shapes and reads these sources as files.
package canonical
