package builtin

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"archive/zip"},"kwargs":"sha256:1f57c332afa6b5b8d53944798e8d8898a43bb91d7c9480ad2c7226b55b01a41c","body":"sha256:d313d2d9a0178875efd7a917bc5a898ac82d6ac825476f6784d6a74e7e61b8e1"}
import "archive/zip"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"bytes"},"kwargs":"sha256:20d7ebfa4df60fa6df4dbf6cbb66406789614ccb961248f47ca66a8a7435f190","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"fmt"},"kwargs":"sha256:e717a6641c9566ea079c737695ab562492d086c218c6ced72acdf1ef2b32d0c6","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"io"},"kwargs":"sha256:39eb732aec16ea55f76b54e3019846aaca9eee727edf049c60be6884f01d69a7","body":"sha256:ef8e8abf4ef5ee40a19dc85b49bf09c5df1fe32abfe834e0ab4c2175baf80ddf"}
import "io"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"sort"},"kwargs":"sha256:80e69c5246d01ab9f9a787021d67acceabe7d6b481000363167f9d5abf3a271b","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/builtin/builtin.go","path":"strings"},"kwargs":"sha256:685b043f96c25a8e437c35d5484152728de810f28d56cb2b5a85c40d7583fe94","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin.go","name":"Names"},"kwargs":"sha256:8eacddfe0b33a8343007d8bf8231a5c453a5921d2d2c2a9300bda8a2bccd4f2e","body":"sha256:acf49b2d695864fd06b8b92e380eb37988f60a04dd6c18e22bc7655807bd066e"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin.go","name":"Files"},"kwargs":"sha256:21f0aec1df9caee7ad8e4f4169c61064b5e964b84a3de81ed24db2e26362efb6","body":"sha256:267311fb5dabc97358a0113bea618f53a36a9aefeab1d7df412f446cc95cc4a2"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-fc6b6621","gen":"sha256:5e3114492de2056ad3dbbed3f453b183f17b7bfead4e52e181503be5cb4b12ff","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/builtin/builtin.go","name":"Paths"},"kwargs":"sha256:f4944cbbe3d4e77ed01ba889eb5f3aabf5654f933509ae1d24bf6eb3d7318672","body":"sha256:8e6ed6c2e2fa29f55d0734526c7da73e3b9853239a891737021cba96a11495ed"}
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
