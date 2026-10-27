package scaffold

// Hook is one git hook sedum installs: the file name and the script text.
type Hook struct {
	Name string
	Body string
}

// Hooks returns the hooks sedum hooks install writes, commit-msg then pre-push.
func Hooks() []Hook {
	return []Hook{
		{Name: "commit-msg", Body: commitMsgHook},
		{Name: "pre-push", Body: prePushHook},
	}
}

const commitMsgHook = `#!/bin/sh
# sedum commit-msg hook, installed by ` + "`" + `sedum hooks install` + "`" + `.
# Checks the staged change with sedum check, once for each provenance config
# root that has staged files under it.

# $SEDUM, else sedum on PATH. Never a file found by name in the repository.
if [ -n "$SEDUM" ]; then
    SED="$SEDUM"
elif command -v sedum >/dev/null 2>&1; then
    SED="sedum"
else
    echo "sedum hook: no sedum binary. Set SEDUM or put sedum on PATH." >&2
    exit 1
fi

MSG="$1"

# A .linespec.yml is a root only when it has a provenance block.
dirs=$(find . \( -name .provenance.yml -o -name .linespec.yml \) -not -path "./.git/*" 2>/dev/null | while read -r f; do
    case "$f" in
        */.provenance.yml) echo "$f" ;;
        *) grep -q '^provenance:' "$f" && echo "$f" ;;
    esac
done | xargs -n1 dirname 2>/dev/null | sort -u)

failed=0
if [ -z "$dirs" ]; then
    "$SED" check --staged --message-file "$MSG" || failed=1
else
    for dir in $dirs; do
        rel=$(echo "$dir" | sed 's|^\./||')
        if [ -f "$dir/.provenance.yml" ]; then
            config="$dir/.provenance.yml"
        else
            config="$dir/.linespec.yml"
        fi
        config=$(echo "$config" | sed 's|^\./||')
        if [ "$rel" != "." ]; then
            if [ -z "$(git diff --cached --name-only | grep "^$rel/" | head -n 1)" ]; then
                continue
            fi
        fi
        "$SED" check --staged --message-file "$MSG" --config "$config" || failed=1
    done
fi

if [ "$failed" -ne 0 ]; then
    echo ""
    echo "Commit blocked due to sedum check failures"
    exit 1
fi
exit 0
`

const prePushHook = `#!/bin/sh
# sedum pre-push hook, installed by ` + "`" + `sedum hooks install` + "`" + `.
# Runs provenance's own pre-push hook, then checks the commits being pushed
# with sedum check, one range for each pushed ref.

# $SEDUM, else sedum on PATH. Never a file found by name in the repository.
if [ -n "$SEDUM" ]; then
    SED="$SEDUM"
elif command -v sedum >/dev/null 2>&1; then
    SED="sedum"
else
    echo "sedum hook: no sedum binary. Set SEDUM or put sedum on PATH." >&2
    exit 1
fi

# $PROVENANCE, else provenance on PATH.
if [ -n "$PROVENANCE" ]; then
    PROV="$PROVENANCE"
elif command -v provenance >/dev/null 2>&1; then
    PROV="provenance"
else
    echo "sedum hook: no provenance binary. Set PROVENANCE or put provenance on PATH." >&2
    exit 1
fi

saved=$(mktemp)
trap 'rm -f "$saved"' EXIT
cat > "$saved"

failed=0
"$PROV" hooks pre-push "$@" < "$saved" || failed=1

checked=" "
while read -r lref lsha rref rsha rest; do
    [ -n "$rsha" ] || continue
    case "$lsha" in
        *[!0]*) ;;
        *) continue ;;
    esac
    case "$checked" in
        *" $lsha "*) continue ;;
    esac
    checked="$checked$lsha "

    range=""
    if case "$rsha" in *[!0]*) true ;; *) false ;; esac && git cat-file -e "$rsha^{commit}" 2>/dev/null; then
        range="$rsha..$lsha"
    else
        oldest=$(git rev-list --reverse "$lsha" --not --remotes 2>/dev/null | head -n 1)
        [ -n "$oldest" ] || continue
        if parent=$(git rev-parse --verify -q "$oldest^1" 2>/dev/null); then
            range="$parent..$lsha"
        else
            "$PROV" check --commit "$oldest" < /dev/null || failed=1
            continue
        fi
    fi
    "$SED" check --range "$range" < /dev/null || failed=1
done < "$saved"

if [ "$failed" -ne 0 ]; then
    echo ""
    echo "Push blocked due to failed checks"
    exit 1
fi
exit 0
`
