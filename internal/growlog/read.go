package growlog

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read.go","path":"bufio"}}
import "bufio"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read.go","path":"encoding/json"}}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/growlog/read.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Read returns the entries of the grow log at path in file order, and the\n// number of malformed lines it skipped. Blank lines are ignored. A missing file\n// is an error.\nfunc Read(path string) ([]Entry, int, error) {\n\tf, err := os.Open(path)\n\tif err != nil {\n\t\treturn nil, 0, fmt.Errorf(\"open grow log: %w\", err)\n\t}\n\tdefer f.Close()\n\tvar entries []Entry\n\tmalformed := 0\n\tsc := bufio.NewScanner(f)\n\tsc.Buffer(make([]byte, 0, 64*1024), 1024*1024)\n\tfor sc.Scan() {\n\t\tline := strings.TrimSpace(sc.Text())\n\t\tif line == \"\" {\n\t\t\tcontinue\n\t\t}\n\t\tvar e Entry\n\t\tif err := json.Unmarshal([]byte(line), &e); err != nil {\n\t\t\tmalformed++\n\t\t\tcontinue\n\t\t}\n\t\tentries = append(entries, e)\n\t}\n\tif err := sc.Err(); err != nil {\n\t\treturn entries, malformed, fmt.Errorf(\"read grow log %s: %w\", path, err)\n\t}\n\treturn entries, malformed, nil\n}","file":"internal/growlog/read.go","name":"Read"}}
// Read returns the entries of the grow log at path in file order, and the
// number of malformed lines it skipped. Blank lines are ignored. A missing file
// is an error.
func Read(path string) ([]Entry, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, fmt.Errorf("open grow log: %w", err)
	}
	defer f.Close()
	var entries []Entry
	malformed := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			malformed++
			continue
		}
		entries = append(entries, e)
	}
	if err := sc.Err(); err != nil {
		return entries, malformed, fmt.Errorf("read grow log %s: %w", path, err)
	}
	return entries, malformed, nil
}

// /sedum:setFunction
