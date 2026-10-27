package scaffold

import (
	"os"
	"path/filepath"
)

// InstallHook writes the hook text as the file name in the hooks directory,
// creating the directory when needed. A hook already holding exactly the text is
// left as it is ("unchanged"); a different one is moved to <name>.bak when none
// exists yet ("replaced (previous hook kept as <name>.bak)") and otherwise
// overwritten ("replaced"). The file ends up mode 0755 either way.
func InstallHook(hooksDir, name, text string) (string, error) {
	path := filepath.Join(hooksDir, name)
	result := "installed"
	existing, err := os.ReadFile(path)
	switch {
	case err == nil:
		if string(existing) == text {
			if err := os.Chmod(path, 0o755); err != nil {
				return "", err
			}
			return "unchanged", nil
		}
		backup := path + ".bak"
		if _, err := os.Stat(backup); os.IsNotExist(err) {
			if err := os.Rename(path, backup); err != nil {
				return "", err
			}
			result = "replaced (previous hook kept as " + name + ".bak)"
		} else if err != nil {
			return "", err
		} else {
			result = "replaced"
		}
	case !os.IsNotExist(err):
		return "", err
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(text), 0o755); err != nil {
		return "", err
	}
	if err := os.Chmod(path, 0o755); err != nil {
		return "", err
	}
	return result, nil
}
