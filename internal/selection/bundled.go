package selection

import (
	"os"
	"path/filepath"
)

// Bundled asset filenames, sibling to Sedum's own executable in a release
// archive. hack/fetch-release-assets.sh and .goreleaser.yml are what put
// them there; this is only where Sedum looks for them.
//
// The model's own filename is Sedum's, not the checkpoint's - a re-tuned
// default ships under the same name, so nothing here has to change with it.
const (
	bundledServerName = "goinfer-serve"
	bundledModelName  = "sedum-default-model.gguf"
)

// BundledLocalConfig looks for a goinfer-serve binary and a gguf next to
// Sedum's own executable, and reports whether both are present.
//
// go install and a from-source build carry neither, so ok is false there,
// deliberately: a bundled default is a fact about a release archive, not
// something this code can conjure without one.
func BundledLocalConfig() (cfg LocalConfig, ok bool) {
	exe, err := os.Executable()
	if err != nil {
		return LocalConfig{}, false
	}
	dir := filepath.Dir(exe)

	server := filepath.Join(dir, bundledServerName)
	model := filepath.Join(dir, bundledModelName)
	if !isFile(server) || !isFile(model) {
		return LocalConfig{}, false
	}
	return LocalConfig{ServerPath: server, ModelPath: model}, true
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
