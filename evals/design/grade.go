package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/grade.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/grade.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/grade.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// tail keeps the last n bytes of s.\nfunc tail(s string, n int) string {\n\tif len(s) <= n {\n\t\treturn s\n\t}\n\treturn s[len(s)-n:]\n}","file":"evals/design/grade.go","name":"tail"}}
// tail keeps the last n bytes of s.
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Grade checks what a grow left in dir. It runs go build, which is\n// \"compiled\"; then copies the task's hidden tests in and runs go test, which\n// is \"passed\". The returned output is the tail of whatever failed.\nfunc Grade(t Task, dir string) (compiled, passed bool, output string, err error) {\n\trun := func(args ...string) (string, error) {\n\t\tcmd := exec.Command(\"go\", args...)\n\t\tcmd.Dir = dir\n\t\tout, err := cmd.CombinedOutput()\n\t\treturn string(out), err\n\t}\n\tif out, err := run(\"build\", \"./...\"); err != nil {\n\t\treturn false, false, tail(out, 2000), nil\n\t}\n\tcompiled = true\n\tfor _, f := range t.Hidden {\n\t\tdata, err := os.ReadFile(filepath.Join(t.Dir, \"hidden\", f))\n\t\tif err != nil {\n\t\t\treturn compiled, false, \"\", err\n\t\t}\n\t\tif err := os.WriteFile(filepath.Join(dir, t.Package, f), data, 0o644); err != nil {\n\t\t\treturn compiled, false, \"\", err\n\t\t}\n\t}\n\tout, testErr := run(\"test\", \"-count=1\", \"./...\")\n\tif testErr != nil {\n\t\treturn compiled, false, tail(out, 2000), nil\n\t}\n\treturn compiled, true, \"\", nil\n}","file":"evals/design/grade.go","name":"Grade"}}
// Grade checks what a grow left in dir. It runs go build, which is
// "compiled"; then copies the task's hidden tests in and runs go test, which
// is "passed". The returned output is the tail of whatever failed.
func Grade(t Task, dir string) (compiled, passed bool, output string, err error) {
	run := func(args ...string) (string, error) {
		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	if out, err := run("build", "./..."); err != nil {
		return false, false, tail(out, 2000), nil
	}
	compiled = true
	for _, f := range t.Hidden {
		data, err := os.ReadFile(filepath.Join(t.Dir, "hidden", f))
		if err != nil {
			return compiled, false, "", err
		}
		if err := os.WriteFile(filepath.Join(dir, t.Package, f), data, 0o644); err != nil {
			return compiled, false, "", err
		}
	}
	out, testErr := run("test", "-count=1", "./...")
	if testErr != nil {
		return compiled, false, tail(out, 2000), nil
	}
	return compiled, true, "", nil
}

// /sedum:setFunction
