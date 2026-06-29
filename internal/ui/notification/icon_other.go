//go:build !darwin

package notification

import (
	_ "embed"
)

//go:embed glash-icon-solo.png
var Icon []byte
