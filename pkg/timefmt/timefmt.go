package timefmt

import (
	"strings"
	"time"
)

var FuncMap = map[string]interface{}{
	"prettyTime": PrettyTime,
}

func PrettyTime(t time.Time) string {
	// e.g. Mon Jan 02 15:04:05 -0700 2006
	s := t.Local().Format(time.RubyDate)
	fields := strings.Split(s, " ")

	// e.g. Jan 02 15:04:05
	return strings.Join(fields[1:4], " ")
}
