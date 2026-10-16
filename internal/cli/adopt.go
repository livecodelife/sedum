package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt.go","path":"github.com/livecodelife/sedum/internal/adopt"}}
import "github.com/livecodelife/sedum/internal/adopt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/adopt.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// AdoptConfig holds the flags of sedum adopt.\ntype AdoptConfig struct {\n\tGenerators, Records, Output, Recording    string\n\tTests, AllowCommandAdopters, DryRun, JSON bool\n}","file":"internal/cli/adopt.go","name":"AdoptConfig"}}
// AdoptConfig holds the flags of sedum adopt.
type AdoptConfig struct {
	Generators, Records, Output, Recording    string
	Tests, AllowCommandAdopters, DryRun, JSON bool
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func newAdoptCommand() *cobra.Command {\n\tvar cfg AdoptConfig\n\tcmd := &cobra.Command{\n\t\tUse:   \"adopt <path>...\",\n\t\tShort: \"Bring existing files under Sedum without a model\",\n\t\tLong:  \"Bring existing files under Sedum without a model. For each file, sedum splits it with the adopter its generator package declares, builds a recording, replays it, and keeps the file only if stripping the markers gives back the original. Files it cannot adopt are skipped and named with the reason. --dry-run writes only the recording; no model is called.\",\n\t\tArgs:  cobra.MinimumNArgs(1),\n\t\tRunE: func(cmd *cobra.Command, args []string) error {\n\t\t\treturn runAdopt(cmd.OutOrStdout(), cfg, args)\n\t\t},\n\t}\n\tcmd.Flags().StringVar(&cfg.Generators, \"generators\", \"\", \"directory of generator packages\")\n\tcmd.Flags().StringVar(&cfg.Records, \"records\", \"\", \"directory of records\")\n\tcmd.Flags().StringVar(&cfg.Output, \"output\", defaultOutput, \"the project root the paths are relative to\")\n\tcmd.Flags().StringVar(&cfg.Recording, \"recording\", \"\", \"where to write the recording (default .sedum/adopt/adopt.json under --output)\")\n\tcmd.Flags().BoolVar(&cfg.Tests, \"tests\", false, \"also adopt test files\")\n\tcmd.Flags().BoolVar(&cfg.AllowCommandAdopters, \"allow-command-adopters\", false, \"allow adopters that run an external command\")\n\tcmd.Flags().BoolVar(&cfg.DryRun, \"dry-run\", false, \"write only the recording; change no files\")\n\tcmd.Flags().BoolVar(&cfg.JSON, \"json\", false, \"print the report as JSON\")\n\tmustMarkRequired(cmd, \"generators\", \"records\")\n\treturn cmd\n}","file":"internal/cli/adopt.go","name":"newAdoptCommand"}}
func newAdoptCommand() *cobra.Command {
	var cfg AdoptConfig
	cmd := &cobra.Command{
		Use:   "adopt <path>...",
		Short: "Bring existing files under Sedum without a model",
		Long:  "Bring existing files under Sedum without a model. For each file, sedum splits it with the adopter its generator package declares, builds a recording, replays it, and keeps the file only if stripping the markers gives back the original. Files it cannot adopt are skipped and named with the reason. --dry-run writes only the recording; no model is called.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdopt(cmd.OutOrStdout(), cfg, args)
		},
	}
	cmd.Flags().StringVar(&cfg.Generators, "generators", "", "directory of generator packages")
	cmd.Flags().StringVar(&cfg.Records, "records", "", "directory of records")
	cmd.Flags().StringVar(&cfg.Output, "output", defaultOutput, "the project root the paths are relative to")
	cmd.Flags().StringVar(&cfg.Recording, "recording", "", "where to write the recording (default .sedum/adopt/adopt.json under --output)")
	cmd.Flags().BoolVar(&cfg.Tests, "tests", false, "also adopt test files")
	cmd.Flags().BoolVar(&cfg.AllowCommandAdopters, "allow-command-adopters", false, "allow adopters that run an external command")
	cmd.Flags().BoolVar(&cfg.DryRun, "dry-run", false, "write only the recording; change no files")
	cmd.Flags().BoolVar(&cfg.JSON, "json", false, "print the report as JSON")
	mustMarkRequired(cmd, "generators", "records")
	return cmd
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runAdopt(out io.Writer, cfg AdoptConfig, paths []string) error {\n\trep, err := adopt.Run(adopt.Config{\n\t\tRoot:          cfg.Output,\n\t\tGenerators:    cfg.Generators,\n\t\tRecords:       cfg.Records,\n\t\tPaths:         paths,\n\t\tTests:         cfg.Tests,\n\t\tAllowCommand:  cfg.AllowCommandAdopters,\n\t\tDryRun:        cfg.DryRun,\n\t\tRecordingPath: cfg.Recording,\n\t})\n\tif err != nil {\n\t\treturn err\n\t}\n\tif cfg.JSON {\n\t\tb, err := rep.JSON()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif _, err := out.Write(b); err != nil {\n\t\t\treturn err\n\t\t}\n\t} else {\n\t\tif _, err := io.WriteString(out, rep.Text()); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\tif rep.Totals.Rejected > 0 {\n\t\treturn fmt.Errorf(\"adopt: %d file(s) rejected\", rep.Totals.Rejected)\n\t}\n\treturn nil\n}","file":"internal/cli/adopt.go","name":"runAdopt"}}
func runAdopt(out io.Writer, cfg AdoptConfig, paths []string) error {
	rep, err := adopt.Run(adopt.Config{
		Root:          cfg.Output,
		Generators:    cfg.Generators,
		Records:       cfg.Records,
		Paths:         paths,
		Tests:         cfg.Tests,
		AllowCommand:  cfg.AllowCommandAdopters,
		DryRun:        cfg.DryRun,
		RecordingPath: cfg.Recording,
	})
	if err != nil {
		return err
	}
	if cfg.JSON {
		b, err := rep.JSON()
		if err != nil {
			return err
		}
		if _, err := out.Write(b); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(out, rep.Text()); err != nil {
			return err
		}
	}
	if rep.Totals.Rejected > 0 {
		return fmt.Errorf("adopt: %d file(s) rejected", rep.Totals.Rejected)
	}
	return nil
}

// /sedum:setFunction
