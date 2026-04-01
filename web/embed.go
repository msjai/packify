package web

import "embed"

// FS embeds all static files from the web directory.
//
//go:embed *.html *.js
var FS embed.FS