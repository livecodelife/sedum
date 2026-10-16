package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare.go","path":"io/fs"}}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/prepare.go","path":"gopkg.in/yaml.v3"}}
import "gopkg.in/yaml.v3"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// growPrompt is the same advisory text for every variant, so the variants\n// differ only in the intent.\nconst growPrompt = \"# Context\\n\\nA small Go package. Standard library only (plus what the context files import). Write each target file in full, gofmt-clean, with a doc comment on every exported identifier.\\n\"","file":"evals/design/prepare.go","name":"growPrompt"}}
// growPrompt is the same advisory text for every variant, so the variants
// differ only in the intent.
const growPrompt = "# Context\n\nA small Go package. Standard library only (plus what the context files import). Write each target file in full, gofmt-clean, with a doc comment on every exported identifier.\n"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// record is the v1 record a sample's intent is written into.\ntype record struct {\n\tID            string   `yaml:\"id\"`\n\tTitle         string   `yaml:\"title\"`\n\tStatus        string   `yaml:\"status\"`\n\tCreatedAt     string   `yaml:\"created_at\"`\n\tAuthor        string   `yaml:\"author\"`\n\tIntent        string   `yaml:\"intent\"`\n\tConstraints   []string `yaml:\"constraints\"`\n\tAffectedScope []string `yaml:\"affected_scope\"`\n\tForbidden     []string `yaml:\"forbidden_scope\"`\n\tType          string   `yaml:\"type\"`\n}","file":"evals/design/prepare.go","name":"record"}}
// record is the v1 record a sample's intent is written into.
type record struct {
	ID            string   `yaml:"id"`
	Title         string   `yaml:"title"`
	Status        string   `yaml:"status"`
	CreatedAt     string   `yaml:"created_at"`
	Author        string   `yaml:"author"`
	Intent        string   `yaml:"intent"`
	Constraints   []string `yaml:"constraints"`
	AffectedScope []string `yaml:"affected_scope"`
	Forbidden     []string `yaml:"forbidden_scope"`
	Type          string   `yaml:"type"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// copyTree copies the directory src into dst.\nfunc copyTree(src, dst string) error {\n\treturn filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\trel, err := filepath.Rel(src, p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\ttarget := filepath.Join(dst, rel)\n\t\tif d.IsDir() {\n\t\t\treturn os.MkdirAll(target, 0o755)\n\t\t}\n\t\tdata, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn os.WriteFile(target, data, 0o644)\n\t})\n}","file":"evals/design/prepare.go","name":"copyTree"}}
// copyTree copies the directory src into dst.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:b32ae71e464781793d843f280f11ab448c95c71fc8ef71ed4bb5f4c31f3e0b69","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Prepare writes the project for one sample into dir: a go.mod, the task's\n// context files, a copy of the generator packages, the grow prompt, and one\n// open record whose intent is the variant's text and whose scope is the target\n// files. The hidden tests are not written; Grade adds them after the grow.\nfunc Prepare(t Task, v Variant, generators, dir string) error {\n\tintent, err := t.Intent(v)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif err := t.Isolated(); err != nil {\n\t\treturn err\n\t}\n\tmod := fmt.Sprintf(\"module design/%s\\n\\ngo 1.25\\n\", t.ID)\n\tif err := os.WriteFile(filepath.Join(dir, \"go.mod\"), []byte(mod), 0o644); err != nil {\n\t\treturn err\n\t}\n\tfor _, f := range t.Context {\n\t\tdata, err := os.ReadFile(filepath.Join(t.Dir, \"context\", f))\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif err := os.WriteFile(filepath.Join(dir, f), data, 0o644); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\tif err := copyTree(generators, filepath.Join(dir, \"generators\")); err != nil {\n\t\treturn fmt.Errorf(\"copy generators: %w\", err)\n\t}\n\tif err := os.WriteFile(filepath.Join(dir, \"PROMPT.md\"), []byte(growPrompt), 0o644); err != nil {\n\t\treturn err\n\t}\n\trec := record{\n\t\tID:            fmt.Sprintf(\"design-%s-%s\", t.ID, v),\n\t\tTitle:         fmt.Sprintf(\"design eval: %s (%s)\", t.ID, v),\n\t\tStatus:        \"open\",\n\t\tCreatedAt:     \"2026-10-16\",\n\t\tAuthor:        \"eval@example.com\",\n\t\tIntent:        intent,\n\t\tConstraints:   []string{},\n\t\tAffectedScope: t.Target,\n\t\tForbidden:     []string{},\n\t\tType:          \"blueprint\",\n\t}\n\tdata, err := yaml.Marshal(rec)\n\tif err != nil {\n\t\treturn err\n\t}\n\tif err := os.MkdirAll(filepath.Join(dir, \"provenance\"), 0o755); err != nil {\n\t\treturn err\n\t}\n\treturn os.WriteFile(filepath.Join(dir, \"provenance\", rec.ID+\".yml\"), data, 0o644)\n}","file":"evals/design/prepare.go","name":"Prepare"}}
// Prepare writes the project for one sample into dir: a go.mod, the task's
// context files, a copy of the generator packages, the grow prompt, and one
// open record whose intent is the variant's text and whose scope is the target
// files. The hidden tests are not written; Grade adds them after the grow.
func Prepare(t Task, v Variant, generators, dir string) error {
	intent, err := t.Intent(v)
	if err != nil {
		return err
	}
	if err := t.Isolated(); err != nil {
		return err
	}
	mod := fmt.Sprintf("module design/%s\n\ngo 1.25\n", t.ID)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		return err
	}
	for _, f := range t.Context {
		data, err := os.ReadFile(filepath.Join(t.Dir, "context", f))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, f), data, 0o644); err != nil {
			return err
		}
	}
	if err := copyTree(generators, filepath.Join(dir, "generators")); err != nil {
		return fmt.Errorf("copy generators: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "PROMPT.md"), []byte(growPrompt), 0o644); err != nil {
		return err
	}
	rec := record{
		ID:            fmt.Sprintf("design-%s-%s", t.ID, v),
		Title:         fmt.Sprintf("design eval: %s (%s)", t.ID, v),
		Status:        "open",
		CreatedAt:     "2026-10-16",
		Author:        "eval@example.com",
		Intent:        intent,
		Constraints:   []string{},
		AffectedScope: t.Target,
		Forbidden:     []string{},
		Type:          "blueprint",
	}
	data, err := yaml.Marshal(rec)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "provenance"), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "provenance", rec.ID+".yml"), data, 0o644)
}

// /sedum:setFunction
