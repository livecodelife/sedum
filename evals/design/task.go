package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"fmt"},"kwargs":"sha256:c896c3ee415ab56358394b5cd3c761169d96612f9c129735d1d9ab0fbe75364a","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"os"},"kwargs":"sha256:9050607a26a9ed3be9bafe99efb46ab02e9d5465839ce1224287132dc81277ad","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"path/filepath"},"kwargs":"sha256:37b0ce75cd2ae5a2e543f2fd8cac9bdf7bf1709ddf57df9afdaee7f4b3fa7863","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"sort"},"kwargs":"sha256:71db451fdb9465e5bb2ff9000f1f9209b30afb34c4ab4e34c9a33728abaa1b8e","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"strings"},"kwargs":"sha256:9ef5a33ff8ec6d249f0f3967dd152a1cc0021de2a0aa5b6290eab28cb024a15a","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/task.go","path":"gopkg.in/yaml.v3"},"kwargs":"sha256:12123e00d06ef5a12103ef1b4749927e608379f4729499c5bf78cf71cd8abaec","body":"sha256:a4214cbff710690b396b73d2b1ccc9aa8b55a807d07f73aa80f10af0c41e13b9"}
import "gopkg.in/yaml.v3"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"evals/design/task.go","name":"Task"},"kwargs":"sha256:1b7471347233968e8ac7a58df9695e5bc93eba04f5b38ba4268fa604956c92f4","body":"sha256:f4d3f39d2cafacc47335e918128e179c42a8c26ca2ce672bb39e30ec22b7171d"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task.go","name":"LoadTask"},"kwargs":"sha256:c167f9c281105d0813c5d767fdacbdf8a0d0a06a0d74cc931471e547514585f7","body":"sha256:b0ac7504c9aa74345021590bcf6d478f7a0493b1d9a53560449a3627d27f83e0"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/task.go","name":"LoadTasks"},"kwargs":"sha256:4f259b2e97ea693e6c1ef608dbfb1b774d832e48268e12ac5484df523172c02f","body":"sha256:c9c717b32ecc35528fae5172eff3f89986a38810cfc1726b8a111c052a145539"}
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
