package main

import "github.com/spf13/cobra"

func newReviseCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"revise","short":"Edit a record's content in $EDITOR and save it as a new version","use":"revise [id]"}}
	cmd.Use = "revise [id]"
	cmd.Short = "Edit a record's content in $EDITOR and save it as a new version"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"revise"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runRevise(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:max_one {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"max_one","bin":"provenance","command":"revise"}}
	cmd.Args = cobra.MaximumNArgs(1)
	// /sedum:setArgs:max_one
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"revise","flag":"record","type":"string","usage":"Record to revise (or give it as the argument)"}}
	cmd.Flags().String("record", "", "Record to revise (or give it as the argument)")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"revise","flag":"reason","type":"string","usage":"Why this version exists (required for a new version)"}}
	cmd.Flags().String("reason", "", "Why this version exists (required for a new version)")
	// /sedum:addFlag:string
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-5f69977b","gen":"sha256:642b9ab38aa420dda21257a4b6af5f066c57a46f187200ab03c86ce74c73b219","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"revise","flag":"fix","type":"bool","usage":"Mark this version as a fix (reason_kind: fix)"}}
	cmd.Flags().Bool("fix", false, "Mark this version as a fix (reason_kind: fix)")
	// /sedum:addFlag:bool
	// sedum:anchor:subcommands
	return cmd
}
