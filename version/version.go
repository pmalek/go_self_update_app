package version

import (
	"strconv"
)

// GetNumber returns the version number as parsed integer from the filename.
//
// NOTE:
// This could be changed to support e.g. semantic versioning but the example
// included an integral version so it was kept that way.
func GetNumber(filename string) (int, error) {
	ret := r.ReplaceAllString(filename, "${Version}")
	v, err := strconv.ParseInt(ret, 10, 64)
	return int(v), err
}
