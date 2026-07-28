package frontend

import (
	"io/fs"
	"testing"
)

func TestFiles(t *testing.T) {
	for _, name := range []string{"index.php", "assets/style.css", "include/Template.php", "templates/main.tpl"} {
		if _, err := fs.ReadFile(Files, name); err != nil {
			t.Errorf("read embedded file %q: %v", name, err)
		}
	}
}
