package main

import "github.com/spf13/cobra"

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"lint","short":"Validate provenance records","use":"lint"}}
	cmd.Use = "lint"
	cmd.Short = "Validate provenance records"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"lint"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runLint(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"lint"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"record","type":"string","usage":"Lint a single record"}}
	cmd.Flags().String("record", "", "Lint a single record")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"enforcement","type":"string","usage":"Override enforcement (none|warn|strict)"}}
	cmd.Flags().String("enforcement", "", "Override enforcement (none|warn|strict)")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"format","type":"string","usage":"Output format (human|json)"}}
	cmd.Flags().String("format", "", "Output format (human|json)")
	// /sedum:addFlag:string
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"warn","type":"bool","usage":"Show only warnings in output"}}
	cmd.Flags().Bool("warn", false, "Show only warnings in output")
	// /sedum:addFlag:bool
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"info","type":"bool","usage":"Show only informational hints in output"}}
	cmd.Flags().Bool("info", false, "Show only informational hints in output")
	// /sedum:addFlag:bool
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-a03e6bc3","gen":"sha256:aec5a8b3297a63ed73463101f5712b37b803293254ec536f0c5009748de88f09","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"lint","flag":"all","type":"bool","usage":"Show all output (errors, warnings, and hints)"}}
	cmd.Flags().Bool("all", false, "Show all output (errors, warnings, and hints)")
	// /sedum:addFlag:bool
	// sedum:anchor:subcommands
	return cmd
}
