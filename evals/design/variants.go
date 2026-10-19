package design

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"evals/design/variants.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Variant is how much of the answer an intent gives away.\ntype Variant string","file":"evals/design/variants.go","name":"Variant"}}
// Variant is how much of the answer an intent gives away.
type Variant string

// /sedum:setType
// sedum:setConst {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// The three intent variants, from least to most given.\nconst (\n\tVariantBehaviour  Variant = \"behaviour\"\n\tVariantSignatures Variant = \"signatures\"\n\tVariantCode       Variant = \"code\"\n)","file":"evals/design/variants.go","name":"VariantBehaviour"}}
// The three intent variants, from least to most given.
const (
	VariantBehaviour  Variant = "behaviour"
	VariantSignatures Variant = "signatures"
	VariantCode       Variant = "code"
)

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Variants returns every variant in order.\nfunc Variants() []Variant {\n\treturn []Variant{VariantBehaviour, VariantSignatures, VariantCode}\n}","file":"evals/design/variants.go","name":"Variants"}}
// Variants returns every variant in order.
func Variants() []Variant {
	return []Variant{VariantBehaviour, VariantSignatures, VariantCode}
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Intent returns the task's intent text for v.\nfunc (t Task) Intent(v Variant) (string, error) {\n\tswitch v {\n\tcase VariantBehaviour:\n\t\treturn t.Behaviour, nil\n\tcase VariantSignatures:\n\t\treturn t.Signatures, nil\n\tcase VariantCode:\n\t\treturn t.Code, nil\n\t}\n\treturn \"\", fmt.Errorf(\"unknown variant %q\", v)\n}","file":"evals/design/variants.go","name":"Intent","receiver":"Task"}}
// Intent returns the task's intent text for v.
func (t Task) Intent(v Variant) (string, error) {
	switch v {
	case VariantBehaviour:
		return t.Behaviour, nil
	case VariantSignatures:
		return t.Signatures, nil
	case VariantCode:
		return t.Code, nil
	}
	return "", fmt.Errorf("unknown variant %q", v)
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Isolated returns an error when hidden test text could reach the model: a\n// hidden file's name, or any of its \"func Test\" lines, appearing in any of the\n// task's intents.\nfunc (t Task) Isolated() error {\n\tfor _, v := range Variants() {\n\t\tintent, err := t.Intent(v)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfor _, f := range t.Hidden {\n\t\t\tif strings.Contains(intent, f) {\n\t\t\t\treturn fmt.Errorf(\"task %s: the %s intent names hidden file %s\", t.ID, v, f)\n\t\t\t}\n\t\t\tdata, err := os.ReadFile(filepath.Join(t.Dir, \"hidden\", f))\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\tfor _, line := range strings.Split(string(data), \"\\n\") {\n\t\t\t\tif strings.HasPrefix(line, \"func Test\") && strings.Contains(intent, line) {\n\t\t\t\t\treturn fmt.Errorf(\"task %s: the %s intent contains the hidden test line %q\", t.ID, v, line)\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n\treturn nil\n}","file":"evals/design/variants.go","name":"Isolated","receiver":"Task"}}
// Isolated returns an error when hidden test text could reach the model: a
// hidden file's name, or any of its "func Test" lines, appearing in any of the
// task's intents.
func (t Task) Isolated() error {
	for _, v := range Variants() {
		intent, err := t.Intent(v)
		if err != nil {
			return err
		}
		for _, f := range t.Hidden {
			if strings.Contains(intent, f) {
				return fmt.Errorf("task %s: the %s intent names hidden file %s", t.ID, v, f)
			}
			data, err := os.ReadFile(filepath.Join(t.Dir, "hidden", f))
			if err != nil {
				return err
			}
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "func Test") && strings.Contains(intent, line) {
					return fmt.Errorf("task %s: the %s intent contains the hidden test line %q", t.ID, v, line)
				}
			}
		}
	}
	return nil
}

// /sedum:setMethod
