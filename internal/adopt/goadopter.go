package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"go/ast"}}
import "go/ast"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"go/format"}}
import "go/format"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"go/parser"}}
import "go/parser"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"go/token"}}
import "go/token"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"strconv"}}
import "strconv"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// GoAdopter is the built-in adopter for Go source (adopt: {builtin: go}). It\n// needs gofmt-clean input and splits a file with go/parser into a header, one\n// unit per declaration, and residual spans for everything else.\ntype GoAdopter struct{}","file":"internal/adopt/goadopter.go","name":"GoAdopter"}}
// GoAdopter is the built-in adopter for Go source (adopt: {builtin: go}). It
// needs gofmt-clean input and splits a file with go/parser into a header, one
// unit per declaration, and residual spans for everything else.
type GoAdopter struct{}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// goRecvName returns the receiver's base type name: no star, parentheses or type parameters.\nfunc goRecvName(e ast.Expr) string {\n\tswitch t := e.(type) {\n\tcase *ast.StarExpr:\n\t\treturn goRecvName(t.X)\n\tcase *ast.ParenExpr:\n\t\treturn goRecvName(t.X)\n\tcase *ast.IndexExpr:\n\t\treturn goRecvName(t.X)\n\tcase *ast.IndexListExpr:\n\t\treturn goRecvName(t.X)\n\tcase *ast.Ident:\n\t\treturn t.Name\n\t}\n\treturn \"\"\n}","file":"internal/adopt/goadopter.go","name":"goRecvName"}}
// goRecvName returns the receiver's base type name: no star, parentheses or type parameters.
func goRecvName(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return goRecvName(t.X)
	case *ast.ParenExpr:
		return goRecvName(t.X)
	case *ast.IndexExpr:
		return goRecvName(t.X)
	case *ast.IndexListExpr:
		return goRecvName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// goGenName returns the name of a type, var or const declaration: the first\n// name of its first spec. It returns \"\" for a declaration with no specs.\nfunc goGenName(d *ast.GenDecl) string {\n\tif len(d.Specs) == 0 {\n\t\treturn \"\"\n\t}\n\tswitch s := d.Specs[0].(type) {\n\tcase *ast.TypeSpec:\n\t\treturn s.Name.Name\n\tcase *ast.ValueSpec:\n\t\treturn s.Names[0].Name\n\t}\n\treturn \"\"\n}","file":"internal/adopt/goadopter.go","name":"goGenName"}}
// goGenName returns the name of a type, var or const declaration: the first
// name of its first spec. It returns "" for a declaration with no specs.
func goGenName(d *ast.GenDecl) string {
	if len(d.Specs) == 0 {
		return ""
	}
	switch s := d.Specs[0].(type) {
	case *ast.TypeSpec:
		return s.Name.Name
	case *ast.ValueSpec:
		return s.Names[0].Name
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// goLineEnd returns the offset just past the line that contains offset off.\nfunc goLineEnd(src []byte, off int) int {\n\tif i := bytes.IndexByte(src[off:], '\\n'); i >= 0 {\n\t\treturn off + i + 1\n\t}\n\treturn len(src)\n}","file":"internal/adopt/goadopter.go","name":"goLineEnd"}}
// goLineEnd returns the offset just past the line that contains offset off.
func goLineEnd(src []byte, off int) int {
	if i := bytes.IndexByte(src[off:], '\n'); i >= 0 {
		return off + i + 1
	}
	return len(src)
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// goPlainImport reports whether d is one ungrouped, unaliased, uncommented\n// import that is not cgo, and returns its path.\nfunc goPlainImport(d *ast.GenDecl) (string, bool) {\n\tif d.Lparen.IsValid() || len(d.Specs) != 1 || d.Doc != nil {\n\t\treturn \"\", false\n\t}\n\ts := d.Specs[0].(*ast.ImportSpec)\n\tif s.Name != nil || s.Doc != nil || s.Comment != nil {\n\t\treturn \"\", false\n\t}\n\tp, err := strconv.Unquote(s.Path.Value)\n\tif err != nil || p == \"C\" {\n\t\treturn \"\", false\n\t}\n\treturn p, true\n}","file":"internal/adopt/goadopter.go","name":"goPlainImport"}}
// goPlainImport reports whether d is one ungrouped, unaliased, uncommented
// import that is not cgo, and returns its path.
func goPlainImport(d *ast.GenDecl) (string, bool) {
	if d.Lparen.IsValid() || len(d.Specs) != 1 || d.Doc != nil {
		return "", false
	}
	s := d.Specs[0].(*ast.ImportSpec)
	if s.Name != nil || s.Doc != nil || s.Comment != nil {
		return "", false
	}
	p, err := strconv.Unquote(s.Path.Value)
	if err != nil || p == "C" {
		return "", false
	}
	return p, true
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:7f08b8faeeb76720edcd889b195a10c494f118783a6b47c71f35e954ccbbc0bb","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Split implements Adopter.\nfunc (GoAdopter) Split(path string, src []byte) ([]Span, error) {\n\tformatted, err := format.Source(src)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"parse %s: %w\", path, err)\n\t}\n\tif !bytes.Equal(formatted, src) {\n\t\treturn nil, fmt.Errorf(\"%s: not gofmt-clean\", path)\n\t}\n\tfset := token.NewFileSet()\n\tfile, err := parser.ParseFile(fset, path, src, parser.ParseComments)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"parse %s: %w\", path, err)\n\t}\n\toff := func(p token.Pos) int { return fset.Position(p).Offset }\n\tvar spans []Span\n\tpos := 0\n\tgap := func(to int) {\n\t\tif to > pos {\n\t\t\tspans = append(spans, Span{Kind: KindResidual, Text: string(src[pos:to])})\n\t\t\tpos = to\n\t\t}\n\t}\n\tgap(off(file.Package))\n\theaderEnd := goLineEnd(src, off(file.Name.End()))\n\tspans = append(spans, Span{Kind: KindHeader, Text: string(src[pos:headerEnd])})\n\tpos = headerEnd\n\tseen := map[string]int{}\n\timports := 0\n\tfor _, decl := range file.Decls {\n\t\tstart := off(decl.Pos())\n\t\tvar doc *ast.CommentGroup\n\t\tvar action, name string\n\t\tkwargs := map[string]any{\"file\": path}\n\t\tswitch d := decl.(type) {\n\t\tcase *ast.FuncDecl:\n\t\t\tdoc = d.Doc\n\t\t\tname = d.Name.Name\n\t\t\taction = \"setFunction\"\n\t\t\tif d.Recv != nil && len(d.Recv.List) == 1 {\n\t\t\t\taction = \"setMethod\"\n\t\t\t\tkwargs[\"receiver\"] = goRecvName(d.Recv.List[0].Type)\n\t\t\t}\n\t\t\tif d.Recv == nil && name == \"init\" {\n\t\t\t\tseen[\"init\"]++\n\t\t\t\tname = fmt.Sprintf(\"init#%d\", seen[\"init\"])\n\t\t\t}\n\t\tcase *ast.GenDecl:\n\t\t\tdoc = d.Doc\n\t\t\tswitch d.Tok {\n\t\t\tcase token.IMPORT:\n\t\t\t\tif p, ok := goPlainImport(d); ok {\n\t\t\t\t\taction = \"addImport\"\n\t\t\t\t\tkwargs[\"path\"] = p\n\t\t\t\t} else if len(d.Specs) == 1 && d.Specs[0].(*ast.ImportSpec).Path.Value == `\"C\"` {\n\t\t\t\t\taction = \"\"\n\t\t\t\t} else {\n\t\t\t\t\timports++\n\t\t\t\t\taction = \"setImports\"\n\t\t\t\t\tname = fmt.Sprintf(\"imports#%d\", imports)\n\t\t\t\t}\n\t\t\tcase token.TYPE:\n\t\t\t\taction, name = \"setType\", goGenName(d)\n\t\t\tcase token.VAR:\n\t\t\t\taction, name = \"setVar\", goGenName(d)\n\t\t\tcase token.CONST:\n\t\t\t\taction, name = \"setConst\", goGenName(d)\n\t\t\t}\n\t\t\tif d.Tok != token.IMPORT && name == \"\" {\n\t\t\t\taction = \"\"\n\t\t\t}\n\t\t}\n\t\tif doc != nil {\n\t\t\tstart = off(doc.Pos())\n\t\t}\n\t\tend := goLineEnd(src, off(decl.End()))\n\t\tgap(start)\n\t\ttext := string(src[start:end])\n\t\tif action == \"\" {\n\t\t\tspans = append(spans, Span{Kind: KindResidual, Text: text})\n\t\t\tpos = end\n\t\t\tcontinue\n\t\t}\n\t\tif action != \"addImport\" && !(action == \"setFunction\" && len(name) > 5 && name[:5] == \"init#\") {\n\t\t\tkey := action + \"/\" + name\n\t\t\tif r, ok := kwargs[\"receiver\"].(string); ok {\n\t\t\t\tkey += \"/\" + r\n\t\t\t}\n\t\t\tseen[key]++\n\t\t\tif seen[key] > 1 {\n\t\t\t\tname = fmt.Sprintf(\"%s#%d\", name, seen[key])\n\t\t\t}\n\t\t}\n\t\tif action != \"addImport\" {\n\t\t\tkwargs[\"name\"] = name\n\t\t}\n\t\tif action != \"addImport\" {\n\t\t\tkwargs[\"code\"] = strings.TrimSuffix(text, \"\\n\")\n\t\t}\n\t\tspans = append(spans, Span{Kind: KindUnit, Text: text, Action: action, Kwargs: kwargs})\n\t\tpos = end\n\t}\n\tgap(len(src))\n\treturn spans, nil\n}","file":"internal/adopt/goadopter.go","name":"Split","receiver":"GoAdopter"}}
// Split implements Adopter.
func (GoAdopter) Split(path string, src []byte) ([]Span, error) {
	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if !bytes.Equal(formatted, src) {
		return nil, fmt.Errorf("%s: not gofmt-clean", path)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	off := func(p token.Pos) int { return fset.Position(p).Offset }
	var spans []Span
	pos := 0
	gap := func(to int) {
		if to > pos {
			spans = append(spans, Span{Kind: KindResidual, Text: string(src[pos:to])})
			pos = to
		}
	}
	gap(off(file.Package))
	headerEnd := goLineEnd(src, off(file.Name.End()))
	spans = append(spans, Span{Kind: KindHeader, Text: string(src[pos:headerEnd])})
	pos = headerEnd
	seen := map[string]int{}
	imports := 0
	for _, decl := range file.Decls {
		start := off(decl.Pos())
		var doc *ast.CommentGroup
		var action, name string
		kwargs := map[string]any{"file": path}
		switch d := decl.(type) {
		case *ast.FuncDecl:
			doc = d.Doc
			name = d.Name.Name
			action = "setFunction"
			if d.Recv != nil && len(d.Recv.List) == 1 {
				action = "setMethod"
				kwargs["receiver"] = goRecvName(d.Recv.List[0].Type)
			}
			if d.Recv == nil && name == "init" {
				seen["init"]++
				name = fmt.Sprintf("init#%d", seen["init"])
			}
		case *ast.GenDecl:
			doc = d.Doc
			switch d.Tok {
			case token.IMPORT:
				if p, ok := goPlainImport(d); ok {
					action = "addImport"
					kwargs["path"] = p
				} else if len(d.Specs) == 1 && d.Specs[0].(*ast.ImportSpec).Path.Value == `"C"` {
					action = ""
				} else {
					imports++
					action = "setImports"
					name = fmt.Sprintf("imports#%d", imports)
				}
			case token.TYPE:
				action, name = "setType", goGenName(d)
			case token.VAR:
				action, name = "setVar", goGenName(d)
			case token.CONST:
				action, name = "setConst", goGenName(d)
			}
			if d.Tok != token.IMPORT && name == "" {
				action = ""
			}
		}
		if doc != nil {
			start = off(doc.Pos())
		}
		end := goLineEnd(src, off(decl.End()))
		gap(start)
		text := string(src[start:end])
		if action == "" {
			spans = append(spans, Span{Kind: KindResidual, Text: text})
			pos = end
			continue
		}
		if action != "addImport" && !(action == "setFunction" && len(name) > 5 && name[:5] == "init#") {
			key := action + "/" + name
			if r, ok := kwargs["receiver"].(string); ok {
				key += "/" + r
			}
			seen[key]++
			if seen[key] > 1 {
				name = fmt.Sprintf("%s#%d", name, seen[key])
			}
		}
		if action != "addImport" {
			kwargs["name"] = name
		}
		if action != "addImport" {
			kwargs["code"] = strings.TrimSuffix(text, "\n")
		}
		spans = append(spans, Span{Kind: KindUnit, Text: text, Action: action, Kwargs: kwargs})
		pos = end
	}
	gap(len(src))
	return spans, nil
}

// /sedum:setMethod
