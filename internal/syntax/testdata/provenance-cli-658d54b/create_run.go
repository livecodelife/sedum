package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"crypto/rand"}}
import "crypto/rand"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"encoding/hex"}}
import "encoding/hex"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"github.com/livecodelife/provenance/record"}}
import "github.com/livecodelife/provenance/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/create_run.go","path":"io"}}
import "io"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func newID(year int, suffix string, taken map[string]bool) (string, error) {\n\tfor i := 0; i < 100; i++ {\n\t\tb := make([]byte, 4)\n\t\tif _, err := rand.Read(b); err != nil {\n\t\t\treturn \"\", err\n\t\t}\n\t\tid := fmt.Sprintf(\"prov-%d-%s\", year, hex.EncodeToString(b))\n\t\tif suffix != \"\" {\n\t\t\tid += \"-\" + suffix\n\t\t}\n\t\tif !taken[id] {\n\t\t\treturn id, nil\n\t\t}\n\t}\n\treturn \"\", fmt.Errorf(\"could not generate a unique id after 100 attempts\")\n}","file":"cmd/provenance/create_run.go","name":"newID"}}
func newID(year int, suffix string, taken map[string]bool) (string, error) {
	for i := 0; i < 100; i++ {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		id := fmt.Sprintf("prov-%d-%s", year, hex.EncodeToString(b))
		if suffix != "" {
			id += "-" + suffix
		}
		if !taken[id] {
			return id, nil
		}
	}
	return "", fmt.Errorf("could not generate a unique id after 100 attempts")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runCreate(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\ttitle, _ := cmd.Flags().GetString(\"title\")\n\tintent, _ := cmd.Flags().GetString(\"intent\")\n\ttyp, _ := cmd.Flags().GetString(\"type\")\n\ttags, _ := cmd.Flags().GetStringArray(\"tag\")\n\tscope, _ := cmd.Flags().GetStringArray(\"scope\")\n\timplements, _ := cmd.Flags().GetString(\"implements\")\n\tlocked, _ := cmd.Flags().GetBool(\"locked\")\n\tnoEdit, _ := cmd.Flags().GetBool(\"no-edit\")\n\tsuffix, _ := cmd.Flags().GetString(\"id-suffix\")\n\n\tif typ != \"\" && !record.Type(typ).Valid(2) {\n\t\treturn fmt.Errorf(\"invalid type %q (brief|blueprint|imprint)\", typ)\n\t}\n\n\tcleanTags := []string{}\n\tfor _, t := range tags {\n\t\tfor _, p := range strings.Split(t, \",\") {\n\t\t\tp = strings.TrimSpace(p)\n\t\t\tif p != \"\" {\n\t\t\t\tcleanTags = append(cleanTags, p)\n\t\t\t}\n\t\t}\n\t}\n\n\ttaken := map[string]bool{}\n\tids, err := r.Store.List()\n\tif err != nil {\n\t\treturn err\n\t}\n\tfor _, existing := range ids {\n\t\ttaken[existing] = true\n\t\tloaded, _, err := r.Store.Load(existing)\n\t\tif err != nil {\n\t\t\tcontinue\n\t\t}\n\t\tfor _, a := range loaded.Aliases {\n\t\t\ttaken[a] = true\n\t\t}\n\t}\n\n\tid, err := newID(time.Now().Year(), suffix, taken)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tvar implementsID string\n\tvar implementsVersion int\n\tif implements != \"\" {\n\t\tparent, _, err := r.resolve(implements)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\timplementsID = parent.ID\n\t\timplementsVersion = parent.Version\n\t}\n\n\tif scope == nil {\n\t\tscope = []string{}\n\t}\n\n\trec := &record.Record{\n\t\tSchema:            2,\n\t\tID:                id,\n\t\tTitle:             title,\n\t\tType:              record.Type(typ),\n\t\tStatus:            record.StatusDraft,\n\t\tCreatedAt:         time.Now().Format(\"2006-01-02\"),\n\t\tAuthor:            r.author(),\n\t\tVersion:           1,\n\t\tIntent:            intent,\n\t\tLocked:            locked,\n\t\tImplements:        implementsID,\n\t\tImplementsVersion: implementsVersion,\n\t\tConstraints:       []string{},\n\t\tAffectedScope:     scope,\n\t\tForbiddenScope:    []string{},\n\t\tRelated:           []string{},\n\t\tAssociatedTraces:  []string{},\n\t\tMonitors:          []string{},\n\t\tTags:              cleanTags,\n\t\tAssociatedSpecs:   []record.AssociatedSpec{},\n\t}\n\n\tif err := r.Store.Save(rec, nil); err != nil {\n\t\treturn err\n\t}\n\n\tpath := filepath.Join(r.Root, r.recordPath(id))\n\n\tif !noEdit {\n\t\teditor := os.Getenv(\"EDITOR\")\n\t\tif editor == \"\" {\n\t\t\teditor = \"vi\"\n\t\t}\n\t\tc := exec.Command(editor, path)\n\t\tc.Stdin = os.Stdin\n\t\tc.Stdout = os.Stdout\n\t\tc.Stderr = os.Stderr\n\t\tif err := c.Run(); err != nil {\n\t\t\tfmt.Fprintf(os.Stderr, \"Warning: Could not open editor: %v\\n\", err)\n\t\t}\n\t}\n\n\tfmt.Fprintf(out, \"\\n✓ Created %s\\n\\n\", path)\n\treturn nil\n}","file":"cmd/provenance/create_run.go","name":"runCreate"}}
func runCreate(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}

	title, _ := cmd.Flags().GetString("title")
	intent, _ := cmd.Flags().GetString("intent")
	typ, _ := cmd.Flags().GetString("type")
	tags, _ := cmd.Flags().GetStringArray("tag")
	scope, _ := cmd.Flags().GetStringArray("scope")
	implements, _ := cmd.Flags().GetString("implements")
	locked, _ := cmd.Flags().GetBool("locked")
	noEdit, _ := cmd.Flags().GetBool("no-edit")
	suffix, _ := cmd.Flags().GetString("id-suffix")

	if typ != "" && !record.Type(typ).Valid(2) {
		return fmt.Errorf("invalid type %q (brief|blueprint|imprint)", typ)
	}

	cleanTags := []string{}
	for _, t := range tags {
		for _, p := range strings.Split(t, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				cleanTags = append(cleanTags, p)
			}
		}
	}

	taken := map[string]bool{}
	ids, err := r.Store.List()
	if err != nil {
		return err
	}
	for _, existing := range ids {
		taken[existing] = true
		loaded, _, err := r.Store.Load(existing)
		if err != nil {
			continue
		}
		for _, a := range loaded.Aliases {
			taken[a] = true
		}
	}

	id, err := newID(time.Now().Year(), suffix, taken)
	if err != nil {
		return err
	}

	var implementsID string
	var implementsVersion int
	if implements != "" {
		parent, _, err := r.resolve(implements)
		if err != nil {
			return err
		}
		implementsID = parent.ID
		implementsVersion = parent.Version
	}

	if scope == nil {
		scope = []string{}
	}

	rec := &record.Record{
		Schema:            2,
		ID:                id,
		Title:             title,
		Type:              record.Type(typ),
		Status:            record.StatusDraft,
		CreatedAt:         time.Now().Format("2006-01-02"),
		Author:            r.author(),
		Version:           1,
		Intent:            intent,
		Locked:            locked,
		Implements:        implementsID,
		ImplementsVersion: implementsVersion,
		Constraints:       []string{},
		AffectedScope:     scope,
		ForbiddenScope:    []string{},
		Related:           []string{},
		AssociatedTraces:  []string{},
		Monitors:          []string{},
		Tags:              cleanTags,
		AssociatedSpecs:   []record.AssociatedSpec{},
	}

	if err := r.Store.Save(rec, nil); err != nil {
		return err
	}

	path := filepath.Join(r.Root, r.recordPath(id))

	if !noEdit {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open editor: %v\n", err)
		}
	}

	fmt.Fprintf(out, "\n✓ Created %s\n\n", path)
	return nil
}

// /sedum:setFunction
