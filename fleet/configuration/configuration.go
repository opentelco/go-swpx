package configuration

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/mitchellh/hashstructure/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(repo Repository, logger hclog.Logger) configurationpb.ConfigurationServiceServer {
	return &config{
		repo:   repo,
		logger: logger.Named("fleet"),
	}
}

type config struct {
	repo   Repository
	logger hclog.Logger

	configurationpb.UnimplementedConfigurationServiceServer
}

// Get a device configuration by its ID, this is used to get a specific device configuration
func (c *config) GetByID(ctx context.Context, params *configurationpb.GetByIDParameters) (*configurationpb.Configuration, error) {
	if params.Id == "" {
		return nil, ErrInvalidArgumentID
	}

	return c.repo.GetByID(ctx, params.Id)
}

// Compare compares the configuration of a device with the configuration in the database and returns the changes
// if no specific configuration is specified the latest configuration is used to compare with
func (c *config) Compare(ctx context.Context, params *configurationpb.CompareParameters) (*configurationpb.CompareResponse, error) {
	return nil, nil
}

// List lists all configurations for a device
func (c *config) List(ctx context.Context, params *configurationpb.ListParameters) (*configurationpb.ListResponse, error) {
	res, err := c.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	return &configurationpb.ListResponse{
		Configurations: res,
	}, nil
}

// Create a device configuration in the fleet (this is used to store the configuration of a device)
func (c *config) Create(ctx context.Context, params *configurationpb.CreateParameters) (*configurationpb.Configuration, error) {

	if params.Hash == "" {
		hash, err := Hash(params.Configuration)
		if err != nil {
			return nil, fmt.Errorf("could not hash config before save: %w", err)
		}
		params.Hash = hash
	}
	conf := &configurationpb.Configuration{
		DeviceId:      params.DeviceId,
		Configuration: params.Configuration,
		Hash:          params.Hash,
	}

	return c.repo.Upsert(ctx, conf)

}

// Delete a device configuration from the fleet (removes the configuration from the database)
func (c *config) Delete(ctx context.Context, params *configurationpb.DeleteParameters) (*emptypb.Empty, error) {
	err := c.repo.Delete(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// used by the repo to generate a new ID for a device or configuration
func NewID() string {
	guid, _ := uuid.GenerateUUID()
	return guid
}

func Hash(data interface{}) (string, error) {
	hash, err := hashstructure.Hash(data, hashstructure.FormatV2, nil)
	if err != nil {
		return "", fmt.Errorf("counld not create hash of config: %w", err)
	}
	return fmt.Sprintf("%d", hash), nil
}
