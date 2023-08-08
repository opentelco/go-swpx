package internal

import (
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/commonpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreatePaginationInfo(limit *int64, offset *int64, totalCount int) *commonpb.PageInfo {
	return &commonpb.PageInfo{
		Total:  int64(totalCount),
		Limit:  limit,
		Offset: offset,
	}
}

func ToPtr[T any](v T) *T {
	return &v
}

func PointerIntToPointerInt64(i *int) *int64 {
	if i == nil {
		return nil
	}
	x := int64(*i)
	return &x
}

func PointerInt64ToPointerInt(i *int64) *int {
	if i == nil {
		return nil
	}
	x := int(*i)
	return &x
}

func TimestampPBToTimePointer(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		x := time.Time{}
		return &x
	}
	x := t.AsTime()
	return &x

}
