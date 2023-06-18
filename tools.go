//go:build tools
// +build tools

package troupe

import (
	_ "golang.org/x/tools/cmd/stringer"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
