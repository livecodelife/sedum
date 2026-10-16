package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"path"}}
import "path"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"github.com/livecodelife/sedum/internal/pathpat"}}
import "github.com/livecodelife/sedum/internal/pathpat"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/select.go","path":"github.com/livecodelife/sedum/internal/resolve"}}
import "github.com/livecodelife/sedum/internal/resolve"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Reason says why Select skipped a file.\ntype Reason string","file":"internal/adopt/select.go","name":"Reason"}}
// Reason says why Select skipped a file.
type Reason string

// /sedum:setType
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonNoPackage means no package claims the file's extension.\nconst ReasonNoPackage Reason = \"no_package\"","file":"internal/adopt/select.go","name":"ReasonNoPackage"}}
// ReasonNoPackage means no package claims the file's extension.
const ReasonNoPackage Reason = "no_package"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonUnmanaged means the file is unmanaged by its package.\nconst ReasonUnmanaged Reason = \"unmanaged\"","file":"internal/adopt/select.go","name":"ReasonUnmanaged"}}
// ReasonUnmanaged means the file is unmanaged by its package.
const ReasonUnmanaged Reason = "unmanaged"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonTest means the file is a test path and tests were not requested.\nconst ReasonTest Reason = \"test\"","file":"internal/adopt/select.go","name":"ReasonTest"}}
// ReasonTest means the file is a test path and tests were not requested.
const ReasonTest Reason = "test"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonGenerated means the file's header marks it as generated.\nconst ReasonGenerated Reason = \"generated\"","file":"internal/adopt/select.go","name":"ReasonGenerated"}}
// ReasonGenerated means the file's header marks it as generated.
const ReasonGenerated Reason = "generated"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonAlreadyAdopted means the file already carries sedum markers.\nconst ReasonAlreadyAdopted Reason = \"already_adopted\"","file":"internal/adopt/select.go","name":"ReasonAlreadyAdopted"}}
// ReasonAlreadyAdopted means the file already carries sedum markers.
const ReasonAlreadyAdopted Reason = "already_adopted"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonNoTemplate means resolution produced no template.\nconst ReasonNoTemplate Reason = \"no_template\"","file":"internal/adopt/select.go","name":"ReasonNoTemplate"}}
// ReasonNoTemplate means resolution produced no template.
const ReasonNoTemplate Reason = "no_template"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonNoImprint means no imprint record authorizes the file.\nconst ReasonNoImprint Reason = \"no_imprint\"","file":"internal/adopt/select.go","name":"ReasonNoImprint"}}
// ReasonNoImprint means no imprint record authorizes the file.
const ReasonNoImprint Reason = "no_imprint"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonAmbiguousImprint means more than one imprint record authorizes the file.\nconst ReasonAmbiguousImprint Reason = \"ambiguous_imprint\"","file":"internal/adopt/select.go","name":"ReasonAmbiguousImprint"}}
// ReasonAmbiguousImprint means more than one imprint record authorizes the file.
const ReasonAmbiguousImprint Reason = "ambiguous_imprint"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonHeader means the file was skipped as a header span.\nconst ReasonHeader Reason = \"header\"","file":"internal/adopt/select.go","name":"ReasonHeader"}}
// ReasonHeader means the file was skipped as a header span.
const ReasonHeader Reason = "header"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonFloatingComment means the file has a floating comment.\nconst ReasonFloatingComment Reason = \"floating_comment\"","file":"internal/adopt/select.go","name":"ReasonFloatingComment"}}
// ReasonFloatingComment means the file has a floating comment.
const ReasonFloatingComment Reason = "floating_comment"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonOtherResidual means the file has non-whitespace residual text.\nconst ReasonOtherResidual Reason = \"other_residual\"","file":"internal/adopt/select.go","name":"ReasonOtherResidual"}}
// ReasonOtherResidual means the file has non-whitespace residual text.
const ReasonOtherResidual Reason = "other_residual"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ReasonUnresolved means resolution failed for the file.\nconst ReasonUnresolved Reason = \"unresolved\"","file":"internal/adopt/select.go","name":"ReasonUnresolved"}}
// ReasonUnresolved means resolution failed for the file.
const ReasonUnresolved Reason = "unresolved"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Skip records one file that Select did not adopt, and why.\ntype Skip struct {\n\tPath   string\n\tReason Reason\n\tDetail string\n}","file":"internal/adopt/select.go","name":"Skip"}}
// Skip records one file that Select did not adopt, and why.
type Skip struct {
	Path   string
	Reason Reason
	Detail string
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Candidate is one file Select chose for adoption.\ntype Candidate struct {\n\tPath     string\n\tPackage  *genpkg.Package\n\tTemplate string\n\tCaptures map[string]string\n\tRecord   *record.Record\n\tSrc      []byte\n}","file":"internal/adopt/select.go","name":"Candidate"}}
// Candidate is one file Select chose for adoption.
type Candidate struct {
	Path     string
	Package  *genpkg.Package
	Template string
	Captures map[string]string
	Record   *record.Record
	Src      []byte
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// SelectConfig holds the inputs to Select.\ntype SelectConfig struct {\n\tRoot    string\n\tSet     *genpkg.Set\n\tRecords *record.Set\n\tTests   bool\n}","file":"internal/adopt/select.go","name":"SelectConfig"}}
// SelectConfig holds the inputs to Select.
type SelectConfig struct {
	Root    string
	Set     *genpkg.Set
	Records *record.Set
	Tests   bool
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-37216199","gen":"sha256:abd8e1b2415a66e4ae979ddcf465e28126cbec7e6cc481d2dc1830782039e90d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Select chooses the files under cfg.Root that adoption can bring under\n// Sedum, and records a Skip with its reason for every other file it visits.\n// It never writes.\nfunc Select(cfg SelectConfig, paths []string) ([]Candidate, []Skip, error) {\n\tvar cands []Candidate\n\tvar skips []Skip\n\tseen := map[string]bool{}\n\tskip := func(p string, r Reason, detail string) {\n\t\tskips = append(skips, Skip{Path: p, Reason: r, Detail: detail})\n\t}\n\tvar visit func(rel string) error\n\tvisit = func(rel string) error {\n\t\tabs := filepath.Join(cfg.Root, filepath.FromSlash(rel))\n\t\tinfo, err := os.Stat(abs)\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"stat %s: %w\", rel, err)\n\t\t}\n\t\tif info.IsDir() {\n\t\t\tentries, err := os.ReadDir(abs)\n\t\t\tif err != nil {\n\t\t\t\treturn fmt.Errorf(\"read dir %s: %w\", rel, err)\n\t\t\t}\n\t\t\tfor _, e := range entries {\n\t\t\t\tname := e.Name()\n\t\t\t\tif e.IsDir() && strings.HasPrefix(name, \".\") {\n\t\t\t\t\tcontinue\n\t\t\t\t}\n\t\t\t\tif err := visit(path.Join(rel, name)); err != nil {\n\t\t\t\t\treturn err\n\t\t\t\t}\n\t\t\t}\n\t\t\treturn nil\n\t\t}\n\t\tif !info.Mode().IsRegular() || seen[rel] {\n\t\t\treturn nil\n\t\t}\n\t\tseen[rel] = true\n\t\tpkgs := cfg.Set.ForExtension(filepath.Ext(rel))\n\t\tif len(pkgs) == 0 {\n\t\t\tskip(rel, ReasonNoPackage, \"\")\n\t\t\treturn nil\n\t\t}\n\t\tif pkg, entry, ok := cfg.Set.Unmanaged(rel); ok {\n\t\t\tskip(rel, ReasonUnmanaged, pkg+\" \"+entry)\n\t\t\treturn nil\n\t\t}\n\t\tpkg := pkgs[0]\n\t\tsrc, err := os.ReadFile(abs)\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"read %s: %w\", rel, err)\n\t\t}\n\t\tif !cfg.Tests {\n\t\t\tif _, ok := pathpat.MatchAny(pkg.TestPaths, rel); ok {\n\t\t\t\tskip(rel, ReasonTest, \"\")\n\t\t\t\treturn nil\n\t\t\t}\n\t\t}\n\t\thead := strings.SplitN(string(src), \"\\n\", 21)\n\t\tif len(head) > 20 {\n\t\t\thead = head[:20]\n\t\t}\n\t\tfor _, line := range head {\n\t\t\tif strings.Contains(line, \"DO NOT EDIT\") || strings.Contains(line, \"@generated\") {\n\t\t\t\tskip(rel, ReasonGenerated, \"\")\n\t\t\t\treturn nil\n\t\t\t}\n\t\t}\n\t\tif strings.Contains(string(src), pkg.Comment().Prefix+\" sedum:\") {\n\t\t\tskip(rel, ReasonAlreadyAdopted, \"\")\n\t\t\treturn nil\n\t\t}\n\t\tvar imprints []*record.Record\n\t\tfor _, r := range cfg.Records.Records {\n\t\t\tif r.Deprecated || r.Type != \"imprint\" || !r.Authorizes(rel) {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif _, forbidden := r.Forbids(rel); forbidden {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\timprints = append(imprints, r)\n\t\t}\n\t\tswitch {\n\t\tcase len(imprints) == 0:\n\t\t\tskip(rel, ReasonNoImprint, \"\")\n\t\t\treturn nil\n\t\tcase len(imprints) > 1:\n\t\t\tids := make([]string, 0, len(imprints))\n\t\t\tfor _, r := range imprints {\n\t\t\t\tids = append(ids, r.ID)\n\t\t\t}\n\t\t\tskip(rel, ReasonAmbiguousImprint, strings.Join(ids, \", \"))\n\t\t\treturn nil\n\t\t}\n\t\timp := imprints[0]\n\t\tres, _, err := resolve.Paths(cfg.Set, &record.Set{Records: []*record.Record{{ID: imp.ID, Paths: []string{rel}}}}, nil)\n\t\tif err != nil {\n\t\t\tskip(rel, ReasonUnresolved, err.Error())\n\t\t\treturn nil\n\t\t}\n\t\tif len(res) == 0 {\n\t\t\tskip(rel, ReasonUnresolved, \"no resolution\")\n\t\t\treturn nil\n\t\t}\n\t\tr0 := res[0]\n\t\tif r0.Unmanaged {\n\t\t\tskip(rel, ReasonUnmanaged, \"\")\n\t\t\treturn nil\n\t\t}\n\t\tif r0.Template == \"\" {\n\t\t\tskip(rel, ReasonNoTemplate, \"\")\n\t\t\treturn nil\n\t\t}\n\t\tcands = append(cands, Candidate{Path: rel, Package: r0.Package, Template: r0.Template, Captures: r0.Captures, Record: imp, Src: src})\n\t\treturn nil\n\t}\n\tfor _, p := range paths {\n\t\tclean := path.Clean(filepath.ToSlash(p))\n\t\tif clean == \"..\" || strings.HasPrefix(clean, \"../\") || path.IsAbs(clean) {\n\t\t\treturn nil, nil, fmt.Errorf(\"path %q escapes root\", p)\n\t\t}\n\t\tif err := visit(clean); err != nil {\n\t\t\treturn nil, nil, err\n\t\t}\n\t}\n\treturn cands, skips, nil\n}","file":"internal/adopt/select.go","name":"Select"}}
// Select chooses the files under cfg.Root that adoption can bring under
// Sedum, and records a Skip with its reason for every other file it visits.
// It never writes.
func Select(cfg SelectConfig, paths []string) ([]Candidate, []Skip, error) {
	var cands []Candidate
	var skips []Skip
	seen := map[string]bool{}
	skip := func(p string, r Reason, detail string) {
		skips = append(skips, Skip{Path: p, Reason: r, Detail: detail})
	}
	var visit func(rel string) error
	visit = func(rel string) error {
		abs := filepath.Join(cfg.Root, filepath.FromSlash(rel))
		info, err := os.Stat(abs)
		if err != nil {
			return fmt.Errorf("stat %s: %w", rel, err)
		}
		if info.IsDir() {
			entries, err := os.ReadDir(abs)
			if err != nil {
				return fmt.Errorf("read dir %s: %w", rel, err)
			}
			for _, e := range entries {
				name := e.Name()
				if e.IsDir() && strings.HasPrefix(name, ".") {
					continue
				}
				if err := visit(path.Join(rel, name)); err != nil {
					return err
				}
			}
			return nil
		}
		if !info.Mode().IsRegular() || seen[rel] {
			return nil
		}
		seen[rel] = true
		pkgs := cfg.Set.ForExtension(filepath.Ext(rel))
		if len(pkgs) == 0 {
			skip(rel, ReasonNoPackage, "")
			return nil
		}
		if pkg, entry, ok := cfg.Set.Unmanaged(rel); ok {
			skip(rel, ReasonUnmanaged, pkg+" "+entry)
			return nil
		}
		pkg := pkgs[0]
		src, err := os.ReadFile(abs)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		if !cfg.Tests {
			if _, ok := pathpat.MatchAny(pkg.TestPaths, rel); ok {
				skip(rel, ReasonTest, "")
				return nil
			}
		}
		head := strings.SplitN(string(src), "\n", 21)
		if len(head) > 20 {
			head = head[:20]
		}
		for _, line := range head {
			if strings.Contains(line, "DO NOT EDIT") || strings.Contains(line, "@generated") {
				skip(rel, ReasonGenerated, "")
				return nil
			}
		}
		if strings.Contains(string(src), pkg.Comment().Prefix+" sedum:") {
			skip(rel, ReasonAlreadyAdopted, "")
			return nil
		}
		var imprints []*record.Record
		for _, r := range cfg.Records.Records {
			if r.Deprecated || r.Type != "imprint" || !r.Authorizes(rel) {
				continue
			}
			if _, forbidden := r.Forbids(rel); forbidden {
				continue
			}
			imprints = append(imprints, r)
		}
		switch {
		case len(imprints) == 0:
			skip(rel, ReasonNoImprint, "")
			return nil
		case len(imprints) > 1:
			ids := make([]string, 0, len(imprints))
			for _, r := range imprints {
				ids = append(ids, r.ID)
			}
			skip(rel, ReasonAmbiguousImprint, strings.Join(ids, ", "))
			return nil
		}
		imp := imprints[0]
		res, _, err := resolve.Paths(cfg.Set, &record.Set{Records: []*record.Record{{ID: imp.ID, Paths: []string{rel}}}}, nil)
		if err != nil {
			skip(rel, ReasonUnresolved, err.Error())
			return nil
		}
		if len(res) == 0 {
			skip(rel, ReasonUnresolved, "no resolution")
			return nil
		}
		r0 := res[0]
		if r0.Unmanaged {
			skip(rel, ReasonUnmanaged, "")
			return nil
		}
		if r0.Template == "" {
			skip(rel, ReasonNoTemplate, "")
			return nil
		}
		cands = append(cands, Candidate{Path: rel, Package: r0.Package, Template: r0.Template, Captures: r0.Captures, Record: imp, Src: src})
		return nil
	}
	for _, p := range paths {
		clean := path.Clean(filepath.ToSlash(p))
		if clean == ".." || strings.HasPrefix(clean, "../") || path.IsAbs(clean) {
			return nil, nil, fmt.Errorf("path %q escapes root", p)
		}
		if err := visit(clean); err != nil {
			return nil, nil, err
		}
	}
	return cands, skips, nil
}

// /sedum:setFunction
