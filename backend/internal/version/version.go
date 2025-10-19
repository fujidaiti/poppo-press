// Package version exposes build-time version information for diagnostics and
// observability endpoints.
package version

// Version is the semantic version of the binary, set at build time.
var Version = "dev"

// Commit is the VCS commit hash the binary was built from.
var Commit = "none"

// Date is the build timestamp, typically in ISO-8601 format.
var Date = "unknown"
