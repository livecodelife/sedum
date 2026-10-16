# Context for growing package adopt (advisory; not part of gen_hash)

Module: github.com/livecodelife/sedum. Package `adopt` (internal/adopt) is the
language-agnostic core of `sedum adopt`: it brings existing files under Sedum
without a model. It contains NO language knowledge: no go/parser, no gofmt, no
regex that mentions a language. Standard library only, plus the internal
packages listed below. Return errors wrapped with fmt.Errorf("...: %w", err).
Never panic. Every exported identifier has a doc comment. All code must be
gofmt-clean. Tests are table-driven, use t.TempDir(), and never run git.

## Types that every file in the package shares (defined in span.go)

    type Kind string
    const (
        KindHeader   Kind = "header"
        KindUnit     Kind = "unit"
        KindResidual Kind = "residual"
    )
    // Span is one contiguous piece of a source file, in file order.
    type Span struct {
        Kind   Kind           `json:"kind"`
        Text   string         `json:"text"`
        Action string         `json:"action,omitempty"` // KindUnit only
        Kwargs map[string]any `json:"kwargs,omitempty"` // KindUnit only; the action's kwargs
    }
    // Adopter splits one file into spans.
    type Adopter interface {
        Split(path string, src []byte) ([]Span, error)
    }

## Internal packages and their exact signatures

package record (github.com/livecodelife/sedum/internal/record):
    func Load(dir string, opts Options) (*Set, []string, error)   // Options is a struct; zero value is fine
    type Set struct { Records []*Record; Known []*Record }
    func (s *Set) Lookup(id string) (*Record, bool)
    type Record struct { ID, Source, Intent string; Constraints, Paths, Patterns, Forbidden []string;
        Gen, CreatedAt string; Aliases []string; Deprecated bool; Type, Implements string }
    func (r *Record) Authorizes(target string) bool
    func (r *Record) Forbids(target string) (string, bool)

package genpkg (internal/genpkg):
    func Load(dir string, opts Options) (*Set, []Finding, error)  // Options zero value is fine
    type Set struct { Packages []*Package }
    func (s *Set) ForExtension(ext string) []*Package   // ext includes the dot, e.g. ".go"
    func (s *Set) Lookup(name string) (*Package, bool)
    func (s *Set) Unmanaged(target string) (pkg, entry string, ok bool)
    type Package struct { Name, Dir string; Extensions []string; CommentPrefix, CommentSuffix string;
        Unmanaged, TestPaths []string; Syntax string; ... }
    type Comment struct { Prefix, Suffix string }  // methods: Wrap(body string) string; Unwrap(text string) string; Delimited() bool

package pathpat (internal/pathpat): path-pattern grammar for scope entries (** * ? [..]).
    check the package for its exported Match function before use

package recording (internal/recording):
    type Recording struct { SedumVersion string `json:"sedum_version"`; Packages map[string]Package `json:"packages"`;
        Records []Record `json:"records"`; Variables map[string]string `json:"variables,omitempty"` ... }
    type Package struct { Extensions []string `json:"extensions"` }
    type Record struct { RecordID string; Gen, CreatedAt string; Files []File; Phases []Phase }
    type File struct { Path, Package, Template string; Captures map[string]string }
    type Phase struct { Name string; Invocations []Invocation }
    type Invocation struct { Action string; Kwargs map[string]any }
    const DefaultPhase = "default"; const Version = "0.8.0"
    func Write(w io.Writer, r Recording) error

package resolve (internal/resolve):
    func Paths(set *genpkg.Set, records *record.Set, prefer []string) ([]Resolution, []string, error)
    type Resolution struct { RecordID, Path string; Package *genpkg.Package; Template string;
        Captures map[string]string; Unmanaged bool; UnmanagedBy, UnmanagedAs string; Default bool }

package pipeline (internal/pipeline):
    type ReplayConfig struct { Generators, Output, Records string; Variables map[string]string;
        RestampGen, DryRun bool; StopAfterPhase int; Log *runlog.Log }
    func Replay(rec recording.Recording, cfg ReplayConfig) (*Result, error)   // writes files under cfg.Output
    type Result struct { Files []resolve.File; Warnings []string; ... }

package inject (internal/inject):
    func FindRegions(comment genpkg.Comment, content string) ([]Region, error)
    type Region struct { Marker Marker; Start, End int }  // content[Start:End] is the opening marker line, the body, and the closing marker line

package runlog (internal/runlog): func Discard() *Log

## Behaviour contract shared by the stage (blueprint prov-2026-95d9ae7e)

- Coverage: span texts concatenated in order equal the file byte for byte.
- v1 policy: a file with a residual span that is not all whitespace is skipped, never partly written.
- No model call anywhere. No external program except a command adopter the user allowed.
