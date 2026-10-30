package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read.go","path":"bufio"},"kwargs":"sha256:1eaff5ed057c8eabd0c108e79f3c345e16d69ab59c18e6878f4eee6fde221a2d","body":"sha256:f392686aa91926837d96af56e8fce4f794d441bdf829d8efba66160f33dfd873"}
import "bufio"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read.go","path":"encoding/json"},"kwargs":"sha256:da03f4d700a7e768a17b28f6a912ded5c245e5be085aab6bcd6c496e08d180a8","body":"sha256:872d7c3e9dee608e79a7f8ef57930d24b5e30f2a57d891a52549179ace45946a"}
import "encoding/json"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read.go","path":"fmt"},"kwargs":"sha256:c8b2ed9a66d54cae49c3a519bbba64a3bf9a76c8e9c6130bb0d3e4e6c2da9bba","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read.go","path":"os"},"kwargs":"sha256:972aeb2a17754eb701d8e030ce82abb17de85c38e638fcf3eb9ad4129c2a6046","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/read.go","path":"strings"},"kwargs":"sha256:9ad6433b166217d72a9c09238a7e4490d50a894c4c9b05c2f431309251b56e66","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2f2f6e70","gen":"sha256:d25ad2a156c461218ef01ef0f757aeb7da6b411cceba9cf0fd04cfb042ca0771","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/read.go","name":"Read"},"kwargs":"sha256:48d04cbe7db6a4d2f9aeefa1a53a5533037f19a03c291f019eb86564e69aef39","body":"sha256:7b2164d67533c29a750b302d70adc261df11971b598348b6fc82f9cb2a1e8cfa"}
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
