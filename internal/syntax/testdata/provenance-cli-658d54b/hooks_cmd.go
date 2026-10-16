package main

import "github.com/spf13/cobra"

func newHooksCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"hooks","short":"Manage the provenance git hooks","use":"hooks"}}
	cmd.Use = "hooks"
	cmd.Short = "Manage the provenance git hooks"
	// /sedum:describeCommand
	// sedum:anchor:flags
	// sedum:anchor:subcommands
	// sedum:registerSubcommand {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:8a9e0a8ee811e7c629264e26f09a626602890b11fd82710cf10ae73e3da1cf7a","kwargs":{"bin":"provenance","command":"hooks-install","parent":"hooks"}}
	cmd.AddCommand(newHooksInstallCmd())
	// /sedum:registerSubcommand
	return cmd
}
