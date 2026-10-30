package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval.go","path":"context"},"kwargs":"sha256:3087f9a3fca111c468734df66b2083814f026207fef8c73287418c785fab2237","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval.go","path":"fmt"},"kwargs":"sha256:4abea3ec02b724080d86ed346d7c519fc513b9dc18eedd3f5d70041c9b5a7d53","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval.go","path":"os"},"kwargs":"sha256:668c7c28e0e257a2ed69a4b2f08faed9c38998936d2f5840b73a9c49aa988c25","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval.go","path":"os/exec"},"kwargs":"sha256:14f49df914d3670fe3491153ed3871413f7873e3ae54ecca9e86cc25196e2af9","body":"sha256:5970e9ca793f55e2beb08de623e98bdbc0ced178860fedd953eef7e15026ee19"}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/haiku_eval.go","path":"strings"},"kwargs":"sha256:d31231cfba6c38c2c1c8a27b1a7d643ef09e1ae6cad8b1634f2f8b9663e29b92","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-860c8854","gen":"sha256:f20fe66ab2f4aae8f070d727f9f8a3fdf95b3e3e50ff7503ddf688fd4e4d5f4e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/haiku_eval.go","name":"HaikuGrow"},"kwargs":"sha256:c57e36615d320b85e79a32ada175a9a8bd3b545888af75ec42919a211edcbd83","body":"sha256:ed22e8f93c20b0af82366191f5d792887c4266b7574426adfb4c8efafb2a6a66"}
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
