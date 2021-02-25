package static

import "embed"

// nolint: gochecknoglobals
//go:embed index.html
var files embed.FS

// Files returns a filesystem with static files.
func Files() embed.FS {
	return files
}
