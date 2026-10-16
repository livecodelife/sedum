# Context for sedum grow (package genpkg)

Pass with: sedum-haiku grow --generators generators --records provenance --prompt-file .sedum/context.md --only <id>
Advisory only: behaviour belongs in the record. Editing this file does not regenerate anything.

Conventions
- Package genpkg (module github.com/livecodelife/sedum). Generated Go is gofmt-clean and calls only identifiers listed here or in the record.
- Each import is its own addImport. Never alter, reorder or add anything beyond what the record lists.
- Code kwargs are the declaration verbatim: signature, body and doc comment, gofmt-formatted.

Identifiers that exist in package genpkg (call, do not redefine)
- type Package (fields Name, Actions map[string]*Action, Transforms, OpExceptions) with method ActionTemplate(path string) (string, bool)
- type Action (fields Name, Kwargs, Template, Templates, Discriminator, Variants, InjectsInto, Anchor, AnchorStart, AnchorEnd, AnchorPattern, Identity, PkgHash); type Kwarg (Type, Required, Default)
- const SoleTemplate

Package pipeline (internal/pipeline)
- Identifiers that exist in package pipeline (call, do not redefine): replayGen(rec recording.Recording, records *record.Set, recordID string) string; genOf(records *record.Set, recordID string) string; anyManaged(files []resolve.File) bool; targetsOf(files []resolve.File) []inject.Target; type Selection (fields RecordID string, Files []resolve.File, Invocations).
- Test helpers that exist in package pipeline tests: config(t) Config; recorded(t, cfg) (recording.Recording, string); replayConfig(t, cfg) ReplayConfig; writeTree(t, map[string]string) string; readFile(t, path) string; markerGen(t, output) string; type stub with a response field; the constant withField.
- inject.Result has Path, Action, Variant, Record string and Pruned, Held, Skipped, Replaced, Changed, Orphan, Unselected bool.
Package cli (internal/cli)
- Identifiers that exist in package cli (call, do not redefine): describeRef(r pipeline.RegionRef) string.
- pipeline.Orphan has fields Path, Action, Variant, Record string and Authorized, Pruned bool; pipeline.Result has Orphans []pipeline.Orphan; pipeline.RegionRef has Path, Action, Variant string.
