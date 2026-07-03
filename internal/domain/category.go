package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateCategoryTree checks that every segment of a category tree is safe
// to use as a path component when building an on-disk location for a memo.
//
// A segment is invalid if it is empty, ".", "..", contains a path separator
// ("/" or "\"), or is itself an absolute path. Category trees originate from
// hand-edited YAML frontmatter or free-typed TUI input, so they must be
// rejected outright rather than silently sanitized before being joined into
// a filesystem path.
func ValidateCategoryTree(tree []string) error {
	for _, segment := range tree {
		if segment == "" {
			return fmt.Errorf("invalid category segment %q: segment is empty", segment)
		}
		if segment == "." || segment == ".." {
			return fmt.Errorf("invalid category segment %q: segment must not be %q", segment, segment)
		}
		if strings.ContainsAny(segment, "/\\") {
			return fmt.Errorf("invalid category segment %q: segment must not contain a path separator", segment)
		}
		if filepath.IsAbs(segment) {
			return fmt.Errorf("invalid category segment %q: segment must not be an absolute path", segment)
		}
	}
	return nil
}
