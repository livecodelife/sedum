package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt.go","path":"fmt"},"kwargs":"sha256:fe9666ccf95a3f3d342a47f984c024a73724d7455dee49b2a26efa791e0492a2","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt.go","path":"io"},"kwargs":"sha256:d7dbd756ae5ee3bf6bcd4e19ba284f3417c443d1561f342d83dccb630b916b63","body":"sha256:ef8e8abf4ef5ee40a19dc85b49bf09c5df1fe32abfe834e0ab4c2175baf80ddf"}
import "io"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt.go","path":"github.com/livecodelife/sedum/internal/adopt"},"kwargs":"sha256:f14d10eda4b5b8a03f3561ecd7d0347301155d0eb4f4ebb24e5e0002988b9604","body":"sha256:14c711162319c41280adf8c125370979e81aede4a5448dd7800444f91371acd1"}
import "github.com/livecodelife/sedum/internal/adopt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/adopt.go","path":"github.com/spf13/cobra"},"kwargs":"sha256:1b63536416b07f6f18fd34caab41a5bde9afd4dc02b7ac7b1d2b2fb447126ebe","body":"sha256:f361c02d0f8e1999cf3fb715393e16310eeec9debd5aaa13936af7617e2b46aa"}
import "github.com/spf13/cobra"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/cli/adopt.go","name":"AdoptConfig"},"kwargs":"sha256:8478822c323ea020ba674f6c971a340d7d3bed2de977177bb40da2e9fc6efbb4","body":"sha256:88731b9e9350d8db0dab87f392be7c1d04f6777f1f01db3ad2f7fd637b8c95ec"}
// AdoptConfig holds the flags of sedum adopt.
type AdoptConfig struct {
	Generators, Records, Output, Recording    string
	Tests, AllowCommandAdopters, DryRun, JSON bool
}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt.go","name":"newAdoptCommand"},"kwargs":"sha256:3727a8490e47a3aff6e3702853e0a4393b3cd0bc55548b458598ebff98074383","body":"sha256:3f75c3a613214bc3b13b4642cd76c27c5002785da8f107b25ae527cfffac6408"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a77c2e9f","gen":"sha256:52f45c66197b92b8180e46912fc1ad38be50faebe69122f0b9cb64741d610d29","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/adopt.go","name":"runAdopt"},"kwargs":"sha256:f26edd4b899a6f939e6505261c748157d1351557ddc820f83c73095951489a34","body":"sha256:8201d6dd7a55ef7a689a3f3cd2d016192dff6747ddafd6cb02ade5f126f25890"}
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
