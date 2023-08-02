package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/commonpb"
)

func ToPageInfo(info *commonpb.PageInfo) *PageInfo {
	return &PageInfo{PageInfo: info}
}

type PageInfo struct{ *commonpb.PageInfo }

func (p PageInfo) ToGQL() *model.PageInfo {
	total := int(p.Total)
	count := int(p.Count)
	return &model.PageInfo{
		Count:  &count,
		Total:  &total,
		Limit:  internal.PointerInt64ToPointerInt(p.Limit),
		Offset: internal.PointerInt64ToPointerInt(p.Offset),
	}
}
