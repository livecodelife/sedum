package main

import "github.com/spf13/cobra"

func newHooksInstallCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"hooks-install","short":"Install the pre-commit and commit-msg hooks (existing hooks are kept as .bak)","use":"install"}}
	cmd.Use = "install"
	cmd.Short = "Install the pre-commit and commit-msg hooks (existing hooks are kept as .bak)"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"hooks-install"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runHooksInstall(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"hooks-install"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:anchor:subcommands
	return cmd
}
