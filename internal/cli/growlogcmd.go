package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"strconv"}}
import "strconv"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"github.com/livecodelife/sedum/internal/growlog"}}
import "github.com/livecodelife/sedum/internal/growlog"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlogcmd.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// parseSince reads --since: a Go duration (\"36h\"), a number of days (\"7d\"), a\n// date (\"2026-10-01\") or an RFC3339 time. The empty string is the zero time.\nfunc parseSince(s string, now time.Time) (time.Time, error) {\n\tif s == \"\" {\n\t\treturn time.Time{}, nil\n\t}\n\tif days, ok := strings.CutSuffix(s, \"d\"); ok {\n\t\tif n, err := strconv.Atoi(days); err == nil {\n\t\t\treturn now.AddDate(0, 0, -n), nil\n\t\t}\n\t}\n\tif d, err := time.ParseDuration(s); err == nil {\n\t\treturn now.Add(-d), nil\n\t}\n\tfor _, layout := range []string{\"2006-01-02\", time.RFC3339} {\n\t\tif t, err := time.Parse(layout, s); err == nil {\n\t\t\treturn t, nil\n\t\t}\n\t}\n\treturn time.Time{}, fmt.Errorf(\"--since %q is not a duration (36h), a number of days (7d), a date (2026-10-01) or an RFC3339 time\", s)\n}","file":"internal/cli/growlogcmd.go","name":"parseSince"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// runGrowLog prints the summary of the grow log at path.\nfunc runGrowLog(out io.Writer, path string, asJSON bool, since, record, model string, now time.Time) error {\n\tfrom, err := parseSince(since, now)\n\tif err != nil {\n\t\treturn err\n\t}\n\tentries, malformed, err := growlog.Read(path)\n\tif err != nil {\n\t\treturn err\n\t}\n\ts := growlog.Summarize(entries, growlog.Filter{Since: from, Record: record, Model: model})\n\ts.Malformed = malformed\n\tif asJSON {\n\t\tb, err := s.JSON()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t_, err = out.Write(b)\n\t\treturn err\n\t}\n\t_, err = io.WriteString(out, s.Text())\n\treturn err\n}","file":"internal/cli/growlogcmd.go","name":"runGrowLog"}}
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
// sedum:setFunction {"tier":"owned","record":"prov-2026-87743cfa","gen":"sha256:4cc82a967421f4715a1da93b5093f4959c47a025a161533ab78ae7307539a5a9","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// newGrowLogCommand builds `sedum grow-log`.\nfunc newGrowLogCommand() *cobra.Command {\n\tvar asJSON bool\n\tvar since, record, model string\n\tcmd := &cobra.Command{\n\t\tUse:   \"grow-log [path]\",\n\t\tShort: \"Summarize the persistent grow log\",\n\t\tLong:  \"Summarize the persistent grow log (default \" + defaultGrowLogPath + \"): per model, how many grows, failures, model calls, retries (answers the model had refused), repairs (verification send-backs), tokens and the share of invocations that carried model-written code.\",\n\t\tArgs:  cobra.MaximumNArgs(1),\n\t\tRunE: func(cmd *cobra.Command, args []string) error {\n\t\t\tpath := defaultGrowLogPath\n\t\t\tif len(args) == 1 {\n\t\t\t\tpath = args[0]\n\t\t\t}\n\t\t\treturn runGrowLog(cmd.OutOrStdout(), path, asJSON, since, record, model, time.Now())\n\t\t},\n\t}\n\tcmd.Flags().BoolVar(&asJSON, \"json\", false, \"print the summary as JSON\")\n\tcmd.Flags().StringVar(&since, \"since\", \"\", \"only entries since a duration (36h), days (7d), a date (2026-10-01) or an RFC3339 time\")\n\tcmd.Flags().StringVar(&record, \"record\", \"\", \"only entries for this record id\")\n\tcmd.Flags().StringVar(&model, \"model\", \"\", \"only entries for this model\")\n\treturn cmd\n}","file":"internal/cli/growlogcmd.go","name":"newGrowLogCommand"}}
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
