package version

import "regexp"

var (
	r = regexp.MustCompile(`server_v(?P<Version>[0-9]+)\.exe$`)
)
