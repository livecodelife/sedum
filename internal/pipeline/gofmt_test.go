package pipeline

import (
	"go/format"
	"path/filepath"
	"testing"

	"github.com/livecodelife/sedum/internal/recording"
)

// Generated Go passes gofmt with no diff (prov-2026-ab24df5a). The chi fixture
// is the Go reference package, so its output is checked with go/format: an
// import, two handlers sharing an anchor, and their routes, which together
// cover the one layout gofmt rewrites - a top-level declaration followed by
// its closing marker. The check lives in a test, never in core: Sedum knows
// no target language.
func TestChiOutputIsGofmtClean(t *testing.T) {
	handler := func(verb string) recording.Invocation {
		return recording.Invocation{Action: "createHandler", Kwargs: map[string]any{"resource": "users", "verb": verb}}
	}
	rec := recording.Recording{
		SedumVersion: recording.Version,
		Packages:     map[string]recording.Package{"chi": {Extensions: []string{".go"}}},
		Records: []recording.Record{{
			RecordID: "prov-2026-aaaaaaaa",
			Files: []recording.File{{
				Path: "internal/handlers/users.go", Package: "chi",
				Template: "internal/handlers/{name}.go", Captures: map[string]string{"name": "users"},
			}},
			Phases: []recording.Phase{{Name: recording.DefaultPhase, Invocations: []recording.Invocation{
				{Action: "addImport", Kwargs: map[string]any{"resource": "users", "path": "net/http"}},
				handler("list"),
				handler("show"),
			}}},
		}},
	}

	output := t.TempDir()
	if _, err := Replay(rec, ReplayConfig{
		Generators: filepath.Join("..", "..", "testdata", "generators"),
		Output:     output,
	}); err != nil {
		t.Fatalf("Replay: %v", err)
	}

	got := readFile(t, filepath.Join(output, "internal", "handlers", "users.go"))
	formatted, err := format.Source([]byte(got))
	if err != nil {
		t.Fatalf("generated Go does not parse: %v\n%s", err, got)
	}
	if string(formatted) != got {
		t.Errorf("generated Go is not gofmt-clean:\n got %q\nwant %q", got, formatted)
	}
}
