// Package syntax is the closed registry of target-language checks
// (prov-2026-5403d162).
//
// Sedum's pipeline contains no language-specific knowledge. Where a target
// language has to be understood - generated Go must parse and be gofmt-clean,
// because model-written code must never land broken - the knowledge lives
// here and only here: named, opted into by a package's `syntax:` key, shipped
// with Sedum, and running no external program. The pipeline asks the registry
// by name and never branches on a language itself.
//
// The registry is closed because the alternative, a package-declared command,
// lets a generator package run an executable on every grow, and nothing from
// a copied or shared source may run automatically.
package syntax

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// File is one whole file to check, with the other files beside it in its
// directory that share its extension. A checker may read siblings to learn
// what the file's package declares elsewhere; it never resolves further.
type File struct {
	Path     string
	Content  string
	Siblings []string
}

// Checker checks one whole file of its language.
type Checker interface {
	// Check returns nil when the file is acceptable, and otherwise a *Problem
	// naming its path.
	Check(f File) error
}

// Siblings returns the contents of the files beside target in its directory
// under output that share its extension, target excluded. A file the run holds
// in memory is read from overrides, keyed by the same slash-separated path
// relative to output, in preference to disk.
func Siblings(output, target string, overrides map[string]string) []string {
	dir := path.Dir(target)
	ext := path.Ext(target)
	seen := map[string]bool{target: true}
	var out []string
	for p, content := range overrides {
		if path.Dir(p) == dir && path.Ext(p) == ext && !seen[p] {
			seen[p] = true
			out = append(out, content)
		}
	}
	entries, err := os.ReadDir(filepath.Join(output, filepath.FromSlash(dir)))
	if err != nil {
		return out
	}
	for _, e := range entries {
		p := path.Join(dir, e.Name())
		if e.IsDir() || path.Ext(p) != ext || seen[p] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(output, filepath.FromSlash(p)))
		if err == nil {
			seen[p] = true
			out = append(out, string(data))
		}
	}
	return out
}

// Problem is a file a checker rejected.
type Problem struct {
	Path   string
	Reason string
	// Formatted is the file as the language's formatter would write it, when
	// the file parses but is not in that form. Empty when it does not parse.
	// It is what a re-prompt shows the model, so it can bind exactly that.
	Formatted string
}

func (p *Problem) Error() string {
	return fmt.Sprintf("%s: %s", p.Path, p.Reason)
}

var registry = map[string]Checker{
	"go": goChecker{},
}

// Lookup returns the checker registered under name.
func Lookup(name string) (Checker, bool) {
	c, ok := registry[name]
	return c, ok
}

// Names returns every registered syntax, sorted.
func Names() []string {
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// goChecker accepts Go source that parses and is exactly what gofmt would
// write. It uses go/format, the library gofmt itself is built on.
type goChecker struct{}

func (goChecker) Check(f File) error {
	formatted, err := format.Source([]byte(f.Content))
	if err != nil {
		return &Problem{Path: f.Path, Reason: "does not parse as Go: " + err.Error()}
	}
	if string(formatted) != f.Content {
		return &Problem{Path: f.Path, Reason: "is not gofmt-clean", Formatted: string(formatted)}
	}
	if missing := unimportedQualifiers(f); len(missing) > 0 {
		return &Problem{Path: f.Path, Reason: fmt.Sprintf(
			"uses %s as a package qualifier without importing it; import the package each one names",
			strings.Join(missing, ", "))}
	}
	if locals, imports := unused(f); len(locals)+len(imports) > 0 {
		var parts []string
		if len(locals) > 0 {
			parts = append(parts, "declares and never uses "+strings.Join(locals, ", "))
		}
		if len(imports) > 0 {
			parts = append(parts, "imports and never uses "+strings.Join(imports, ", "))
		}
		return &Problem{Path: f.Path, Reason: strings.Join(parts, "; ") +
			"; Go does not compile either, so remove each one or use it"}
	}
	return nil
}

// unused returns the function-local variables the file declares and never
// uses, and the imports nothing in it uses, each as "name (line N)"
// (prov-2026-ed9124e1). Both are compile errors in Go and both are visible in
// the file alone, so go/parser's per-file resolution is enough; like the
// qualifier check it errs toward silence, so a correct answer is never
// re-prompted.
//
// A local is used when any reference resolves to it other than its own
// declaration and the bare target of a plain assignment: x = 1 is not a use,
// while x += 1, x++, s.f = 1 and m[k] = v are, as the compiler counts them.
// Parameters, named results, the blank identifier and package-level names are
// never candidates.
func unused(f File) (locals, imports []string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, f.Path, f.Content, 0)
	if err != nil {
		return nil, nil
	}
	at := func(id *ast.Ident) string {
		return fmt.Sprintf("%s (line %d)", id.Name, fset.Position(id.Pos()).Line)
	}

	declared := map[*ast.Object]*ast.Ident{}
	declaring := map[*ast.Ident]bool{}
	candidate := func(id *ast.Ident, decl any) {
		if id == nil || id.Name == "_" || id.Obj == nil || id.Obj.Kind != ast.Var {
			return
		}
		if decl != nil && id.Obj.Decl != decl {
			return // a redeclaration in :=, not a new variable
		}
		declared[id.Obj] = id
		declaring[id] = true
	}
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			switch s := n.(type) {
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE {
					for _, lhs := range s.Lhs {
						id, _ := lhs.(*ast.Ident)
						candidate(id, s)
					}
				}
			case *ast.RangeStmt:
				if s.Tok == token.DEFINE {
					for _, e := range []ast.Expr{s.Key, s.Value} {
						if id, ok := e.(*ast.Ident); ok {
							candidate(id, nil)
						}
					}
				}
			case *ast.ValueSpec:
				for _, id := range s.Names {
					candidate(id, s)
				}
			}
			return true
		})
	}

	assigned := map[*ast.Ident]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if s, ok := n.(*ast.AssignStmt); ok && s.Tok == token.ASSIGN {
			for _, lhs := range s.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					assigned[id] = true
				}
			}
		}
		return true
	})
	used := map[*ast.Object]bool{}
	qualifiers := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj != nil && !declaring[x] && !assigned[x] {
				used[x.Obj] = true
			}
		case *ast.SelectorExpr:
			if id, ok := x.X.(*ast.Ident); ok && id.Obj == nil {
				qualifiers[id.Name] = true
			}
		}
		return true
	})
	var unusedLocals []*ast.Ident
	for obj, id := range declared {
		if !used[obj] {
			unusedLocals = append(unusedLocals, id)
		}
	}
	sort.Slice(unusedLocals, func(i, j int) bool { return unusedLocals[i].Pos() < unusedLocals[j].Pos() })
	for _, id := range unusedLocals {
		locals = append(locals, at(id))
	}

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "`\"")
		var names []string
		switch {
		case imp.Name != nil && (imp.Name.Name == "_" || imp.Name.Name == "."):
			continue
		case imp.Name != nil:
			names = []string{imp.Name.Name}
		case path == "C":
			continue // cgo
		default:
			names = importNames(path)
		}
		usedImport, inferable := false, false
		for _, name := range names {
			if token.IsIdentifier(name) {
				inferable = true
			}
			if qualifiers[name] {
				usedImport = true
			}
		}
		if !usedImport && inferable {
			id := &ast.Ident{Name: names[0], NamePos: imp.Pos()}
			imports = append(imports, at(id))
		}
	}
	return locals, imports
}

// unimportedQualifiers returns, sorted, the identifiers used as X in an X.Sel
// selector that nothing in reach supplies (prov-2026-0ce03d1c): the parser
// leaves them unresolved in the file, no import of the file can supply the
// name, and no sibling file declares it at top level. It uses go/parser only
// and errs toward silence - a false alarm re-prompts a correct answer until its
// retries run out - so a dot import turns the check off and an import whose
// name cannot be inferred still supplies every name it plausibly could.
func unimportedQualifiers(f File) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, f.Path, f.Content, 0)
	if err != nil {
		return nil
	}

	supplied := map[string]bool{}
	for _, imp := range file.Imports {
		if imp.Name != nil {
			if imp.Name.Name == "." {
				return nil
			}
			supplied[imp.Name.Name] = true
			continue
		}
		for _, name := range importNames(strings.Trim(imp.Path.Value, "`\"")) {
			supplied[name] = true
		}
	}
	for _, sibling := range f.Siblings {
		for name := range topLevelNames(sibling) {
			supplied[name] = true
		}
	}

	unresolved := map[*ast.Ident]bool{}
	for _, id := range file.Unresolved {
		unresolved[id] = true
	}
	missing := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if x, ok := sel.X.(*ast.Ident); ok && unresolved[x] && !supplied[x.Name] {
			missing[x.Name] = true
		}
		return true
	})

	out := make([]string, 0, len(missing))
	for name := range missing {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

var majorVersion = regexp.MustCompile(`^v[0-9]+$`)
var dotVersion = regexp.MustCompile(`\.v[0-9]+$`)

// importNames are the package names an import path plausibly gives, without
// type information: the last element, and that element without a .vN suffix,
// a -go suffix or a go- prefix; and for a path ending in a /vN major-version
// element, the element before it, treated the same way.
func importNames(importPath string) []string {
	elems := strings.Split(importPath, "/")
	last := elems[len(elems)-1]
	candidates := []string{last}
	if majorVersion.MatchString(last) && len(elems) > 1 {
		candidates = append(candidates, elems[len(elems)-2])
	}
	var out []string
	for _, c := range candidates {
		c = dotVersion.ReplaceAllString(c, "")
		out = append(out, c, strings.TrimPrefix(c, "go-"), strings.TrimSuffix(c, "-go"))
	}
	return out
}

// topLevelNames are the names a Go file declares at package level.
func topLevelNames(content string) map[string]bool {
	out := map[string]bool{}
	file, err := parser.ParseFile(token.NewFileSet(), "", content, parser.SkipObjectResolution)
	if err != nil {
		return out
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil {
				out[d.Name.Name] = true
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for _, n := range s.Names {
						out[n.Name] = true
					}
				case *ast.TypeSpec:
					out[s.Name.Name] = true
				}
			}
		}
	}
	return out
}
