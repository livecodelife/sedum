package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/grade.go","path":"os"},"kwargs":"sha256:7ce9ec40e1ac8d635938e0a2533764cb780a552c2ac8f9fcc2ac1e52c3238484","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/grade.go","path":"os/exec"},"kwargs":"sha256:91b8bc15e443f85cf4061fba0448d420237e57b63a0e69427064acc113808262","body":"sha256:5970e9ca793f55e2beb08de623e98bdbc0ced178860fedd953eef7e15026ee19"}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/grade.go","path":"path/filepath"},"kwargs":"sha256:cee206797df7e805fa5227e7f6a6e96db74c538e519ea26fcca8928f65cb1cba","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/grade.go","name":"tail"},"kwargs":"sha256:65fb428961103276946b1f0f52195bdb41511cde765984162167938469c11644","body":"sha256:01f41d503ed921b1e1f710a1b92f1329e86b0f81da3a420cc0bfe19e90ec660d"}
// tail keeps the last n bytes of s.
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3576b512","gen":"sha256:8fdf237b1cfc5b39c2b61b113e19348370b6fe3df88a1d96d16e3acf278f4cef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/grade.go","name":"Grade"},"kwargs":"sha256:3cdcd8ae7812a57adfa6dec364d823d5058de6a9564e91c617b08cc15466a73c","body":"sha256:6ed77ae10f2bd8b6faf0a1c59eeae8f8cc6834db68de74e311bf1eaab3232040"}
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
