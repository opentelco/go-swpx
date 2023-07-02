package shared

import (
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
)

func ValidateEOLTimeout(req *resourcepb.Request, defaultDuration time.Duration) time.Duration {
	dur, _ := time.ParseDuration(req.Timeout)

	if dur.Seconds() == 0 {
		dur = defaultDuration

	}

	return dur

}
