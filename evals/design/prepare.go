package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"fmt"},"kwargs":"sha256:b257cdea935ce463d86174fed588207491e65175d7437315463528e184f708cf","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"io/fs"},"kwargs":"sha256:452b93f71adf44ca587c740a7929bd6bcc19a1aaed06bf6864eb13441c13682f","body":"sha256:80f0597eb523a3ace402ac68fefdf64bf00074dda5ec249d04e8f1e525952547"}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"os"},"kwargs":"sha256:c942067d388ca784cc2b94dd2ce9cd2083692fd4382f1779ee71773aee5c8225","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"path/filepath"},"kwargs":"sha256:4067d24d84bf68e5cd4a466a02630d07f0ef9401ca5384f266c0c3c44b37b7bc","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"gopkg.in/yaml.v3"},"kwargs":"sha256:6cf5713f052e29e684d16f98fe15d8e565ea511f864ccb75f21156ee8ee6a6f8","body":"sha256:a4214cbff710690b396b73d2b1ccc9aa8b55a807d07f73aa80f10af0c41e13b9"}
import "gopkg.in/yaml.v3"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/prepare.go","path":"strings"},"kwargs":"sha256:b0f4c60907f90278431626ad4a631b562616127828d35ca83b3e9b8bcdadaf95","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"evals/design/prepare.go","name":"growPrompt"},"kwargs":"sha256:af52f2aca0b126c8e71eccf1d732a013e2559e0e169c0e2e329a7f204ad25e02","body":"sha256:9974897893e44c24cb0082d91b4a5b30bdd4f72c5abe59543bd95d61c20f35bf"}
// growPrompt is the same advisory text for every variant, so the variants
// differ only in the intent.
const growPrompt = "# Context\n\nA small Go package. Standard library only (plus what the context files import). Write each target file in full, gofmt-clean, with a doc comment on every exported identifier.\n"

// /sedum:setConst
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"evals/design/prepare.go","name":"record"},"kwargs":"sha256:602575ad75139ae4dce1c6f07ec2a58ce3f8ee38da76ca179d2f71612aca32f8","body":"sha256:29e87879c2e561513c855557e25dcdaab463815fc53dd5014239647c18a6c9dc"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/prepare.go","name":"copyTree"},"kwargs":"sha256:a9b12adcad4a4a0417d0d3e1e171a6974ee82ee95f10e0029217416265d81354","body":"sha256:1f0dbb2989dffcea94a72f28457abd49d5620c376de83e38962495e5c2c58527"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/prepare.go","name":"Prepare"},"kwargs":"sha256:4098885c1a69daa22b70d5e33c58c30062f5ad45540cab8a4e86388c5b016a03","body":"sha256:d180f42defb84addf61b1eb491b052c14fbb9da63cd036b2547cfc677c590757"}
// Prepare writes the project for one sample into dir: a go.mod, the task's
// context files (in the directory named for the task's package), a copy of the generator packages, the grow prompt, and one
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
	if err := os.MkdirAll(filepath.Join(dir, t.Package), 0o755); err != nil {
		return err
	}
	for _, f := range t.Context {
		data, err := os.ReadFile(filepath.Join(t.Dir, "context", f))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, t.Package, f), data, 0o644); err != nil {
			return err
		}
	}
	if err := copyTree(generators, filepath.Join(dir, "generators")); err != nil {
		return fmt.Errorf("copy generators: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "PROMPT.md"), []byte(growPrompt+"\nThe target files, as project paths:\n"+targetPaths(t)), 0o644); err != nil {
		return err
	}
	scope := make([]string, 0, len(t.Target))
	for _, f := range t.Target {
		scope = append(scope, t.Package+"/"+f)
	}
	rec := record{
		ID:            fmt.Sprintf("design-%s-%s", t.ID, v),
		Title:         fmt.Sprintf("design eval: %s (%s)", t.ID, v),
		Status:        "open",
		CreatedAt:     "2026-10-16",
		Author:        "eval@example.com",
		Intent:        intent,
		Constraints:   []string{},
		AffectedScope: scope,
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/prepare.go","name":"targetPaths"},"kwargs":"sha256:660874c58a5e361a077d06a8b3860845c1699457a399b1b320690734ee8aa819","body":"sha256:74a2c5627aaabffbe579bc35e0fd9d98048e03d384e836896b780fea9a252f8d"}
// targetPaths lists the task's target files as project paths, one per line.
func targetPaths(t Task) string {
	var b strings.Builder
	for _, f := range t.Target {
		b.WriteString(t.Package + "/" + f + "\n")
	}
	return b.String()
}

// /sedum:setFunction
