package main

import "github.com/spf13/cobra"

func newOpenCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"open","short":"Transition a record from draft to open (enables enforcement)","use":"open"}}
	cmd.Use = "open"
	cmd.Short = "Transition a record from draft to open (enables enforcement)"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"open"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runOpen(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"open"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-6db88b75","gen":"sha256:b330aeacbc6612174595999e1495b3d3a8b2ec816017256a5b834b49a30d8c20","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"open","flag":"record","type":"string","usage":"Record to transition from draft to open"}}
	cmd.Flags().String("record", "", "Record to transition from draft to open")
	// /sedum:addFlag:string
	// sedum:anchor:subcommands
	return cmd
}
