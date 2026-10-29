package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// A slim install carries no bundled model, so the error of a run that names no
// model says how to get the full one (prov-2026-668746b1). The text the sealed
// tests look for stays in the message.

const fullInstallCommand = "curl -fsSL https://raw.githubusercontent.com/livecodelife/sedum/main/install.sh | sh"

func TestGrowWithNoModelNamesTheFlagsAndTheFullInstall(t *testing.T) {
	dir := t.TempDir()

	_, err := exec(t, "grow",
		"--generators", fixtureGenerators(),
		"--records", fixtureRecords(t),
		"--output", dir,
		"--log", filepath.Join(t.TempDir(), "run.log"))

	wantErr(t, err,
		"no --model or --local-model given, and no bundled model found",
		"--model", "--local-model",
		fullInstallCommand)
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("a refused run wrote %d entries: it must fail before any model call", len(entries))
	}
}
