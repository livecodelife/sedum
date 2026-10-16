package main

import "github.com/spf13/cobra"

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{}
	// sedum:anchor:command
	// sedum:describeCommand {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:cb672494161bfe1ff7809e5af7923edd680b02925293c87ba9726637c12e86af","kwargs":{"bin":"provenance","command":"create","short":"Create a new provenance record (defaults to draft status)","use":"create"}}
	cmd.Use = "create"
	cmd.Short = "Create a new provenance record (defaults to draft status)"
	// /sedum:describeCommand
	// sedum:setRun {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:5e369d6d2a4af49acd0fc42910733893fbbd57e90acd09046714692a30b3819e","kwargs":{"bin":"provenance","command":"create"}}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runCreate(cmd.OutOrStdout(), cmd, args)
	}
	// /sedum:setRun
	// sedum:setArgs:none {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:ca644378490be5811d57a9072cb3caf3c807f5d2c6deee0170e1bc58e0c2d3f9","kwargs":{"arity":"none","bin":"provenance","command":"create"}}
	cmd.Args = cobra.NoArgs
	// /sedum:setArgs:none
	// sedum:anchor:flags
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"title","type":"string","usage":"Pre-populate the title field"}}
	cmd.Flags().String("title", "", "Pre-populate the title field")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"intent","type":"string","usage":"Pre-populate the intent field"}}
	cmd.Flags().String("intent", "", "Pre-populate the intent field")
	// /sedum:addFlag:string
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"type","type":"string","usage":"Set the tier type of the record (brief|blueprint|imprint)"}}
	cmd.Flags().String("type", "", "Set the tier type of the record (brief|blueprint|imprint)")
	// /sedum:addFlag:string
	// sedum:addFlag:strings {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"tag","type":"strings","usage":"Pre-populate tags (comma-separated; repeatable)"}}
	cmd.Flags().StringArray("tag", nil, "Pre-populate tags (comma-separated; repeatable)")
	// /sedum:addFlag:strings
	// sedum:addFlag:strings {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"scope","type":"strings","usage":"Pre-populate affected_scope (repeatable)"}}
	cmd.Flags().StringArray("scope", nil, "Pre-populate affected_scope (repeatable)")
	// /sedum:addFlag:strings
	// sedum:addFlag:string {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"implements","type":"string","usage":"Parent record id; pins implements_version to its current version"}}
	cmd.Flags().String("implements", "", "Parent record id; pins implements_version to its current version")
	// /sedum:addFlag:string
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"locked","type":"bool","usage":"Create the record as a locked layer"}}
	cmd.Flags().Bool("locked", false, "Create the record as a locked layer")
	// /sedum:addFlag:bool
	// sedum:addFlag:bool {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:d873fc5f71035cf617eb56ca6ca05ce36f39842f4f774d2fb6b6816cf95ccc9b","kwargs":{"bin":"provenance","command":"create","flag":"no-edit","type":"bool","usage":"Write without opening editor"}}
	cmd.Flags().Bool("no-edit", false, "Write without opening editor")
	// /sedum:addFlag:bool
	// sedum:addShortFlag:string {"tier":"owned","record":"prov-2026-44a07768","gen":"sha256:ba3e6b98ee7f4297af267684cf4d5385650dd9703772a650529ecba804c273f6","pkg":"sha256:317b97244b0fe957a6e0b80cdfa819a8afaae50b42b10385c9cd6729cb74fd3a","kwargs":{"bin":"provenance","command":"create","flag":"id-suffix","shorthand":"i","type":"string","usage":"Append service suffix to ID (e.g., user-service)"}}
	cmd.Flags().StringP("id-suffix", "i", "", "Append service suffix to ID (e.g., user-service)")
	// /sedum:addShortFlag:string
	// sedum:anchor:subcommands
	return cmd
}
