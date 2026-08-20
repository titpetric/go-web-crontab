package frontend

import (
	"embed"
	"io/fs"
)

//go:embed *.php composer.json composer.lock assets/* templates/*.tpl vendor
var embedded embed.FS

// Files contains the PHP dashboard, templates, and static assets.
var Files fs.FS = embedded
