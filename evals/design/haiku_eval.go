package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/haiku_eval.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// HaikuGrow returns a Grow that runs sedum through the herd's sedum-haiku\n// wrapper, which pins claude-haiku-5-5. The harness passes no model flag, and\n// the wrapper refuses one: no other model is allowed.\nfunc HaikuGrow(wrapper, sedumBin string) func(ctx context.Context, dir string) error {\n\treturn func(ctx context.Context, dir string) error {\n\t\tcmd := exec.CommandContext(ctx, wrapper, \"grow\",\n\t\t\t\"--generators\", \"./generators\", \"--records\", \"./provenance\",\n\t\t\t\"--prompt-file\", \"PROMPT.md\", \"--verify\", \"go vet ./...\")\n\t\tcmd.Dir = dir\n\t\tcmd.Env = append(os.Environ(), \"SEDUM_BIN=\"+sedumBin)\n\t\tout, err := cmd.CombinedOutput()\n\t\tif err != nil {\n\t\t\tlines := strings.Split(strings.TrimSpace(string(out)), \"\\n\")\n\t\t\treturn fmt.Errorf(\"%s\", lines[len(lines)-1])\n\t\t}\n\t\treturn nil\n\t}\n}","file":"evals/design/haiku_eval.go","name":"HaikuGrow"}}
// HaikuGrow returns a Grow that runs sedum through the herd's sedum-haiku
// wrapper, which pins claude-haiku-5-5. The harness passes no model flag, and
// the wrapper refuses one: no other model is allowed.
func HaikuGrow(wrapper, sedumBin string) func(ctx context.Context, dir string) error {
	return func(ctx context.Context, dir string) error {
		cmd := exec.CommandContext(ctx, wrapper, "grow",
			"--generators", "./generators", "--records", "./provenance",
			"--prompt-file", "PROMPT.md", "--verify", "go vet ./...")
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "SEDUM_BIN="+sedumBin)
		out, err := cmd.CombinedOutput()
		if err != nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			return fmt.Errorf("%s", lines[len(lines)-1])
		}
		return nil
	}
}

// /sedum:setFunction
