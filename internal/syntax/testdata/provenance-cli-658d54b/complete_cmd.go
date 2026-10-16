package main

import "github.com/spf13/cobra"

func newCompleteCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"complete","short":"Mark an open record as implemented (seals it at HEAD)","use":"complete"}}
	cmd.Use = "complete"
	cmd.Short = "Mark an open record as implemented (seals it at HEAD)"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"complete"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runComplete(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"complete"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"complete","flag":"record","type":"string","usage":"Record to mark as implemented"}}
	cmd.Flags().String("record", "", "Record to mark as implemented")
	// /sedum:addFlag:string
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-41fb548c","gen":"sha256:7dee87387bdca578508d9dfc08b17eb82cf1299d88e96557bc86ce81c652e1bf","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"complete","flag":"force","type":"bool","usage":"Skip the spec existence check and spec runs"}}
	cmd.Flags().Bool("force", false, "Skip the spec existence check and spec runs")
	// /sedum:addFlag:bool
	// sedum:anchor:subcommands
	return cmd
}
