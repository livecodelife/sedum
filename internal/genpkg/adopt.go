package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt.go","path":"regexp"}}
import "regexp"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Adopt is the adopt: block of sedum.yaml. It brings existing files under Sedum.\ntype Adopt struct {\n\tCommand     []string          `yaml:\"command\"`\n\tUnitStart   string            `yaml:\"unit_start\"`\n\tUnitEnd     string            `yaml:\"unit_end\"`\n\tHeader      string            `yaml:\"header\"`\n\tAttachAbove []string          `yaml:\"attach_above\"`\n\tActions     map[string]string `yaml:\"actions\"`\n}","file":"internal/genpkg/adopt.go","name":"Adopt"}}
// Adopt is the adopt: block of sedum.yaml. It brings existing files under Sedum.
type Adopt struct {
	Command     []string          `yaml:"command"`
	UnitStart   string            `yaml:"unit_start"`
	UnitEnd     string            `yaml:"unit_end"`
	Header      string            `yaml:"header"`
	AttachAbove []string          `yaml:"attach_above"`
	Actions     map[string]string `yaml:"actions"`
}

// /sedum:setType
// sedum:setConst {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// RuleAdoptInvalid is the rule for a malformed adopt: block in sedum.yaml.\nconst RuleAdoptInvalid = \"adopt_invalid\"","file":"internal/genpkg/adopt.go","name":"RuleAdoptInvalid"}}
// RuleAdoptInvalid is the rule for a malformed adopt: block in sedum.yaml.
const RuleAdoptInvalid = "adopt_invalid"

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:235dd3daaefdd4a98dd3fb2bab1e2988d47e6bc2aab63e74e24144bec61de0c5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// checkAdopt reports problems with the adopt: block of the manifest.\nfunc checkAdopt(pkg *Package, r *reporter) {\n\ta := pkg.Adopt\n\tif a == nil {\n\t\treturn\n\t}\n\thasCmd := len(a.Command) > 0\n\thasDecl := a.UnitStart != \"\"\n\tif hasCmd == hasDecl {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: set exactly one of command or unit_start\")\n\t}\n\tif hasCmd {\n\t\tif a.Command[0] == \"\" {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: command first element is empty\")\n\t\t}\n\t\tif a.UnitEnd != \"\" || a.Header != \"\" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: command cannot be combined with unit_end, header, attach_above or actions\")\n\t\t}\n\t\treturn\n\t}\n\tif !hasDecl {\n\t\treturn\n\t}\n\tre, err := regexp.Compile(a.UnitStart)\n\tif err != nil {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_start does not compile: %v\", err)\n\t} else {\n\t\tnames := map[string]bool{}\n\t\tfor _, n := range re.SubexpNames() {\n\t\t\tnames[n] = true\n\t\t}\n\t\tfor _, g := range []string{\"kind\", \"name\"} {\n\t\t\tif !names[g] {\n\t\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_start lacks the named group %q\", g)\n\t\t\t}\n\t\t}\n\t}\n\tif a.UnitEnd != \"dedent\" && a.UnitEnd != \"blank\" {\n\t\tif _, err := regexp.Compile(a.UnitEnd); err != nil {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_end does not compile: %v\", err)\n\t\t}\n\t}\n\tif _, err := regexp.Compile(a.Header); err != nil {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: header does not compile: %v\", err)\n\t}\n\tfor _, p := range a.AttachAbove {\n\t\tif _, err := regexp.Compile(p); err != nil {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: attach_above %q does not compile: %v\", p, err)\n\t\t}\n\t}\n\tif len(a.Actions) == 0 {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: actions must not be empty\")\n\t}\n\tfor kind, name := range a.Actions {\n\t\tact, ok := pkg.Actions[name]\n\t\tif !ok {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: action %q for kind %q is not an action of package %s\", name, kind, pkg.Name)\n\t\t\tcontinue\n\t\t}\n\t\tfor kw, spec := range act.Kwargs {\n\t\t\tif !spec.Required {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif kw != \"file\" && kw != \"name\" && kw != \"code\" {\n\t\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: action %q has required kwarg %q outside file, name, code\", name, kw)\n\t\t\t}\n\t\t}\n\t}\n}","file":"internal/genpkg/adopt.go","name":"checkAdopt"}}
// checkAdopt reports problems with the adopt: block of the manifest.
func checkAdopt(pkg *Package, r *reporter) {
	a := pkg.Adopt
	if a == nil {
		return
	}
	hasCmd := len(a.Command) > 0
	hasDecl := a.UnitStart != ""
	if hasCmd == hasDecl {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: set exactly one of command or unit_start")
	}
	if hasCmd {
		if a.Command[0] == "" {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: command first element is empty")
		}
		if a.UnitEnd != "" || a.Header != "" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: command cannot be combined with unit_end, header, attach_above or actions")
		}
		return
	}
	if !hasDecl {
		return
	}
	re, err := regexp.Compile(a.UnitStart)
	if err != nil {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_start does not compile: %v", err)
	} else {
		names := map[string]bool{}
		for _, n := range re.SubexpNames() {
			names[n] = true
		}
		for _, g := range []string{"kind", "name"} {
			if !names[g] {
				r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_start lacks the named group %q", g)
			}
		}
	}
	if a.UnitEnd != "dedent" && a.UnitEnd != "blank" {
		if _, err := regexp.Compile(a.UnitEnd); err != nil {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_end does not compile: %v", err)
		}
	}
	if _, err := regexp.Compile(a.Header); err != nil {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: header does not compile: %v", err)
	}
	for _, p := range a.AttachAbove {
		if _, err := regexp.Compile(p); err != nil {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: attach_above %q does not compile: %v", p, err)
		}
	}
	if len(a.Actions) == 0 {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: actions must not be empty")
	}
	for kind, name := range a.Actions {
		act, ok := pkg.Actions[name]
		if !ok {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: action %q for kind %q is not an action of package %s", name, kind, pkg.Name)
			continue
		}
		for kw, spec := range act.Kwargs {
			if !spec.Required {
				continue
			}
			if kw != "file" && kw != "name" && kw != "code" {
				r.errorf(manifestFile, RuleAdoptInvalid, "adopt: action %q has required kwarg %q outside file, name, code", name, kw)
			}
		}
	}
}

// /sedum:setFunction
