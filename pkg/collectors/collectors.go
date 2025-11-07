package collectors

import (
	"context"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

type Config struct {
	ProjectID string
	Region    string
	Profile   string
	Binary    string
}

type Runner struct {
	cfg        Config
	cloudAdmin *CloudAdmin
}

func New(cfg Config) *Runner {
	return &Runner{
		cfg:        cfg,
		cloudAdmin: NewCloudAdmin(cfg),
	}
}

func (r *Runner) Run(ctx context.Context) (model.Snapshot, error) {
	var snap model.Snapshot

	instances, err := r.cloudAdmin.Instances(ctx)
	if err != nil {
		return snap, err
	}
	snap.Instances = instances

	tunnels, err := r.cloudAdmin.OVSTunnels(ctx)
	if err != nil {
		return snap, err
	}
	snap.Tunnels = tunnels

	volumes, err := r.cloudAdmin.Volumes(ctx)
	if err != nil {
		return snap, err
	}
	snap.Volumes = volumes

	return snap, nil
}
