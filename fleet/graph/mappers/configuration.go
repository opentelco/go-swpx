package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
)

type Configurations []*configurationpb.Configuration

func (c Configurations) ToGQL() []*model.Configuration {
	result := make([]*model.Configuration, len(c))
	for i, v := range c {
		result[i] = Configuration{v}.ToGQL()
	}
	return result
}

type Configuration struct{ *configurationpb.Configuration }

func (c Configuration) ToGQL() *model.Configuration {
	return &model.Configuration{
		ID: c.Id,
		Device: &model.Device{
			ID: c.DeviceId,
		},
		Configuration: &c.Configuration.Configuration,
		Changes:       &c.Changes,
		Checksum:      &c.Checksum,
		CreatedAt:     c.Created.AsTime(),
	}
}
