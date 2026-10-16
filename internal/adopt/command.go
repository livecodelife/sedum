package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/command.go","path":"time"}}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// commandTimeout bounds one run of the adopter command.\nconst commandTimeout = 60 * time.Second","file":"internal/adopt/command.go","name":"commandTimeout"}}
// commandTimeout bounds one run of the adopter command.
const commandTimeout = 60 * time.Second

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// maxCommandOutput bounds the stdout of one run of the adopter command.\nconst maxCommandOutput = 64 << 20","file":"internal/adopt/command.go","name":"maxCommandOutput"}}
// maxCommandOutput bounds the stdout of one run of the adopter command.
const maxCommandOutput = 64 << 20

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// CommandAdopter splits files by running an external command the user allowed.\ntype CommandAdopter struct {\n\tArgv []string\n}","file":"internal/adopt/command.go","name":"CommandAdopter"}}
// CommandAdopter splits files by running an external command the user allowed.
type CommandAdopter struct {
	Argv []string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// ParseSpans reads one JSON document of the form {\"spans\":[...]} and returns its spans.\n// Unknown fields, trailing data, and an empty or missing spans list are errors.\n// The spans are not validated here; CheckSpans does that.\nfunc ParseSpans(r io.Reader) ([]Span, error) {\n\tdec := json.NewDecoder(r)\n\tdec.DisallowUnknownFields()\n\tvar doc struct {\n\t\tSpans []Span `json:\"spans\"`\n\t}\n\tif err := dec.Decode(&doc); err != nil {\n\t\treturn nil, fmt.Errorf(\"decoding spans: %w\", err)\n\t}\n\tif err := dec.Decode(&struct{}{}); err != io.EOF {\n\t\treturn nil, errors.New(\"trailing data after spans document\")\n\t}\n\tif len(doc.Spans) == 0 {\n\t\treturn nil, errors.New(\"adopter returned no spans\")\n\t}\n\treturn doc.Spans, nil\n}","file":"internal/adopt/command.go","name":"ParseSpans"}}
// ParseSpans reads one JSON document of the form {"spans":[...]} and returns its spans.
// Unknown fields, trailing data, and an empty or missing spans list are errors.
// The spans are not validated here; CheckSpans does that.
func ParseSpans(r io.Reader) ([]Span, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var doc struct {
		Spans []Span `json:"spans"`
	}
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("decoding spans: %w", err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, errors.New("trailing data after spans document")
	}
	if len(doc.Spans) == 0 {
		return nil, errors.New("adopter returned no spans")
	}
	return doc.Spans, nil
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Split runs the command with path appended to its arguments and parses the\n// spans it prints on stdout. The command reads the file itself, so src is not\n// passed to it.\nfunc (c CommandAdopter) Split(path string, src []byte) ([]Span, error) {\n\tif len(c.Argv) == 0 {\n\t\treturn nil, errors.New(\"command adopter: empty command line\")\n\t}\n\tctx, cancel := context.WithTimeout(context.Background(), commandTimeout)\n\tdefer cancel()\n\targs := append(append([]string{}, c.Argv[1:]...), path)\n\tcmd := exec.CommandContext(ctx, c.Argv[0], args...)\n\tvar stdout, stderr bytes.Buffer\n\tcmd.Stdout = &stdout\n\tcmd.Stderr = &stderr\n\tif err := cmd.Run(); err != nil {\n\t\tvar exitErr *exec.ExitError\n\t\tif errors.As(err, &exitErr) {\n\t\t\tmsg := strings.SplitN(stderr.String(), \"\\n\", 2)[0]\n\t\t\treturn nil, fmt.Errorf(\"command adopter %s: exit status %d: %s\", c.Argv[0], exitErr.ExitCode(), msg)\n\t\t}\n\t\treturn nil, fmt.Errorf(\"command adopter %s: %w\", c.Argv[0], err)\n\t}\n\tif stdout.Len() > maxCommandOutput {\n\t\treturn nil, fmt.Errorf(\"command adopter %s: stdout exceeds %d bytes\", c.Argv[0], maxCommandOutput)\n\t}\n\tspans, err := ParseSpans(&stdout)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"command adopter %s: %w\", c.Argv[0], err)\n\t}\n\treturn spans, nil\n}","file":"internal/adopt/command.go","name":"Split","receiver":"CommandAdopter"}}
// Split runs the command with path appended to its arguments and parses the
// spans it prints on stdout. The command reads the file itself, so src is not
// passed to it.
func (c CommandAdopter) Split(path string, src []byte) ([]Span, error) {
	if len(c.Argv) == 0 {
		return nil, errors.New("command adopter: empty command line")
	}
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	args := append(append([]string{}, c.Argv[1:]...), path)
	cmd := exec.CommandContext(ctx, c.Argv[0], args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			msg := strings.SplitN(stderr.String(), "\n", 2)[0]
			return nil, fmt.Errorf("command adopter %s: exit status %d: %s", c.Argv[0], exitErr.ExitCode(), msg)
		}
		return nil, fmt.Errorf("command adopter %s: %w", c.Argv[0], err)
	}
	if stdout.Len() > maxCommandOutput {
		return nil, fmt.Errorf("command adopter %s: stdout exceeds %d bytes", c.Argv[0], maxCommandOutput)
	}
	spans, err := ParseSpans(&stdout)
	if err != nil {
		return nil, fmt.Errorf("command adopter %s: %w", c.Argv[0], err)
	}
	return spans, nil
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-1d15d075","gen":"sha256:2049c79293531ee9f05c025fc651f408d4f08c41b3c7c9847076b368bc2289ee","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// String returns the command line, Argv joined by single spaces.\nfunc (c CommandAdopter) String() string {\n\treturn strings.Join(c.Argv, \" \")\n}","file":"internal/adopt/command.go","name":"String","receiver":"CommandAdopter"}}
// String returns the command line, Argv joined by single spaces.
func (c CommandAdopter) String() string {
	return strings.Join(c.Argv, " ")
}

// /sedum:setMethod
