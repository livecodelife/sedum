package main

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/hooks_install_run.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/hooks_install_run.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/hooks_install_run.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/hooks_install_run.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"cmd/provenance/hooks_install_run.go","path":"github.com/spf13/cobra"}}
import "github.com/spf13/cobra"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const preCommitHook = `#!/bin/sh\n# provenance pre-commit hook, installed by ` + \"`\" + `provenance hooks install` + \"`\" + `.\n# Lints each staged provenance record with the config nearest to it, and runs\n# a record's associated specs when it moves from open to implemented.\n\nif [ -n \"$PROVENANCE\" ]; then\n    PROV=\"$PROVENANCE\"\nelif [ -x \"./provenance\" ]; then\n    PROV=\"./provenance\"\nelse\n    PROV=\"provenance\"\nfi\n\nstaged=$(git diff --cached --name-only --diff-filter=ACM | grep -E \"prov-[0-9]{4}-[^/]*\\.ya?ml$\" | grep -v \"/\\.versions/\")\nif [ -z \"$staged\" ]; then\n    exit 0\nfi\n\n# Walk up from a file to the nearest .provenance.yml, else .linespec.yml.\nfind_config() {\n    dir=$(dirname \"$1\")\n    while true; do\n        for name in .provenance.yml .linespec.yml; do\n            if [ -f \"$dir/$name\" ]; then\n                echo \"$dir/$name\"\n                return 0\n            fi\n        done\n        if [ \"$dir\" = \".\" ] || [ \"$dir\" = \"/\" ]; then\n            return 1\n        fi\n        dir=$(dirname \"$dir\")\n    done\n}\n\nstatus_of() {\n    git show \"$1\" 2>/dev/null | grep \"^status:\" | head -1 | awk '{print $2}'\n}\n\nexit_code=0\nfor file in $staged; do\n    record=$(basename \"$file\" | sed -E 's/\\.ya?ml$//')\n    config=$(find_config \"$file\")\n    if [ -n \"$config\" ]; then\n        set -- -c \"$config\"\n    else\n        set --\n    fi\n    if ! \"$PROV\" \"$@\" lint --record \"$record\"; then\n        echo \"Commit blocked due to lint errors in $record\"\n        exit_code=1\n    fi\n    if [ \"$(status_of \"HEAD:$file\")\" = \"open\" ] && [ \"$(status_of \":$file\")\" = \"implemented\" ]; then\n        if ! \"$PROV\" \"$@\" run-specs --record \"$record\"; then\n            exit_code=1\n        fi\n    fi\ndone\n\nexit $exit_code\n`","file":"cmd/provenance/hooks_install_run.go","name":"preCommitHook"}}
const preCommitHook = `#!/bin/sh
# provenance pre-commit hook, installed by ` + "`" + `provenance hooks install` + "`" + `.
# Lints each staged provenance record with the config nearest to it, and runs
# a record's associated specs when it moves from open to implemented.

if [ -n "$PROVENANCE" ]; then
    PROV="$PROVENANCE"
elif [ -x "./provenance" ]; then
    PROV="./provenance"
else
    PROV="provenance"
fi

staged=$(git diff --cached --name-only --diff-filter=ACM | grep -E "prov-[0-9]{4}-[^/]*\.ya?ml$" | grep -v "/\.versions/")
if [ -z "$staged" ]; then
    exit 0
fi

# Walk up from a file to the nearest .provenance.yml, else .linespec.yml.
find_config() {
    dir=$(dirname "$1")
    while true; do
        for name in .provenance.yml .linespec.yml; do
            if [ -f "$dir/$name" ]; then
                echo "$dir/$name"
                return 0
            fi
        done
        if [ "$dir" = "." ] || [ "$dir" = "/" ]; then
            return 1
        fi
        dir=$(dirname "$dir")
    done
}

status_of() {
    git show "$1" 2>/dev/null | grep "^status:" | head -1 | awk '{print $2}'
}

exit_code=0
for file in $staged; do
    record=$(basename "$file" | sed -E 's/\.ya?ml$//')
    config=$(find_config "$file")
    if [ -n "$config" ]; then
        set -- -c "$config"
    else
        set --
    fi
    if ! "$PROV" "$@" lint --record "$record"; then
        echo "Commit blocked due to lint errors in $record"
        exit_code=1
    fi
    if [ "$(status_of "HEAD:$file")" = "open" ] && [ "$(status_of ":$file")" = "implemented" ]; then
        if ! "$PROV" "$@" run-specs --record "$record"; then
            exit_code=1
        fi
    fi
done

exit $exit_code
`

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const commitMsgHook = `#!/bin/sh\n# provenance commit-msg hook, installed by ` + \"`\" + `provenance hooks install` + \"`\" + `.\n# Checks the staged files against provenance scope for every config they fall\n# under. A config outside the root applies only when files under its\n# directory are staged.\n\nif [ -n \"$PROVENANCE\" ]; then\n    PROV=\"$PROVENANCE\"\nelif [ -x \"./provenance\" ]; then\n    PROV=\"./provenance\"\nelse\n    PROV=\"provenance\"\nfi\n\nCOMMIT_MSG_FILE=\"$1\"\n\ndirs=$(find . \\( -name .provenance.yml -o -name .linespec.yml \\) -not -path \"./.git/*\" 2>/dev/null | xargs -n1 dirname 2>/dev/null | sort -u)\n\nif [ -z \"$dirs\" ]; then\n    if ! \"$PROV\" check --staged --message-file \"$COMMIT_MSG_FILE\"; then\n        echo \"\"\n        echo \"Commit blocked due to provenance scope violations\"\n        exit 1\n    fi\n    exit 0\nfi\n\nexit_code=0\nfor dir in $dirs; do\n    if [ -f \"$dir/.provenance.yml\" ]; then\n        config=\"$dir/.provenance.yml\"\n    else\n        config=\"$dir/.linespec.yml\"\n    fi\n    rel=$(echo \"$dir\" | sed 's|^\\./||')\n    if [ \"$rel\" != \".\" ]; then\n        if [ -z \"$(git diff --cached --name-only | grep \"^$rel/\" | head -1)\" ]; then\n            continue\n        fi\n    fi\n    if ! \"$PROV\" -c \"$config\" check --staged --message-file \"$COMMIT_MSG_FILE\"; then\n        exit_code=1\n    fi\ndone\n\nif [ $exit_code -ne 0 ]; then\n    echo \"\"\n    echo \"Commit blocked due to provenance scope violations\"\nfi\n\nexit $exit_code\n`","file":"cmd/provenance/hooks_install_run.go","name":"commitMsgHook"}}
const commitMsgHook = `#!/bin/sh
# provenance commit-msg hook, installed by ` + "`" + `provenance hooks install` + "`" + `.
# Checks the staged files against provenance scope for every config they fall
# under. A config outside the root applies only when files under its
# directory are staged.

if [ -n "$PROVENANCE" ]; then
    PROV="$PROVENANCE"
elif [ -x "./provenance" ]; then
    PROV="./provenance"
else
    PROV="provenance"
fi

COMMIT_MSG_FILE="$1"

dirs=$(find . \( -name .provenance.yml -o -name .linespec.yml \) -not -path "./.git/*" 2>/dev/null | xargs -n1 dirname 2>/dev/null | sort -u)

if [ -z "$dirs" ]; then
    if ! "$PROV" check --staged --message-file "$COMMIT_MSG_FILE"; then
        echo ""
        echo "Commit blocked due to provenance scope violations"
        exit 1
    fi
    exit 0
fi

exit_code=0
for dir in $dirs; do
    if [ -f "$dir/.provenance.yml" ]; then
        config="$dir/.provenance.yml"
    else
        config="$dir/.linespec.yml"
    fi
    rel=$(echo "$dir" | sed 's|^\./||')
    if [ "$rel" != "." ]; then
        if [ -z "$(git diff --cached --name-only | grep "^$rel/" | head -1)" ]; then
            continue
        fi
    fi
    if ! "$PROV" -c "$config" check --staged --message-file "$COMMIT_MSG_FILE"; then
        exit_code=1
    fi
done

if [ $exit_code -ne 0 ]; then
    echo ""
    echo "Commit blocked due to provenance scope violations"
fi

exit $exit_code
`

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func installHook(dir, name, body string) (string, error) {\n\tpath := filepath.Join(dir, name)\n\tresult := \"installed\"\n\tif existing, err := os.ReadFile(path); err == nil {\n\t\tif string(existing) == body {\n\t\t\treturn \"unchanged\", nil\n\t\t}\n\t\tif _, err := os.Stat(path + \".bak\"); os.IsNotExist(err) {\n\t\t\tif err := os.Rename(path, path+\".bak\"); err != nil {\n\t\t\t\treturn \"\", err\n\t\t\t}\n\t\t\tresult = \"replaced (previous hook kept as \" + name + \".bak)\"\n\t\t} else {\n\t\t\tresult = \"replaced\"\n\t\t}\n\t}\n\tif err := os.MkdirAll(dir, 0o755); err != nil {\n\t\treturn \"\", err\n\t}\n\tif err := os.WriteFile(path, []byte(body), 0o755); err != nil {\n\t\treturn \"\", err\n\t}\n\tif err := os.Chmod(path, 0o755); err != nil {\n\t\treturn \"\", err\n\t}\n\treturn result, nil\n}","file":"cmd/provenance/hooks_install_run.go","name":"installHook"}}
func installHook(dir, name, body string) (string, error) {
	path := filepath.Join(dir, name)
	result := "installed"
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == body {
			return "unchanged", nil
		}
		if _, err := os.Stat(path + ".bak"); os.IsNotExist(err) {
			if err := os.Rename(path, path+".bak"); err != nil {
				return "", err
			}
			result = "replaced (previous hook kept as " + name + ".bak)"
		} else {
			result = "replaced"
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		return "", err
	}
	if err := os.Chmod(path, 0o755); err != nil {
		return "", err
	}
	return result, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d1d789de","gen":"sha256:8c961e2ea0d8b299cc389434247f854067085a0d5b858b6e6070771cba42d950","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func runHooksInstall(out io.Writer, cmd *cobra.Command, args []string) error {\n\tr, err := openRepo(cmd)\n\tif err != nil {\n\t\treturn err\n\t}\n\tdir, err := r.git(\"rev-parse\", \"--git-path\", \"hooks\")\n\tif err != nil {\n\t\treturn err\n\t}\n\tif !filepath.IsAbs(dir) {\n\t\tdir = filepath.Join(r.Root, dir)\n\t}\n\tpreResult, err := installHook(dir, \"pre-commit\", preCommitHook)\n\tif err != nil {\n\t\treturn err\n\t}\n\tmsgResult, err := installHook(dir, \"commit-msg\", commitMsgHook)\n\tif err != nil {\n\t\treturn err\n\t}\n\tfmt.Fprintf(out, \"\\n✓ Installed git hooks to %s\\n\\n\", dir)\n\tfmt.Fprintf(out, \"  pre-commit (%s):\\n    · Lints staged provenance records (nearest .provenance.yml, else .linespec.yml)\\n    · Runs a record's associated specs when it moves open -> implemented\\n\\n\", preResult)\n\tfmt.Fprintf(out, \"  commit-msg (%s):\\n    · Checks staged files against provenance scope (provenance check --staged)\\n    · Validates provenance tags, including [id@version]\\n\\n\", msgResult)\n\tfmt.Fprint(out, \"  The hooks run $PROVENANCE, else ./provenance, else provenance on PATH.\\n\\n\")\n\treturn nil\n}","file":"cmd/provenance/hooks_install_run.go","name":"runHooksInstall"}}
func runHooksInstall(out io.Writer, cmd *cobra.Command, args []string) error {
	r, err := openRepo(cmd)
	if err != nil {
		return err
	}
	dir, err := r.git("rev-parse", "--git-path", "hooks")
	if err != nil {
		return err
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(r.Root, dir)
	}
	preResult, err := installHook(dir, "pre-commit", preCommitHook)
	if err != nil {
		return err
	}
	msgResult, err := installHook(dir, "commit-msg", commitMsgHook)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "\n✓ Installed git hooks to %s\n\n", dir)
	fmt.Fprintf(out, "  pre-commit (%s):\n    · Lints staged provenance records (nearest .provenance.yml, else .linespec.yml)\n    · Runs a record's associated specs when it moves open -> implemented\n\n", preResult)
	fmt.Fprintf(out, "  commit-msg (%s):\n    · Checks staged files against provenance scope (provenance check --staged)\n    · Validates provenance tags, including [id@version]\n\n", msgResult)
	fmt.Fprint(out, "  The hooks run $PROVENANCE, else ./provenance, else provenance on PATH.\n\n")
	return nil
}

// /sedum:setFunction
