package shared

import (
	"time"

	"go.opentelco.io/go-swpx/proto/go/resourcepb"
)

func ValidateEOLTimeout(req *resourcepb.Request, defaultDuration time.Duration) time.Duration {
	dur, _ := time.ParseDuration(req.Timeout)

	if dur.Seconds() == 0 {
		dur = defaultDuration

	}

	return dur

}
