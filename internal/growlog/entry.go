package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/entry.go","path":"unicode/utf8"}}
import "unicode/utf8"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// Version is the schema version written in every line of the grow log.\nconst Version = 1","file":"internal/growlog/entry.go","name":"Version"}}
// Version is the schema version written in every line of the grow log.
const Version = 1

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// maxError bounds the error text kept in one line.\nconst maxError = 300","file":"internal/growlog/entry.go","name":"maxError"}}
// maxError bounds the error text kept in one line.
const maxError = 300

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Regions counts what one grow did to a record's regions.\ntype Regions struct {\n\tInjected  int `json:\"injected\"`\n\tReplaced  int `json:\"replaced\"`\n\tPruned    int `json:\"pruned\"`\n\tUnchanged int `json:\"unchanged\"`\n}","file":"internal/growlog/entry.go","name":"Regions"}}
// Regions counts what one grow did to a record's regions.
type Regions struct {
	Injected  int `json:"injected"`
	Replaced  int `json:"replaced"`
	Pruned    int `json:"pruned"`
	Unchanged int `json:"unchanged"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Entry is one line of the grow log: what one grow did for one record. It\n// holds counts and timings only, never prompt or response text, kwargs, file\n// contents or keys.\ntype Entry struct {\n\tV                int     `json:\"v\"`\n\tTS               string  `json:\"ts\"`\n\tSedum            string  `json:\"sedum\"`\n\tRecord           string  `json:\"record\"`\n\tMode             string  `json:\"mode\"`\n\tModel            string  `json:\"model\"`\n\tCalls            int     `json:\"calls\"`\n\tRejected         int     `json:\"rejected\"`\n\tCompleteness     int     `json:\"completeness\"`\n\tRepairs          int     `json:\"repairs\"`\n\tPromptTokens     int     `json:\"prompt_tokens\"`\n\tCompletionTokens int     `json:\"completion_tokens\"`\n\tInvocations      int     `json:\"invocations\"`\n\tLiteral          int     `json:\"literal\"`\n\tDeclarative      int     `json:\"declarative\"`\n\tRegions          Regions `json:\"regions\"`\n\tFiles            int     `json:\"files\"`\n\tMS               int64   `json:\"ms\"`\n\tOutcome          string  `json:\"outcome\"`\n\tError            string  `json:\"error,omitempty\"`\n}","file":"internal/growlog/entry.go","name":"Entry"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// clipError keeps the first line of s, trimmed, at most maxError bytes and\n// never cutting a UTF-8 sequence.\nfunc clipError(s string) string {\n\ts, _, _ = strings.Cut(s, \"\\n\")\n\ts = strings.TrimSpace(s)\n\tif len(s) <= maxError {\n\t\treturn s\n\t}\n\ts = s[:maxError]\n\tfor !utf8.ValidString(s) {\n\t\ts = s[:len(s)-1]\n\t}\n\treturn s\n}","file":"internal/growlog/entry.go","name":"clipError"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Append adds e as one line to the log at path, creating the file and its\n// directory. V is set to Version when zero and Error is clipped. The line is\n// written with a single write to an O_APPEND file, so concurrent writers never\n// interleave their lines.\nfunc Append(path string, e Entry) error {\n\tif e.V == 0 {\n\t\te.V = Version\n\t}\n\te.Error = clipError(e.Error)\n\tline, err := json.Marshal(e)\n\tif err != nil {\n\t\treturn fmt.Errorf(\"encode grow log entry: %w\", err)\n\t}\n\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\n\t\treturn fmt.Errorf(\"create grow log directory for %s: %w\", path, err)\n\t}\n\tf, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)\n\tif err != nil {\n\t\treturn fmt.Errorf(\"open grow log %s: %w\", path, err)\n\t}\n\tif _, err := f.Write(append(line, '\\n')); err != nil {\n\t\tf.Close()\n\t\treturn fmt.Errorf(\"write grow log %s: %w\", path, err)\n\t}\n\treturn f.Close()\n}","file":"internal/growlog/entry.go","name":"Append"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-21aa9103","gen":"sha256:9c70e093911787cc75b2a943211aa011b72fe7c232f5b504f388170211fee6ef","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Read returns the entries in the log at path, in file order, and the number\n// of lines that could not be decoded. A missing file yields no entries.\nfunc Read(path string) ([]Entry, int, error) {\n\tdata, err := os.ReadFile(path)\n\tif os.IsNotExist(err) {\n\t\treturn nil, 0, nil\n\t}\n\tif err != nil {\n\t\treturn nil, 0, fmt.Errorf(\"read grow log %s: %w\", path, err)\n\t}\n\tvar entries []Entry\n\tbad := 0\n\tfor _, line := range strings.Split(string(data), \"\\n\") {\n\t\tif line == \"\" {\n\t\t\tcontinue\n\t\t}\n\t\tvar e Entry\n\t\tif err := json.Unmarshal([]byte(line), &e); err != nil {\n\t\t\tbad++\n\t\t\tcontinue\n\t\t}\n\t\tentries = append(entries, e)\n\t}\n\treturn entries, bad, nil\n}","file":"internal/growlog/entry.go","name":"Read"}}
// Read returns the entries in the log at path, in file order, and the number
// of lines that could not be decoded. A missing file yields no entries.
func Read(path string) ([]Entry, int, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("read grow log %s: %w", path, err)
	}
	var entries []Entry
	bad := 0
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			bad++
			continue
		}
		entries = append(entries, e)
	}
	return entries, bad, nil
}

// /sedum:setFunction
