package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"fmt"},"kwargs":"sha256:4ff215b38a19c7936af1a3b50c08453e12c7309088a238da0d54c2acea086a14","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"io"},"kwargs":"sha256:61ce0bac05e8a0fa6d8e1e236312c7b3ccc13abe7c2ced36ea0c74dd4ab77576","body":"sha256:ef8e8abf4ef5ee40a19dc85b49bf09c5df1fe32abfe834e0ab4c2175baf80ddf"}
import "io"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"strconv"},"kwargs":"sha256:4b3c75ec5c7c126436247bf513c567adcef4410fc5cc48cffadc66d31f283e96","body":"sha256:f33f981c64a42d78a5daa7ff3fca15b10045080cccd282fef544b59eaf7f6dc8"}
import "strconv"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"strings"},"kwargs":"sha256:e0f8d3ee4751dbf6be26d9c435cd36f5408d2ce215c214bac29b1873f50aa868","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"time"},"kwargs":"sha256:4c6da2cf5f7569e6db002ecdd7d7e4a5c0f9edb7b754be90e39e3243de63a4d1","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"github.com/livecodelife/sedum/internal/growlog"},"kwargs":"sha256:56d86235b182d280f659c34907418cc7b70076962bbc6f8bdc24efd9e4f9b06a","body":"sha256:b7968f24df7d76473f917ee9527158f424f930250a91e5407a0140aa4c56c91f"}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlogcmd.go","path":"github.com/spf13/cobra"},"kwargs":"sha256:f05a6d32478032ec81469fca3fb190e15b0aee2f8bf6311c4ef531fcd031b94f","body":"sha256:f361c02d0f8e1999cf3fb715393e16310eeec9debd5aaa13936af7617e2b46aa"}
import "github.com/spf13/cobra"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd.go","name":"parseSince"},"kwargs":"sha256:2a7dea934433354c2e9dd8d4385b920228fe7d8423f2a1fee3c656d0cfbaff7a","body":"sha256:c4c373115f8069f572c6ca90095a8e143d78f2f5f49910cb9a023e686f30d168"}
// parseSince reads --since: a Go duration ("36h"), a number of days ("7d"), a
// date ("2026-10-01") or an RFC3339 time. The empty string is the zero time.
func parseSince(s string, now time.Time) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if days, ok := strings.CutSuffix(s, "d"); ok {
		if n, err := strconv.Atoi(days); err == nil {
			return now.AddDate(0, 0, -n), nil
		}
	}
	if d, err := time.ParseDuration(s); err == nil {
		return now.Add(-d), nil
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("--since %q is not a duration (36h), a number of days (7d), a date (2026-10-01) or an RFC3339 time", s)
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd.go","name":"runGrowLog"},"kwargs":"sha256:8bd9e6fe10d4dd58914c2f7be1eb7560985fdc16095582fbc4ac0e3f01bcf59f","body":"sha256:8910dec3a354286f1d0ac0f2f319a3b552e45f6dbc8d3fb2f97e0ab849b9ded4"}
// runGrowLog prints the summary of the grow log at path.
func runGrowLog(out io.Writer, path string, asJSON bool, since, record, model string, now time.Time) error {
	from, err := parseSince(since, now)
	if err != nil {
		return err
	}
	entries, malformed, err := growlog.Read(path)
	if err != nil {
		return err
	}
	s := growlog.Summarize(entries, growlog.Filter{Since: from, Record: record, Model: model})
	s.Malformed = malformed
	if asJSON {
		b, err := s.JSON()
		if err != nil {
			return err
		}
		_, err = out.Write(b)
		return err
	}
	_, err = io.WriteString(out, s.Text())
	return err
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlogcmd.go","name":"newGrowLogCommand"},"kwargs":"sha256:bdee7094df15681624ccff8d2e00271015eb0b252aa2010530765323c872fffa","body":"sha256:05ed55699b393e2c60a10d19f333f8de36fda7da36fddf1f7e39ef1c813257f4"}
// newGrowLogCommand builds `sedum grow-log`.
func newGrowLogCommand() *cobra.Command {
	var asJSON bool
	var since, record, model string
	cmd := &cobra.Command{
		Use:   "grow-log [path]",
		Short: "Summarize the persistent grow log",
		Long:  "Summarize the persistent grow log (default " + defaultGrowLogPath + "): per model, how many grows, failures, model calls, retries (answers the model had refused), repairs (verification send-backs), tokens and the share of invocations that carried model-written code.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := defaultGrowLogPath
			if len(args) == 1 {
				path = args[0]
			}
			return runGrowLog(cmd.OutOrStdout(), path, asJSON, since, record, model, time.Now())
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print the summary as JSON")
	cmd.Flags().StringVar(&since, "since", "", "only entries since a duration (36h), days (7d), a date (2026-10-01) or an RFC3339 time")
	cmd.Flags().StringVar(&record, "record", "", "only entries for this record id")
	cmd.Flags().StringVar(&model, "model", "", "only entries for this model")
	return cmd
}

// /sedum:setFunction
