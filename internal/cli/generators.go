package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"sort"}}
import "sort"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/generators.go","path":"github.com/livecodelife/sedum/internal/builtin"}}
import "github.com/livecodelife/sedum/internal/builtin"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type installOptions struct {\n\tDir    string\n\tFrom   string\n\tForce  bool\n\tDryRun bool\n}","file":"internal/cli/generators.go","name":"installOptions"}}
type installOptions struct {
	Dir    string
	From   string
	Force  bool
	DryRun bool
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type installEntry struct {\n\tPath   string\n\tStatus string\n\tData   []byte\n}","file":"internal/cli/generators.go","name":"installEntry"}}
type installEntry struct {
	Path   string
	Status string
	Data   []byte
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func newGeneratorsCommand() *cobra.Command {\n\tcmd := &cobra.Command{\n\t\tUse:   \"generators\",\n\t\tShort: \"Manage generator packages\",\n\t}\n\tvar opts installOptions\n\tinstall := &cobra.Command{\n\t\tUse:   \"install <pkg>...\",\n\t\tShort: \"Install generator packages into a directory\",\n\t\tArgs:  cobra.MinimumNArgs(1),\n\t\tRunE: func(cmd *cobra.Command, args []string) error {\n\t\t\treturn applyInstall(cmd.OutOrStdout(), args, opts)\n\t\t},\n\t}\n\tinstall.Flags().StringVar(&opts.Dir, \"dir\", \"generators\", \"directory to install into\")\n\tinstall.Flags().StringVar(&opts.From, \"from\", \"\", \"copy the package from this directory\")\n\tinstall.Flags().BoolVar(&opts.Force, \"force\", false, \"overwrite files that differ\")\n\tinstall.Flags().BoolVar(&opts.DryRun, \"dry-run\", false, \"print the plan without writing\")\n\tcmd.AddCommand(install)\n\treturn cmd\n}","file":"internal/cli/generators.go","name":"newGeneratorsCommand"}}
func newGeneratorsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generators",
		Short: "Manage generator packages",
	}
	var opts installOptions
	install := &cobra.Command{
		Use:   "install <pkg>...",
		Short: "Install generator packages into a directory",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyInstall(cmd.OutOrStdout(), args, opts)
		},
	}
	install.Flags().StringVar(&opts.Dir, "dir", "generators", "directory to install into")
	install.Flags().StringVar(&opts.From, "from", "", "copy the package from this directory")
	install.Flags().BoolVar(&opts.Force, "force", false, "overwrite files that differ")
	install.Flags().BoolVar(&opts.DryRun, "dry-run", false, "print the plan without writing")
	cmd.AddCommand(install)
	return cmd
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func planInstall(names []string, opts installOptions) ([]installEntry, error) {\n\tvar entries []installEntry\n\tfor _, name := range names {\n\t\tif name == \"\" || name == \".\" || name == \"..\" || strings.ContainsAny(name, \"/\\\\\") || strings.ContainsRune(name, 0) {\n\t\t\treturn nil, fmt.Errorf(\"invalid package name %q\", name)\n\t\t}\n\t\tfiles, err := loadPackageFiles(name, opts)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\trels := make([]string, 0, len(files))\n\t\tfor rel := range files {\n\t\t\trels = append(rels, rel)\n\t\t}\n\t\tsort.Strings(rels)\n\t\tfor _, rel := range rels {\n\t\t\tif !filepath.IsLocal(rel) {\n\t\t\t\treturn nil, fmt.Errorf(\"package %q: refusing file path %q\", name, rel)\n\t\t\t}\n\t\t\tpath := filepath.Join(opts.Dir, name, rel)\n\t\t\tstatus := \"create\"\n\t\t\texisting, err := os.ReadFile(path)\n\t\t\tswitch {\n\t\t\tcase err == nil && bytes.Equal(existing, files[rel]):\n\t\t\t\tstatus = \"identical\"\n\t\t\tcase err == nil:\n\t\t\t\tstatus = \"differs\"\n\t\t\tcase !os.IsNotExist(err):\n\t\t\t\treturn nil, err\n\t\t\t}\n\t\t\tentries = append(entries, installEntry{Path: path, Status: status, Data: files[rel]})\n\t\t}\n\t}\n\treturn entries, nil\n}","file":"internal/cli/generators.go","name":"planInstall"}}
func planInstall(names []string, opts installOptions) ([]installEntry, error) {
	var entries []installEntry
	for _, name := range names {
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\") || strings.ContainsRune(name, 0) {
			return nil, fmt.Errorf("invalid package name %q", name)
		}
		files, err := loadPackageFiles(name, opts)
		if err != nil {
			return nil, err
		}
		rels := make([]string, 0, len(files))
		for rel := range files {
			rels = append(rels, rel)
		}
		sort.Strings(rels)
		for _, rel := range rels {
			if !filepath.IsLocal(rel) {
				return nil, fmt.Errorf("package %q: refusing file path %q", name, rel)
			}
			path := filepath.Join(opts.Dir, name, rel)
			status := "create"
			existing, err := os.ReadFile(path)
			switch {
			case err == nil && bytes.Equal(existing, files[rel]):
				status = "identical"
			case err == nil:
				status = "differs"
			case !os.IsNotExist(err):
				return nil, err
			}
			entries = append(entries, installEntry{Path: path, Status: status, Data: files[rel]})
		}
	}
	return entries, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func loadPackageFiles(name string, opts installOptions) (map[string][]byte, error) {\n\tfiles := map[string][]byte{}\n\tif opts.From == \"\" {\n\t\tbuiltinFiles, err := builtin.Files(name)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"%w (use --from <dir> for a package that is not built in)\", err)\n\t\t}\n\t\tfor rel, data := range builtinFiles {\n\t\t\tfiles[filepath.FromSlash(rel)] = data\n\t\t}\n\t\treturn files, nil\n\t}\n\troot := filepath.Join(opts.From, name)\n\tinfo, err := os.Stat(root)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"package %q: %w\", name, err)\n\t}\n\tif !info.IsDir() {\n\t\treturn nil, fmt.Errorf(\"package %q: %s is not a directory\", name, root)\n\t}\n\terr = filepath.Walk(root, func(p string, fi os.FileInfo, walkErr error) error {\n\t\tif walkErr != nil {\n\t\t\treturn walkErr\n\t\t}\n\t\tif !fi.Mode().IsRegular() {\n\t\t\treturn nil\n\t\t}\n\t\trel, err := filepath.Rel(root, p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tdata, err := os.ReadFile(p)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tfiles[rel] = data\n\t\treturn nil\n\t})\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn files, nil\n}","file":"internal/cli/generators.go","name":"loadPackageFiles"}}
func loadPackageFiles(name string, opts installOptions) (map[string][]byte, error) {
	files := map[string][]byte{}
	if opts.From == "" {
		builtinFiles, err := builtin.Files(name)
		if err != nil {
			return nil, fmt.Errorf("%w (use --from <dir> for a package that is not built in)", err)
		}
		for rel, data := range builtinFiles {
			files[filepath.FromSlash(rel)] = data
		}
		return files, nil
	}
	root := filepath.Join(opts.From, name)
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("package %q: %w", name, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("package %q: %s is not a directory", name, root)
	}
	err = filepath.Walk(root, func(p string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		files[rel] = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-3e5a05dc","gen":"sha256:601088eb1aa16c655a8a3610e01f220fcceaa8b83b311c76b4ac9a810345d516","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func applyInstall(out io.Writer, names []string, opts installOptions) error {\n\tentries, err := planInstall(names, opts)\n\tif err != nil {\n\t\treturn err\n\t}\n\tfor _, e := range entries {\n\t\tfmt.Fprintf(out, \"%s %s\\n\", e.Status, e.Path)\n\t}\n\tgoMod := filepath.Join(opts.Dir, \"go.mod\")\n\tgoModExists := false\n\tif _, err := os.Stat(goMod); err == nil {\n\t\tgoModExists = true\n\t}\n\tif goModExists {\n\t\tfmt.Fprintf(out, \"exists %s (left as it is)\\n\", goMod)\n\t} else {\n\t\tfmt.Fprintf(out, \"create %s\\n\", goMod)\n\t}\n\tif opts.DryRun {\n\t\tfmt.Fprintln(out, \"dry run: nothing written\")\n\t\treturn nil\n\t}\n\tvar differs []string\n\tfor _, e := range entries {\n\t\tif e.Status == \"differs\" {\n\t\t\tdiffers = append(differs, e.Path)\n\t\t}\n\t}\n\tif len(differs) > 0 && !opts.Force {\n\t\tshown := differs\n\t\tif len(shown) > 5 {\n\t\t\tshown = shown[:5]\n\t\t}\n\t\treturn fmt.Errorf(\"%d file(s) differ from the package, e.g. %s; rerun with --force to overwrite them\", len(differs), strings.Join(shown, \", \"))\n\t}\n\twritten, identical, overwritten := 0, 0, 0\n\tfor _, e := range entries {\n\t\tif e.Status == \"identical\" {\n\t\t\tidentical++\n\t\t\tcontinue\n\t\t}\n\t\tif e.Status == \"differs\" {\n\t\t\toverwritten++\n\t\t}\n\t\tif err := os.MkdirAll(filepath.Dir(e.Path), 0o755); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif err := os.WriteFile(e.Path, e.Data, 0o644); err != nil {\n\t\t\treturn err\n\t\t}\n\t\twritten++\n\t}\n\tif !goModExists {\n\t\tif err := os.MkdirAll(opts.Dir, 0o755); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tcontent := \"// Generator packages hold templates, not Go code. This module boundary keeps\\n// the enclosing module's ./... from compiling them.\\nmodule generators\\n\\ngo 1.21\\n\"\n\t\tif err := os.WriteFile(goMod, []byte(content), 0o644); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\tfmt.Fprintf(out, \"installed %d file(s) into %s (%d identical, %d overwritten)\\n\", written, opts.Dir, identical, overwritten)\n\treturn nil\n}","file":"internal/cli/generators.go","name":"applyInstall"}}
func applyInstall(out io.Writer, names []string, opts installOptions) error {
	entries, err := planInstall(names, opts)
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Fprintf(out, "%s %s\n", e.Status, e.Path)
	}
	goMod := filepath.Join(opts.Dir, "go.mod")
	goModExists := false
	if _, err := os.Stat(goMod); err == nil {
		goModExists = true
	}
	if goModExists {
		fmt.Fprintf(out, "exists %s (left as it is)\n", goMod)
	} else {
		fmt.Fprintf(out, "create %s\n", goMod)
	}
	if opts.DryRun {
		fmt.Fprintln(out, "dry run: nothing written")
		return nil
	}
	var differs []string
	for _, e := range entries {
		if e.Status == "differs" {
			differs = append(differs, e.Path)
		}
	}
	if len(differs) > 0 && !opts.Force {
		shown := differs
		if len(shown) > 5 {
			shown = shown[:5]
		}
		return fmt.Errorf("%d file(s) differ from the package, e.g. %s; rerun with --force to overwrite them", len(differs), strings.Join(shown, ", "))
	}
	written, identical, overwritten := 0, 0, 0
	for _, e := range entries {
		if e.Status == "identical" {
			identical++
			continue
		}
		if e.Status == "differs" {
			overwritten++
		}
		if err := os.MkdirAll(filepath.Dir(e.Path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(e.Path, e.Data, 0o644); err != nil {
			return err
		}
		written++
	}
	if !goModExists {
		if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
			return err
		}
		content := "// Generator packages hold templates, not Go code. This module boundary keeps\n// the enclosing module's ./... from compiling them.\nmodule generators\n\ngo 1.21\n"
		if err := os.WriteFile(goMod, []byte(content), 0o644); err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "installed %d file(s) into %s (%d identical, %d overwritten)\n", written, opts.Dir, identical, overwritten)
	return nil
}

// /sedum:setFunction
