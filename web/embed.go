package web

import "embed"

// Assets contains everything under web/assets, laid out as
// assets/<library>/<version>/<file>. Paths are served verbatim, so the
// version is part of the URL and a bump is always an explicit change here.
//
//go:embed assets
var Assets embed.FS
