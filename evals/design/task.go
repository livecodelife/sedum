package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/task.go","path":"gopkg.in/yaml.v3"}}
import "gopkg.in/yaml.v3"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Task is one design task: a small Go package the model must write, graded by\n// hidden tests it is never shown. Context files are copied in beside it.\ntype Task struct {\n\tID         string   `yaml:\"id\"`\n\tPackage    string   `yaml:\"package\"`\n\tTarget     []string `yaml:\"target\"`\n\tContext    []string `yaml:\"context\"`\n\tHidden     []string `yaml:\"hidden\"`\n\tBehaviour  string   `yaml:\"behaviour\"`\n\tSignatures string   `yaml:\"signatures\"`\n\tCode       string   `yaml:\"code\"`\n\n\t// Dir is the task directory it was loaded from.\n\tDir string `yaml:\"-\"`\n}","file":"evals/design/task.go","name":"Task"}}
// Task is one design task: a small Go package the model must write, graded by
// hidden tests it is never shown. Context files are copied in beside it.
type Task struct {
	ID         string   `yaml:"id"`
	Package    string   `yaml:"package"`
	Target     []string `yaml:"target"`
	Context    []string `yaml:"context"`
	Hidden     []string `yaml:"hidden"`
	Behaviour  string   `yaml:"behaviour"`
	Signatures string   `yaml:"signatures"`
	Code       string   `yaml:"code"`

	// Dir is the task directory it was loaded from.
	Dir string `yaml:"-"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// LoadTask reads dir/task.yaml and checks it. Context files live under\n// dir/context/ and hidden test files under dir/hidden/; every one named must\n// exist, hidden files must end in _test.go, and no target file may also be a\n// context or hidden file.\nfunc LoadTask(dir string) (Task, error) {\n\tdata, err := os.ReadFile(filepath.Join(dir, \"task.yaml\"))\n\tif err != nil {\n\t\treturn Task{}, fmt.Errorf(\"read task: %w\", err)\n\t}\n\tvar t Task\n\tif err := yaml.Unmarshal(data, &t); err != nil {\n\t\treturn Task{}, fmt.Errorf(\"parse %s: %w\", filepath.Join(dir, \"task.yaml\"), err)\n\t}\n\tt.Dir = dir\n\tif t.ID == \"\" || t.ID != filepath.Base(dir) {\n\t\treturn Task{}, fmt.Errorf(\"task %s: id %q must equal the directory name\", dir, t.ID)\n\t}\n\tif t.Package == \"\" || len(t.Target) == 0 || len(t.Hidden) == 0 {\n\t\treturn Task{}, fmt.Errorf(\"task %s: package, target and hidden are required\", t.ID)\n\t}\n\tif t.Behaviour == \"\" || t.Signatures == \"\" || t.Code == \"\" {\n\t\treturn Task{}, fmt.Errorf(\"task %s: behaviour, signatures and code intents are all required\", t.ID)\n\t}\n\tfor _, f := range t.Hidden {\n\t\tif !strings.HasSuffix(f, \"_test.go\") {\n\t\t\treturn Task{}, fmt.Errorf(\"task %s: hidden file %s must end in _test.go\", t.ID, f)\n\t\t}\n\t\tif _, err := os.Stat(filepath.Join(dir, \"hidden\", f)); err != nil {\n\t\t\treturn Task{}, fmt.Errorf(\"task %s: %w\", t.ID, err)\n\t\t}\n\t}\n\tfor _, f := range t.Context {\n\t\tif _, err := os.Stat(filepath.Join(dir, \"context\", f)); err != nil {\n\t\t\treturn Task{}, fmt.Errorf(\"task %s: %w\", t.ID, err)\n\t\t}\n\t}\n\tfor _, f := range t.Target {\n\t\tfor _, other := range append(append([]string{}, t.Context...), t.Hidden...) {\n\t\t\tif f == other {\n\t\t\t\treturn Task{}, fmt.Errorf(\"task %s: target %s is also a context or hidden file\", t.ID, f)\n\t\t\t}\n\t\t}\n\t}\n\treturn t, nil\n}","file":"evals/design/task.go","name":"LoadTask"}}
// LoadTask reads dir/task.yaml and checks it. Context files live under
// dir/context/ and hidden test files under dir/hidden/; every one named must
// exist, hidden files must end in _test.go, and no target file may also be a
// context or hidden file.
func LoadTask(dir string) (Task, error) {
	data, err := os.ReadFile(filepath.Join(dir, "task.yaml"))
	if err != nil {
		return Task{}, fmt.Errorf("read task: %w", err)
	}
	var t Task
	if err := yaml.Unmarshal(data, &t); err != nil {
		return Task{}, fmt.Errorf("parse %s: %w", filepath.Join(dir, "task.yaml"), err)
	}
	t.Dir = dir
	if t.ID == "" || t.ID != filepath.Base(dir) {
		return Task{}, fmt.Errorf("task %s: id %q must equal the directory name", dir, t.ID)
	}
	if t.Package == "" || len(t.Target) == 0 || len(t.Hidden) == 0 {
		return Task{}, fmt.Errorf("task %s: package, target and hidden are required", t.ID)
	}
	if t.Behaviour == "" || t.Signatures == "" || t.Code == "" {
		return Task{}, fmt.Errorf("task %s: behaviour, signatures and code intents are all required", t.ID)
	}
	for _, f := range t.Hidden {
		if !strings.HasSuffix(f, "_test.go") {
			return Task{}, fmt.Errorf("task %s: hidden file %s must end in _test.go", t.ID, f)
		}
		if _, err := os.Stat(filepath.Join(dir, "hidden", f)); err != nil {
			return Task{}, fmt.Errorf("task %s: %w", t.ID, err)
		}
	}
	for _, f := range t.Context {
		if _, err := os.Stat(filepath.Join(dir, "context", f)); err != nil {
			return Task{}, fmt.Errorf("task %s: %w", t.ID, err)
		}
	}
	for _, f := range t.Target {
		for _, other := range append(append([]string{}, t.Context...), t.Hidden...) {
			if f == other {
				return Task{}, fmt.Errorf("task %s: target %s is also a context or hidden file", t.ID, f)
			}
		}
	}
	return t, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// LoadTasks loads every subdirectory of root that holds a task.yaml, sorted by\n// task id.\nfunc LoadTasks(root string) ([]Task, error) {\n\tentries, err := os.ReadDir(root)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"read tasks: %w\", err)\n\t}\n\tvar tasks []Task\n\tfor _, e := range entries {\n\t\tif !e.IsDir() {\n\t\t\tcontinue\n\t\t}\n\t\tif _, err := os.Stat(filepath.Join(root, e.Name(), \"task.yaml\")); err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tt, err := LoadTask(filepath.Join(root, e.Name()))\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\ttasks = append(tasks, t)\n\t}\n\tsort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })\n\treturn tasks, nil\n}","file":"evals/design/task.go","name":"LoadTasks"}}
// LoadTasks loads every subdirectory of root that holds a task.yaml, sorted by
// task id.
func LoadTasks(root string) ([]Task, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	var tasks []Task
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, e.Name(), "task.yaml")); err != nil {
			continue
		}
		t, err := LoadTask(filepath.Join(root, e.Name()))
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks, nil
}

// /sedum:setFunction
