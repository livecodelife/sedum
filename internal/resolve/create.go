package resolve

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livecodelife/sedum/internal/render"
	"github.com/livecodelife/sedum/internal/runlog"
)

// Options controls Phase 3.
type Options struct {
	// Variables are the run's values for the project facts packages declare.
	// They reach every template the same way, so a value behaves identically
	// wherever an author writes it (prov-2026-6fc3d13d).
	Variables map[string]string

	// Output is the directory authorized paths are created under.
	Output string
	// DryRun runs every decision and writes nothing.
	DryRun bool
	// Log is the run log. A nil log discards.
	Log *runlog.Log
}

// File is one authorized path after Phase 3.
type File struct {
	Resolution

	// Rendered is the template output: what the file was created with, or
	// what it would have been created with had it not already existed. Empty
	// when no template applied.
	Rendered string
	// Existed reports that the path was already on disk and was left as it
	// was.
	Existed bool
}

// Create runs Phase 3: it renders each resolution's file template and writes
// the file, unless the file is already there.
//
// Nothing is created that a provenance record did not authorize. There is no
// sibling expansion and no inference of companion files: if a record names an
// implementation file but omits its header, the header is not created and the
// injection targeting it fails loudly later.
func Create(resolutions []Resolution, opts Options) ([]File, error) {
	log := opts.Log
	if log == nil {
		log = runlog.Discard()
	}

	var (
		out      []File
		problems []error
	)
	for _, res := range resolutions {
		file, err := createOne(res, opts, log)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		out = append(out, file)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return out, nil
}

func createOne(res Resolution, opts Options, log *runlog.Log) (File, error) {
	file := File{Resolution: res}

	// An unmanaged path is carried through rather than dropped, so that a
	// later phase can say the path was declared unmanaged instead of only
	// that it was never created. The two have different fixes.
	if res.Unmanaged {
		log.Info("left unmanaged",
			"path", res.Path, "declared_by", res.UnmanagedBy, "entry", res.UnmanagedAs)
		return file, nil
	}

	rendered, err := Render(res, opts.Variables)
	if err != nil {
		return File{}, err
	}
	file.Rendered = rendered

	full := filepath.Join(opts.Output, filepath.FromSlash(res.Path))

	info, err := os.Stat(full)
	switch {
	case err == nil && info.IsDir():
		return File{}, fmt.Errorf(
			"authorized path %s is a directory in the output tree; Sedum cannot create a file there", res.Path)

	case err == nil:
		// The file is never re-rendered over. Whatever has been injected into
		// it since it was generated survives, which is what makes stopping and
		// resuming a run an ordinary workflow. What may happen is that the
		// lines its template adds are written in (prov-2026-4c49ca46).
		file.Existed = true
		if err := reconcileExisting(res, full, rendered, opts, log); err != nil {
			return File{}, err
		}
		return file, nil

	case !errors.Is(err, fs.ErrNotExist):
		return File{}, fmt.Errorf("stat %s: %w", full, err)
	}

	if res.Template == "" {
		// Not an error: a migration or a plain config file may legitimately
		// start blank, and requiring a template for every extension would be
		// noise. It is logged so that a file with no boilerplate is
		// explicable.
		log.Info("no file template matched and the package ships no default for this extension; creating the file empty",
			"path", res.Path, "package", res.Package.Name)
	} else {
		log.Info("creating file from template",
			"path", res.Path, "package", res.Package.Name, "template", res.Template,
			"default", res.Default, "captures", res.Captures)
	}

	if opts.DryRun {
		return file, nil
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return File{}, fmt.Errorf("create directory for %s: %w", res.Path, err)
	}
	if err := os.WriteFile(full, []byte(rendered), 0o644); err != nil {
		return File{}, fmt.Errorf("write %s: %w", res.Path, err)
	}
	return file, nil
}

// Render renders the file template a path matched, against the captures its
// pattern bound. It returns the empty string for a path that matched nothing.
//
// This is the half of Phase 3 that touches no filesystem, and it is exported so
// that a caller wanting to see what would be written does not have to run the
// half that does. `sedum resolve` is that caller: it reads the generators
// directory and the records directory and nothing else (prov-2026-43808a47).
func Render(res Resolution, variables map[string]string) (string, error) {
	if res.Template == "" {
		return "", nil
	}

	source, ok := res.Package.FileTemplate(res.Template)
	if !ok {
		return "", fmt.Errorf("package %s: file template %s vanished between load and generation",
			res.Package.Name, res.Template)
	}

	tmpl, err := render.Compile(res.Package.Engine, source)
	if err != nil {
		return "", fmt.Errorf("package %s: file template %s: %w", res.Package.Name, res.Template, err)
	}

	// Variables first, captures over them. They cannot actually collide - a
	// package declaring a variable named like one of its own captures is a load
	// error - so the order is a statement of which is more specific rather than
	// a precedence rule anything relies on (prov-2026-6fc3d13d).
	values := make(map[string]any, len(res.Captures)+len(variables))
	for name, value := range variables {
		values[name] = value
	}
	for name, value := range res.Captures {
		values[name] = value
	}

	out, err := tmpl.Render(values)
	if err != nil {
		return "", fmt.Errorf("path %s: file template %s: %w", res.Path, res.Template, err)
	}
	return out, nil
}

// reconcileExisting brings a file that is already on disk up to its template,
// or halts.
//
// A path with no matching template is left alone. There is nothing to reconcile
// against, and a package that ships no boilerplate for an extension is a
// supported shape rather than a problem - the same reasoning that makes an
// unmatched path create an empty file rather than an error.
func reconcileExisting(res Resolution, full, rendered string, opts Options, log *runlog.Log) error {
	if res.Template == "" || rendered == "" {
		log.Info("file already exists and its package ships no template for this path, so it is left as it is",
			"path", res.Path, "package", res.Package.Name)
		return nil
	}

	data, err := os.ReadFile(full)
	if err != nil {
		return fmt.Errorf("read %s: %w", res.Path, err)
	}

	result, err := reconcile(res.Package.CommentPrefix, res.Path, res.Template, string(data), rendered)
	if err != nil {
		return err
	}
	if result.Inserted == 0 {
		log.Info("file already exists and carries everything its template declares, left as it is",
			"path", res.Path, "package", res.Package.Name, "template", res.Template)
		return nil
	}

	// Reported rather than done quietly. Applying a template can add lines a
	// team did not write, and a team adopting Sedum into a repository has to be
	// able to read what it did to their files (prov-2026-4c49ca46).
	log.Info("file already exists; wrote in the lines its template adds",
		"path", res.Path, "package", res.Package.Name, "template", res.Template,
		"lines", result.Inserted)

	if opts.DryRun {
		return nil
	}
	if err := os.WriteFile(full, []byte(result.Content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", res.Path, err)
	}
	return nil
}
