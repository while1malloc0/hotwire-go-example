// Package timefmt implements time formatting utilities
package timefmt

import (
	"strings"
	"time"
)

// FuncMap is a map of the provided functions of this package that can be used
// with the Go template package
var FuncMap = map[string]interface{}{
	"prettyTime": PrettyTime,
}

// PrettyTime returns a truncated version of a RubyTime formatted time.Time
// Note: there's probably a standard-library way of doing this, but Go's time
// formatting is ridiculous, so this was faster than figuring it out.
func PrettyTime(t time.Time) string {
	// e.g. Mon Jan 02 15:04:05 -0700 2006
	s := t.Local().Format(time.RubyDate)
	fields := strings.Split(s, " ")

	// e.g. Jan 02 15:04:05
	return strings.Join(fields[1:4], " ")
}
