package main

import "github.com/spf13/cobra"

func newDeprecateCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"deprecate","short":"Mark a record as deprecated (terminal: it takes no new versions)","use":"deprecate"}}
	cmd.Use = "deprecate"
	cmd.Short = "Mark a record as deprecated (terminal: it takes no new versions)"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"deprecate"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runDeprecate(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"deprecate"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"deprecate","flag":"record","type":"string","usage":"Record to deprecate"}}
	cmd.Flags().String("record", "", "Record to deprecate")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-35439261","gen":"sha256:7a7f9652ab699a1d6305ff28ecf9f26e5a23c00aa71b77596c96dedb6c805fe1","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"deprecate","flag":"reason","type":"string","usage":"Why the record is deprecated (recorded in the commit message)"}}
	cmd.Flags().String("reason", "", "Why the record is deprecated (recorded in the commit message)")
	// /sedum:addFlag:string
	// sedum:anchor:subcommands
	return cmd
}
