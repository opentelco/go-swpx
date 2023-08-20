package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
)

type Stanzas []*stanzapb.Stanza

func (s Stanzas) ToGQL() []*model.Stanza {
	result := make([]*model.Stanza, len(s))
	for i, v := range s {
		result[i] = Stanza{v}.ToGQL()
	}
	return result
}

type Stanza struct{ *stanzapb.Stanza }

func (s Stanza) ToGQL() *model.Stanza {
	st := &model.Stanza{
		ID:             s.Id,
		Name:           &s.Name,
		Description:    &s.Description,
		Template:       &s.Template,
		Content:        &s.Content,
		RevertTemplate: &s.RevertContent,
		RevertContent:  &s.RevertContent,
		DeviceType:     &s.DeviceType,

		UpdatedAt: *internal.TimestampPBToTimePointer(s.UpdatedAt),
		CreatedAt: *internal.TimestampPBToTimePointer(s.CreatedAt),
		AppliedAt: *internal.TimestampPBToTimePointer(s.AppliedAt),
	}
	if s.DeviceId != nil {
		st.Device = &model.Device{
			ID: *s.DeviceId,
		}
	}
	return st

}
