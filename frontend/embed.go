package frontend

import (
	"embed"
	"io/fs"
)

//go:embed *.php assets/* include/*.php templates/*.tpl
var embedded embed.FS

// Files contains the PHP dashboard, templates, and static assets.
var Files fs.FS = embedded
