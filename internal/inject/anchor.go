package inject

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/calebcowen/sedum/internal/genpkg"
)

// Anchors: where in a file an action's output belongs.
//
// The vocabulary is small, closed, and evaluated entirely at the text level.
// There is no parser, no AST, and no syntax awareness of any kind - an anchor is
// found by a marker comment a file template planted or by a regex the action
// declared, and by nothing else. That is what lets Sedum's core know nothing
// about the target language.
//
// A missing anchor is a hard error. It means the file is not shaped the way the
// action assumed, which is a disagreement between configuration and reality that
// the author has to resolve. Auto-creating the anchor would paper over exactly
// the mistake worth surfacing.

// anchorDecl is the text a file template plants to create an anchor point.
//
// genpkg writes and recognizes this same shape when it checks that an action's
// anchor corresponds to a marker some file template plants. The shape is
// repeated here rather than shared because locating the line is this package's
// business and recognizing the name is that one's; markerShapesAgree in the
// tests asserts the two have not drifted apart.
func anchorDecl(commentPrefix, name string) string {
	return commentPrefix + " sedum:anchor:" + name
}

// site is where injected content goes, as a byte offset into the file.
type site struct {
	offset int
	// describes the anchor in a diagnostic, in the author's vocabulary.
	described string
}

// locate finds the point in content where an action's output belongs.
func locate(pkg *genpkg.Package, action *genpkg.Action, content string) (site, error) {
	switch action.Anchor {
	case "":
		return site{}, fmt.Errorf("action %s declares no anchor, so there is nowhere to put its output", action.Name)

	case genpkg.AnchorStartOfFile:
		return site{offset: 0, described: genpkg.AnchorStartOfFile}, nil

	case genpkg.AnchorEndOfFile:
		return site{offset: len(content), described: genpkg.AnchorEndOfFile}, nil

	case genpkg.AnchorMarker:
		return site{}, fmt.Errorf(
			"action %s is anchored to %q without naming a marker", action.Name, genpkg.AnchorMarker)

	case genpkg.AnchorRegion:
		return locateRegion(pkg, action, content)

	case genpkg.AnchorAfterMatch, genpkg.AnchorBeforeMatch:
		return locateMatch(action, content)

	default:
		// Every other value names a marker a file template planted.
		offset, ok := findAnchorLine(pkg.CommentPrefix, action.Anchor, content)
		if !ok {
			return site{}, fmt.Errorf("marker %q is not in the file", action.Anchor)
		}
		return site{offset: offset, described: fmt.Sprintf("marker %q", action.Anchor)}, nil
	}
}

// locateRegion places content at the end of a named region, just before the
// marker that closes it, so that repeated injections into one region accumulate
// in the order they were applied.
func locateRegion(pkg *genpkg.Package, action *genpkg.Action, content string) (site, error) {
	if action.AnchorStart == "" || action.AnchorEnd == "" {
		return site{}, fmt.Errorf(
			"action %s is anchored to a region without naming both anchor_start and anchor_end", action.Name)
	}

	start, ok := findAnchorLineStart(pkg.CommentPrefix, action.AnchorStart, content)
	if !ok {
		return site{}, fmt.Errorf("marker %q, which opens the region, is not in the file", action.AnchorStart)
	}
	end, ok := findAnchorLineStart(pkg.CommentPrefix, action.AnchorEnd, content[start:])
	if !ok {
		return site{}, fmt.Errorf(
			"marker %q, which closes the region opened by %q, is not in the file after it",
			action.AnchorEnd, action.AnchorStart)
	}

	return site{
		offset:    start + end,
		described: fmt.Sprintf("region between %q and %q", action.AnchorStart, action.AnchorEnd),
	}, nil
}

// locateMatch places content relative to a line matching the action's declared
// regex. The offset is a line boundary rather than the match's own bounds, so
// injected content never lands inside an existing line.
//
// The pattern is compiled exactly as the author declared it. An author wanting
// ^ and $ to mean line boundaries rather than the bounds of the whole file
// writes (?m) themselves: silently rewriting a declared expression would make
// the pattern in actions.yaml not the pattern that ran.
func locateMatch(action *genpkg.Action, content string) (site, error) {
	if action.AnchorPattern == "" {
		return site{}, fmt.Errorf("action %s is anchored to %s without declaring anchor_pattern",
			action.Name, action.Anchor)
	}

	re, err := regexp.Compile(action.AnchorPattern)
	if err != nil {
		return site{}, fmt.Errorf("action %s declares anchor_pattern %q, which is not a valid expression: %w",
			action.Name, action.AnchorPattern, err)
	}

	loc := re.FindStringIndex(content)
	if loc == nil {
		return site{}, fmt.Errorf("nothing in the file matches anchor_pattern %q", action.AnchorPattern)
	}

	described := fmt.Sprintf("%s %q", action.Anchor, action.AnchorPattern)
	if action.Anchor == genpkg.AnchorBeforeMatch {
		return site{offset: lineStart(content, loc[0]), described: described}, nil
	}
	return site{offset: lineEnd(content, loc[1]), described: described}, nil
}

// findAnchorLine returns the offset just past the line carrying the named
// anchor marker, which is where content anchored to that marker goes.
func findAnchorLine(commentPrefix, name, content string) (int, bool) {
	start, ok := findAnchorLineStart(commentPrefix, name, content)
	if !ok {
		return 0, false
	}
	return lineEnd(content, start), true
}

// findAnchorLineStart returns the offset of the first byte of the line carrying
// the named anchor marker.
//
// The marker is matched against the line's trimmed text so that a template may
// indent it to suit the file it plants it in, and the name is compared for
// equality rather than as a prefix so that "class_body" does not match
// "class_body_top".
func findAnchorLineStart(commentPrefix, name, content string) (int, bool) {
	decl := anchorDecl(commentPrefix, "")

	for offset := 0; offset < len(content); {
		end := lineTextEnd(content, offset)
		trimmed := strings.TrimSpace(content[offset:end])

		if rest, ok := strings.CutPrefix(trimmed, decl); ok && rest == name {
			return offset, true
		}
		if end == len(content) {
			break
		}
		offset = end + 1
	}
	return 0, false
}

// lineTextEnd returns the offset of the newline ending the line that contains
// offset, or the end of content.
func lineTextEnd(content string, offset int) int {
	if i := strings.IndexByte(content[offset:], '\n'); i >= 0 {
		return offset + i
	}
	return len(content)
}

// lineStart returns the offset of the first byte of the line containing offset.
func lineStart(content string, offset int) int {
	if i := strings.LastIndexByte(content[:offset], '\n'); i >= 0 {
		return i + 1
	}
	return 0
}

// lineEnd returns the offset just past the newline ending the line containing
// offset, or the end of content.
func lineEnd(content string, offset int) int {
	if i := strings.IndexByte(content[offset:], '\n'); i >= 0 {
		return offset + i + 1
	}
	return len(content)
}
