package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
)

// Version is the CLI version string (set by GoReleaser).
var Version = "dev"

// VersionFlag prints the version and exits.
type VersionFlag string

// Decode is a no-op for the boolean version flag.
func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }

// IsBool marks the version flag as boolean.
func (v VersionFlag) IsBool() bool { return true }

// BeforeApply prints the version and exits.
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	_, _ = fmt.Fprintln(app.Stdout, vars["version"])
	app.Exit(0)
	return nil
}
