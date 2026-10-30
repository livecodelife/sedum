package design

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants.go","path":"fmt"},"kwargs":"sha256:c8a910f1ce483e70c4036fadadf2e599736e290ce13f78b23b438978233eddbf","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants.go","path":"os"},"kwargs":"sha256:77636674272fca93e9d51d7df2217e8c190cf1534946212f44abd2bdeb59846c","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants.go","path":"path/filepath"},"kwargs":"sha256:3459b849e9b9d1ec53d994175c4a2f3eafd1cd0616fbbad1ff5273f50c6a5e39","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"evals/design/variants.go","path":"strings"},"kwargs":"sha256:55894839587b51bf63585ba98730eb4396d22584d0816d9c434daf27ded1f1e9","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"evals/design/variants.go","name":"Variant"},"kwargs":"sha256:5f14a3c4c2067eee8228f2692e362bdb37f545f749af53cee6b35a69577d0142","body":"sha256:5d4150fa997e0186305b768f67e942dce8cacba55c19188294e462bc32554d6d"}
// Variant is how much of the answer an intent gives away.
type Variant string

// /sedum:setType
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"evals/design/variants.go","name":"VariantBehaviour"},"kwargs":"sha256:fbcdc7b2cfc49660dcbc9c77756e82f2fcdeb7641dd80c9fcf3f9aa296bc8aac","body":"sha256:19aa0fc0a30cfce11b6067fc226b55f7380b120d4e02e28acbcae6fad928ef93"}
// The three intent variants, from least to most given.
const (
	VariantBehaviour  Variant = "behaviour"
	VariantSignatures Variant = "signatures"
	VariantCode       Variant = "code"
)

// /sedum:setConst
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"evals/design/variants.go","name":"Variants"},"kwargs":"sha256:b622486e61519b46211451b295ba03237f7908f4ae594c27486053ef828b2952","body":"sha256:398f7aa41dfedac514655ef15fb81bd4c19b4d6b370031c9cad20e0091c73043"}
// Variants returns every variant in order.
func Variants() []Variant {
	return []Variant{VariantBehaviour, VariantSignatures, VariantCode}
}

// /sedum:setFunction
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"evals/design/variants.go","name":"Intent","receiver":"Task"},"kwargs":"sha256:575453bda6220de928f854b02ac703f8929c2aca0cddc6e8a9e813507175954c","body":"sha256:c52ad94bc5e11c72babe79ca98e273e2a87b82ea3e125e8d77c4f504efc29ee8"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-84c8ca99","gen":"sha256:bdcd1b73d07f62bb1943c9864fd4f0e31a5dcc1bf700073bf0c10cdb9cac0e01","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"evals/design/variants.go","name":"Isolated","receiver":"Task"},"kwargs":"sha256:8f5d5778e09f121330da39c97dcbe52812446a122abd4427e65063c20fe1593c","body":"sha256:3176fd32823243cd30350f29833947f924017522a16056ad51f82a5f342ab9d3"}
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
