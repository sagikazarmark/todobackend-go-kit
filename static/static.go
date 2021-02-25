package static

import "embed"

//go:embed index.html
var files embed.FS

// Files returns a filesystem with static files.
func Files() embed.FS {
	return files
}
