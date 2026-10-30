package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"encoding/json"},"kwargs":"sha256:573df20ef159ede2daa4c6e7f4a894390a4765300448a5fa0d51a72dd4ea55a3","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"fmt"},"kwargs":"sha256:7eb2177804abf2ae69e72d07544e9743cd45acffb77794299d282bbe345a0725","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"os"},"kwargs":"sha256:c36c4667849b84c4b25b6a93635f90e916f7db716af976a414fcec6a157005fa","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"path/filepath"},"kwargs":"sha256:0de5eae0597aa5ffd77601164e7611d05b86287dbdfb383e52a662c49ec1cac1","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"strings"},"kwargs":"sha256:c81c211d018b00165ec9951a6c51b6addaa61e1f0096a144be4ca20044e524c4","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/entry.go","path":"unicode/utf8"},"kwargs":"sha256:44bd8c0d95118ff0ecb585fc9081cb001db032b1f7d68e4310142c07fc79d846","body":"sha256:329f7334124c7cfe72d8208daedf694eae268e0a4bd205bdde0ead7fbb2a0037"}
import "unicode/utf8"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/growlog/entry.go","name":"Version"},"kwargs":"sha256:afeaeb30ef76c7c8cd2e9a4f357e2d11dfa4dfae4f13a5bb84bbb0a9cca3d896","body":"sha256:1e0cd46d476112c3f7dfec4835756e66a1224ec8097edd0c0092a3e90ffabb28"}
// Version is the schema version written in every line of the grow log.
const Version = 1

// /sedum:setConst
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/growlog/entry.go","name":"maxError"},"kwargs":"sha256:3dad56583e9c6bef9ae2c4e335efb95e731a38e82c6e24b81584ad01d08220d6","body":"sha256:40993f3d7b7157bbc015f1eb4c649f17bc17d3259b31ad8d3021c8555b13eaac"}
// maxError bounds the error text kept in one line.
const maxError = 300

// /sedum:setConst
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/growlog/entry.go","name":"Regions"},"kwargs":"sha256:e5b71dd2f20bf02a5d54169f1d6627feca6adb94e779e844ec08fb1ca0a8f8d0","body":"sha256:700f4f94d2cede9056fdb2d10d2f2e00a37fc03ccf9966d3de56ca224c2cd4e6"}
// Regions counts what one grow did to a record's regions.
type Regions struct {
	Injected  int `json:"injected"`
	Replaced  int `json:"replaced"`
	Pruned    int `json:"pruned"`
	Unchanged int `json:"unchanged"`
}

// /sedum:setType
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/growlog/entry.go","name":"Entry"},"kwargs":"sha256:67c8d15fe1474e3c19457e6e2945550a1afbc90a894384563a83145bcc1166d1","body":"sha256:cfcefc44f75a09461f24b8b201316225a94bf04004fc8497046de121a3e90532"}
// Entry is one line of the grow log: what one grow did for one record. It
// holds counts and timings only, never prompt or response text, kwargs, file
// contents or keys.
type Entry struct {
	V                int     `json:"v"`
	TS               string  `json:"ts"`
	Sedum            string  `json:"sedum"`
	Record           string  `json:"record"`
	Mode             string  `json:"mode"`
	Model            string  `json:"model"`
	Calls            int     `json:"calls"`
	Rejected         int     `json:"rejected"`
	Completeness     int     `json:"completeness"`
	Repairs          int     `json:"repairs"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	Invocations      int     `json:"invocations"`
	Literal          int     `json:"literal"`
	Declarative      int     `json:"declarative"`
	Regions          Regions `json:"regions"`
	Files            int     `json:"files"`
	MS               int64   `json:"ms"`
	Outcome          string  `json:"outcome"`
	Error            string  `json:"error,omitempty"`
}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/entry.go","name":"clipError"},"kwargs":"sha256:20f30e073b2c919356aced1d2765fa658f94cbc47c267936124b403113ea4b4c","body":"sha256:cb16d0c53d0784ff0cc4956f39cc09ddabf865045615d336e01bc6b1d67b4d82"}
// clipError keeps the first line of s, trimmed, at most maxError bytes and
// never cutting a UTF-8 sequence.
func clipError(s string) string {
	s, _, _ = strings.Cut(s, "\n")
	s = strings.TrimSpace(s)
	if len(s) <= maxError {
		return s
	}
	s = s[:maxError]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:886221717c81720fd4120a4228dec3e8266e2e7711207ba2d373f07c72104c66","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/entry.go","name":"Append"},"kwargs":"sha256:3754fc86c86eaa9d6a5c84661e658de356625c290882f5f69c2c8efcdeaa4429","body":"sha256:38d3cc292e21e87da32465743377d724287388bbe80aca71028b295ce4d22b47"}
// Append adds e as one line to the log at path, creating the file and its
// directory. V is set to Version when zero and Error is clipped. The line is
// written with a single write to an O_APPEND file, so concurrent writers never
// interleave their lines.
func Append(path string, e Entry) error {
	if e.V == 0 {
		e.V = Version
	}
	e.Error = clipError(e.Error)
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("encode grow log entry: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create grow log directory for %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open grow log %s: %w", path, err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		f.Close()
		return fmt.Errorf("write grow log %s: %w", path, err)
	}
	return f.Close()
}

// /sedum:setFunction
