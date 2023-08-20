package configuration

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
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
	cfga, err := c.repo.GetByID(ctx, params.ConfigurationAId)
	if err != nil {
		return nil, fmt.Errorf("could not get configuration A: %w", err)
	}
	cfgb, err := c.repo.GetByID(ctx, params.ConfigurationBId)
	if err != nil {
		return nil, fmt.Errorf("could not get configuration B: %w", err)
	}

	resp, err := c.Diff(ctx, &configurationpb.DiffParameters{
		ConfigurationA:   cfga.Configuration,
		ConfigurationAId: &cfga.Id,
		ConfigurationB:   cfgb.Configuration,
		ConfigurationBId: &cfgb.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create diff between configs: %w", err)
	}

	return &configurationpb.CompareResponse{
		ConfigurationAId: params.ConfigurationAId,
		ConfigurationBId: params.ConfigurationBId,
		Changes:          resp.Changes,
	}, nil

}

// Diff creates a diff between two configurations (strings) and returns the changes
// changes are returned in unified format https://www.gnu.org/software/diffutils/manual/html_node/Unified-Format.html
func (c *config) Diff(ctx context.Context, params *configurationpb.DiffParameters) (*configurationpb.DiffResponse, error) {

	var prevConfigName string = "previous"
	var newConfigName string = "new"
	if params.ConfigurationAId != nil {
		prevConfigName = fmt.Sprintf("previous: %s", *params.ConfigurationAId)
	}

	if params.ConfigurationBId != nil {
		newConfigName = fmt.Sprintf("new: %s", *params.ConfigurationBId)
	}

	edits := myers.ComputeEdits(span.URIFromPath("device-config.cfg"), params.ConfigurationA, params.ConfigurationB)

	diff := fmt.Sprint(gotextdiff.ToUnified(prevConfigName, newConfigName, params.ConfigurationA, edits))
	return &configurationpb.DiffResponse{
		Changes: diff,
	}, nil

}

// List lists all configurations for a device
func (c *config) List(ctx context.Context, params *configurationpb.ListParameters) (*configurationpb.ListResponse, error) {
	res, err := c.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Create a device configuration in the fleet (this is used to store the configuration of a device)
func (c *config) Create(ctx context.Context, params *configurationpb.CreateParameters) (*configurationpb.Configuration, error) {

	if params.Checksum == "" {
		hash, err := Hash(params.Configuration)
		if err != nil {
			return nil, fmt.Errorf("could not hash config before save: %w", err)
		}
		params.Checksum = hash
	}
	conf := &configurationpb.Configuration{
		DeviceId:      params.DeviceId,
		Changes:       params.Changes,
		Configuration: params.Configuration,
		Checksum:      params.Checksum,
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
