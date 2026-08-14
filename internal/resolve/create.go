package resolve

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/render"
	"github.com/calebcowen/sedum/internal/runlog"
)

// Options controls Phase 3.
type Options struct {
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

	rendered, err := renderTemplate(res)
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
		// Create-if-absent. Re-rendering would destroy whatever has been
		// injected into the file since it was generated, which is what makes
		// stopping and resuming a run an ordinary workflow.
		file.Existed = true
		if err := checkMarkers(res, full, rendered); err != nil {
			return File{}, err
		}
		log.Info("file already exists, left as it is", "path", res.Path, "package", res.Package.Name)
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

// renderTemplate renders the file template a path matched, against the captures
// its pattern bound.
func renderTemplate(res Resolution) (string, error) {
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

	values := make(map[string]any, len(res.Captures))
	for name, value := range res.Captures {
		values[name] = value
	}

	out, err := tmpl.Render(values)
	if err != nil {
		return "", fmt.Errorf("path %s: file template %s: %w", res.Path, res.Template, err)
	}
	return out, nil
}

// checkMarkers verifies that a file which already exists carries the markers
// its template plants. Injection points exist because a file template created
// them, so a file missing them is one no action can be applied to.
func checkMarkers(res Resolution, full, rendered string) error {
	data, err := os.ReadFile(full)
	if err != nil {
		return fmt.Errorf("read %s: %w", res.Path, err)
	}

	missing := genpkg.MissingMarkers(res.Package.CommentPrefix, rendered, string(data))
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(
		"file %s exists but does not carry %s planted by its file template %s; either something other than Sedum wrote it, or the template changed shape after it was generated",
		res.Path, markerList(missing), res.Template)
}

func markerList(names []string) string {
	quoted := make([]string, 0, len(names))
	for _, n := range names {
		quoted = append(quoted, strconv.Quote(n))
	}
	if len(names) == 1 {
		return "marker " + quoted[0]
	}
	return "markers " + strings.Join(quoted, ", ")
}
