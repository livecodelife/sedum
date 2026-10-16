package builtin

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"archive/zip"}}
import "archive/zip"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/builtin/builtin.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Names returns the sorted names of the built-in packages.\nfunc Names() []string {\n\tentries, err := packagesFS.ReadDir(\"packages\")\n\tif err != nil {\n\t\treturn nil\n\t}\n\tvar names []string\n\tfor _, e := range entries {\n\t\tname := e.Name()\n\t\tif strings.HasSuffix(name, \".zip\") {\n\t\t\tnames = append(names, strings.TrimSuffix(name, \".zip\"))\n\t\t}\n\t}\n\tsort.Strings(names)\n\treturn names\n}","file":"internal/builtin/builtin.go","name":"Names"}}
// Names returns the sorted names of the built-in packages.
func Names() []string {
	entries, err := packagesFS.ReadDir("packages")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".zip") {
			names = append(names, strings.TrimSuffix(name, ".zip"))
		}
	}
	sort.Strings(names)
	return names
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Files returns every file of the named built-in package, keyed by its\n// slash-separated path relative to the package root.\nfunc Files(name string) (map[string][]byte, error) {\n\tknown := Names()\n\tfound := false\n\tfor _, n := range known {\n\t\tif n == name {\n\t\t\tfound = true\n\t\t\tbreak\n\t\t}\n\t}\n\tif !found {\n\t\treturn nil, fmt.Errorf(\"no built-in package %q; built-in packages: %s\", name, strings.Join(known, \", \"))\n\t}\n\tdata, err := packagesFS.ReadFile(\"packages/\" + name + \".zip\")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tzr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tfiles := make(map[string][]byte)\n\tfor _, f := range zr.File {\n\t\tif strings.HasPrefix(f.Name, \"/\") || strings.Contains(f.Name, \"..\") || strings.HasSuffix(f.Name, \"/\") {\n\t\t\treturn nil, fmt.Errorf(\"invalid path %q in built-in package %q\", f.Name, name)\n\t\t}\n\t\trc, err := f.Open()\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tcontent, err := io.ReadAll(rc)\n\t\trc.Close()\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tfiles[f.Name] = content\n\t}\n\treturn files, nil\n}","file":"internal/builtin/builtin.go","name":"Files"}}
// Files returns every file of the named built-in package, keyed by its
// slash-separated path relative to the package root.
func Files(name string) (map[string][]byte, error) {
	known := Names()
	found := false
	for _, n := range known {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("no built-in package %q; built-in packages: %s", name, strings.Join(known, ", "))
	}
	data, err := packagesFS.ReadFile("packages/" + name + ".zip")
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	files := make(map[string][]byte)
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "/") || strings.Contains(f.Name, "..") || strings.HasSuffix(f.Name, "/") {
			return nil, fmt.Errorf("invalid path %q in built-in package %q", f.Name, name)
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		files[f.Name] = content
	}
	return files, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Paths returns the keys of files sorted.\nfunc Paths(files map[string][]byte) []string {\n\tpaths := make([]string, 0, len(files))\n\tfor p := range files {\n\t\tpaths = append(paths, p)\n\t}\n\tsort.Strings(paths)\n\treturn paths\n}","file":"internal/builtin/builtin.go","name":"Paths"}}
// Paths returns the keys of files sorted.
func Paths(files map[string][]byte) []string {
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// /sedum:setFunction
